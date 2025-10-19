package reminder

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/pozedorum/WB_project_4/task3/internal/interfaces"
	"github.com/pozedorum/WB_project_4/task3/internal/models"
	"github.com/pozedorum/WB_project_4/task3/pkg/config"
	"github.com/pozedorum/wbf/rabbitmq" // Замени на актуальный путь
	"github.com/pozedorum/wbf/retry"
	"github.com/rabbitmq/amqp091-go"
)

type rabbitMQClient struct {
	conn         *rabbitmq.Connection
	channel      *rabbitmq.Channel
	publisher    *rabbitmq.Publisher
	queueManager *rabbitmq.QueueManager
	exchange     *rabbitmq.Exchange
	queueName    string
}

func NewRabbitMQClient(cfg config.RabbitMQConfig) (interfaces.QueueClient, error) {
	// Устанавливаем соединение с RabbitMQ
	conn, err := rabbitmq.Connect(cfg.URL, 5, 2*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// Создаем канал
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}

	// Создаем обменник
	exchange := rabbitmq.NewExchange(cfg.Exchange, "direct")
	exchange.Durable = true

	if err = exchange.BindToChannel(ch); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Создаем менеджер очередей
	queueManager := rabbitmq.NewQueueManager(ch)

	// Объявляем ОДНУ основную очередь с DLX
	queue, err := queueManager.DeclareQueue(cfg.QueueName, rabbitmq.QueueConfig{
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		Args: amqp091.Table{
			// DLX для обработки сообщений с истекшим TTL
			"x-dead-letter-exchange":    cfg.Exchange,
			"x-dead-letter-routing-key": "reminders.process", // routing key для обработки после TTL
		},
	})
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	// Привязываем основную очередь к обменнику для отложенных сообщений
	err = ch.QueueBind(
		queue.Name,
		"reminders.delayed", // routing key для отложенных сообщений
		cfg.Exchange,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to bind delayed queue: %w", err)
	}

	// Привязываем ТУ ЖЕ очередь к обменнику для немедленной обработки
	err = ch.QueueBind(
		queue.Name,
		"reminders.process", // routing key для обработки после TTL
		cfg.Exchange,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to bind process queue: %w", err)
	}

	// Создаем публикатор
	publisher := rabbitmq.NewPublisher(ch, cfg.Exchange)

	client := &rabbitMQClient{
		conn:         conn,
		channel:      ch,
		publisher:    publisher,
		queueManager: queueManager,
		exchange:     exchange,
		queueName:    cfg.QueueName,
	}

	return client, nil
}

func (c *rabbitMQClient) PublishReminder(reminder models.ReminderMessage) error {
	jsonData, err := json.Marshal(reminder)
	if err != nil {
		return fmt.Errorf("failed to marshal reminder: %w", err)
	}

	// Вычисляем задержку для напоминания
	now := time.Now()
	var routingKey string
	var options []rabbitmq.PublishingOptions

	if reminder.NotifyTime.After(now) {
		// Если напоминание в будущем - отправляем в отложенную очередь
		delay := reminder.NotifyTime.Sub(now)
		routingKey = "reminders.delayed"

		options = []rabbitmq.PublishingOptions{
			{
				Expiration: delay,
			},
		}
	} else {
		// Если напоминание уже должно быть отправлено - отправляем сразу
		routingKey = "reminders.immediate"
	}

	// Публикуем сообщение с ретраями
	strategy := retry.Strategy{
		Attempts: 3,
		Delay:    1 * time.Second,
	}

	return c.publisher.PublishWithRetry(
		jsonData,
		routingKey,
		"application/json",
		strategy,
		options...,
	)
}

func (c *rabbitMQClient) StartConsuming(ctx context.Context, handler func(models.ReminderMessage) error) error {
	// Создаем консьюмер
	consumerConfig := rabbitmq.NewConsumerConfig(c.queueName)
	consumerConfig.AutoAck = false // Ручное подтверждение
	consumerConfig.Consumer = "reminder-worker"

	consumer := rabbitmq.NewConsumer(c.channel, consumerConfig)

	// Канал для получения сообщений
	msgChan := make(chan []byte, 100)

	// Запускаем потребление в горутине
	go func() {
		strategy := retry.Strategy{
			Attempts: 5,
			Delay:    2 * time.Second,
		}

		if err := consumer.ConsumeWithRetry(msgChan, strategy); err != nil {
			log.Printf("Failed to start consumer: %v", err)
		}
	}()

	// Обрабатываем сообщения
	for {
		select {
		case <-ctx.Done():
			close(msgChan)
			return ctx.Err()
		case msgBody, ok := <-msgChan:
			if !ok {
				return fmt.Errorf("message channel closed")
			}

			var reminder models.ReminderMessage
			if err := json.Unmarshal(msgBody, &reminder); err != nil {
				log.Printf("Failed to unmarshal reminder message: %v", err)
				continue
			}

			if err := handler(reminder); err != nil {
				log.Printf("Failed to process reminder: %v", err)
				// В production системе здесь должна быть логика повторной обработки или перемещения в DLQ
			}
		}
	}
}

func (c *rabbitMQClient) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
