package reminder

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pozedorum/WB_project_4/task3/internal/interfaces"
	"github.com/pozedorum/WB_project_4/task3/internal/models"
	"github.com/streadway/amqp"
)

type rabbitMQClient struct {
	conn            *amqp.Connection
	channel         *amqp.Channel
	queueName       string
	dlxExchangeName string
	dlxQueueName    string
	logger          interfaces.Logger
}

func NewRabbitMQClient(url string, queueName string, logger interfaces.Logger) (interfaces.QueueClient, error) {
	logger.Info("RABBITMQ_INIT", "Initializing RabbitMQ client",
		"url", url, "queue", queueName)

	// Устанавливаем соединение с RabbitMQ
	conn, err := amqp.Dial(url)
	if err != nil {
		logger.Error("RABBITMQ_INIT", "Failed to connect to RabbitMQ",
			"error", err, "url", url)
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	logger.Info("RABBITMQ_INIT", "Connected to RabbitMQ successfully")

	// Создаем канал
	ch, err := conn.Channel()
	if err != nil {

		err2 := conn.Close()

		logger.Error("RABBITMQ_INIT", "Failed to create channel",
			"error", err)
		if err2 != nil {
			logger.Error("RABBITMQ_INIT", "Failed to close connection",
				"error", err2)
		}
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}
	logger.Info("RABBITMQ_INIT", "Channel created successfully")

	// Имена для DLX
	dlxExchangeName := queueName + "_dlx"
	dlxQueueName := queueName + "_delayed"

	// 1. Сначала создаем DLX exchange
	err = ch.ExchangeDeclare(
		dlxExchangeName,
		"direct", // type
		true,     // durable
		false,    // autoDelete
		false,    // internal
		false,    // noWait
		nil,      // args
	)
	if err != nil {
		err2, err3 := ch.Close(), conn.Close()
		logger.Error("RABBITMQ_INIT", "Failed to declare DLX exchange",
			"error", err, "exchange", dlxExchangeName)
		if err2 != nil {
			logger.Error("RABBITMQ_INIT", "Failed to close channel",
				"error", err2)
		}
		if err3 != nil {
			logger.Error("RABBITMQ_INIT", "Failed to close connection",
				"error", err3)
		}
		return nil, fmt.Errorf("failed to declare DLX exchange: %w", err)
	}
	logger.Info("RABBITMQ_INIT", "DLX exchange declared",
		"exchange", dlxExchangeName)

	// 2. Создаем DLX очередь для отложенных сообщений
	_, err = ch.QueueDeclare(
		dlxQueueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		amqp.Table{
			"x-dead-letter-exchange":    "",                         // Используем default exchange
			"x-dead-letter-routing-key": queueName,                  // возвращаем в основную очередь
			"x-message-ttl":             int32(24 * 60 * 60 * 1000), // 24 часа макс TTL
		},
	)
	if err != nil {
		logger.Error("RABBITMQ_INIT", "Failed to declare DLX queue",
			"error", err, "queue", dlxQueueName)
		err2, err3 := ch.Close(), conn.Close()

		if err2 != nil {
			logger.Error("RABBITMQ_INIT", "Failed to close channel",
				"error", err2)
		}
		if err3 != nil {
			logger.Error("RABBITMQ_INIT", "Failed to close connection",
				"error", err3)
		}
		return nil, fmt.Errorf("failed to declare DLX queue: %w", err)
	}
	logger.Info("RABBITMQ_INIT", "DLX queue declared",
		"queue", dlxQueueName)

	// 3. Привязываем DLX очередь к DLX exchange
	err = ch.QueueBind(
		dlxQueueName,
		"", // routing key
		dlxExchangeName,
		false,
		nil,
	)
	if err != nil {
		err2, err3 := ch.Close(), conn.Close()

		if err2 != nil {
			logger.Error("RABBITMQ_INIT", "Failed to close channel",
				"error", err2)
		}
		if err3 != nil {
			logger.Error("RABBITMQ_INIT", "Failed to close connection",
				"error", err3)
		}
		return nil, fmt.Errorf("failed to bind DLX queue: %w", err)
	}
	logger.Info("RABBITMQ_INIT", "DLX queue bound to exchange",
		"queue", dlxQueueName, "exchange", dlxExchangeName)

	// 4. Создаем основную очередь (куда будут возвращаться сообщения после TTL)
	_, err = ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	if err != nil {
		logger.Error("RABBITMQ_INIT", "Failed to declare main queue",
			"error", err, "queue", queueName)
		err2, err3 := ch.Close(), conn.Close()

		if err2 != nil {
			logger.Error("RABBITMQ_INIT", "Failed to close channel",
				"error", err2)
		}
		if err3 != nil {
			logger.Error("RABBITMQ_INIT", "Failed to close connection",
				"error", err3)
		}
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}
	logger.Info("RABBITMQ_INIT", "Main queue declared",
		"queue", queueName)

	client := &rabbitMQClient{
		conn:            conn,
		channel:         ch,
		queueName:       queueName,
		dlxExchangeName: dlxExchangeName,
		dlxQueueName:    dlxQueueName,
		logger:          logger,
	}

	logger.Info("RABBITMQ_INIT", "RabbitMQ client initialized successfully")
	return client, nil
}

func (c *rabbitMQClient) PublishReminder(reminder models.ReminderMessage) error {
	jsonData, err := json.Marshal(reminder)
	if err != nil {
		c.logger.Error("RABBITMQ_PUBLISH", "Failed to marshal reminder",
			"error", err, "event_id", reminder.EventID)
		return fmt.Errorf("failed to marshal reminder: %w", err)
	}

	now := time.Now()
	var expiration string

	if reminder.NotifyTime.After(now) {
		// Напоминание в будущем - устанавливаем TTL
		delay := reminder.NotifyTime.Sub(now)
		expiration = fmt.Sprintf("%d", delay.Milliseconds())
		c.logger.Debug("RABBITMQ_PUBLISH", "Scheduling future reminder",
			"event_id", reminder.EventID,
			"delay_ms", delay.Milliseconds(),
			"notify_time", reminder.NotifyTime)
	} else {
		// Напоминание сейчас или в прошлом - минимальный TTL
		expiration = "1000" // 1 секунда
		c.logger.Debug("RABBITMQ_PUBLISH", "Scheduling immediate reminder",
			"event_id", reminder.EventID)
	}

	// ВАЖНО: Публикуем в DLX очередь с TTL
	err = c.channel.Publish(
		"",             // exchange (default)
		c.dlxQueueName, // DLX очередь для ожидания
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         jsonData,
			Expiration:   expiration, // TTL для отложенной доставки
			DeliveryMode: amqp.Persistent,
			Timestamp:    now,
			Headers: amqp.Table{
				"x-event-id":    reminder.EventID,
				"x-notify-time": reminder.NotifyTime,
			},
		},
	)
	if err != nil {
		c.logger.Error("RABBITMQ_PUBLISH", "Failed to publish reminder",
			"error", err, "event_id", reminder.EventID, "dlx_queue", c.dlxQueueName)
		return fmt.Errorf("failed to publish reminder: %w", err)
	}

	c.logger.Info("RABBITMQ_PUBLISH", "Reminder published to DLX queue",
		"event_id", reminder.EventID,
		"dlx_queue", c.dlxQueueName,
		"expiration_ms", expiration,
		"notify_time", reminder.NotifyTime,
		"current_time", now)

	return nil
}

// Остальные методы (StartConsuming, processMessage, Close) остаются без изменений
func (c *rabbitMQClient) StartConsuming(ctx context.Context, handler func(models.ReminderMessage) error) error {
	c.logger.Info("RABBITMQ_CONSUME", "Starting consumer",
		"queue", c.queueName) // Потребляем из ОСНОВНОЙ очереди

	// Настраиваем QoS - только одно сообщение в обработке за раз
	err := c.channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		c.logger.Error("RABBITMQ_CONSUME", "Failed to set QoS",
			"error", err)
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	// Начинаем потребление сообщений из ОСНОВНОЙ очереди
	msgs, err := c.channel.Consume(
		c.queueName,       // Основная очередь
		"reminder-worker", // consumer
		false,             // autoAck (ручное подтверждение)
		false,             // exclusive
		false,             // noLocal
		false,             // noWait
		nil,               // args
	)
	if err != nil {
		c.logger.Error("RABBITMQ_CONSUME", "Failed to start consuming",
			"error", err, "queue", c.queueName)
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	c.logger.Info("RABBITMQ_CONSUME", "Consumer started successfully",
		"queue", c.queueName)

	// Обрабатываем сообщения
	go func() {
		defer func() {
			if r := recover(); r != nil {
				c.logger.Error("RABBITMQ_CONSUME", "Consumer panic",
					"panic", r)
			}
		}()

		for {
			select {
			case <-ctx.Done():
				c.logger.Info("RABBITMQ_CONSUME", "Consumer stopped by context")
				return

			case msg, ok := <-msgs:
				if !ok {
					c.logger.Info("RABBITMQ_CONSUME", "Messages channel closed")
					return
				}

				// Обрабатываем сообщение
				c.processMessage(msg, handler)
			}
		}
	}()

	return nil
}

func (c *rabbitMQClient) processMessage(delivery amqp.Delivery, handler func(models.ReminderMessage) error) {
	var reminder models.ReminderMessage

	c.logger.Debug("RABBITMQ_PROCESS", "Received message",
		"message_id", delivery.MessageId,
		"body_size", len(delivery.Body),
		"headers", delivery.Headers)

	if err := json.Unmarshal(delivery.Body, &reminder); err != nil {
		c.logger.Error("RABBITMQ_PROCESS", "Failed to unmarshal reminder",
			"error", err, "message_id", delivery.MessageId)
		// Подтверждаем даже при ошибке парсинга, чтобы не застревало
		if err := delivery.Ack(false); err != nil {
			c.logger.Error("RABBITMQ_PROCESS", "Failed to execute ack command",
				"error", err)
		}
		return
	}

	c.logger.Info("RABBITMQ_PROCESS", "Processing reminder",
		"event_id", reminder.EventID,
		"notify_time", reminder.NotifyTime,
		"telegram_id", reminder.TelegramID)

	// Обрабатываем напоминание
	if err := handler(reminder); err != nil {
		c.logger.Error("RABBITMQ_PROCESS", "Failed to process reminder",
			"error", err, "event_id", reminder.EventID)
		// Перепотребляем при ошибке обработки
		if err := delivery.Nack(false, true); err != nil {
			c.logger.Error("RABBITMQ_PROCESS", "Failed to execute nack command",
				"error", err)
		}
	} else {
		c.logger.Info("RABBITMQ_PROCESS", "Successfully processed reminder",
			"event_id", reminder.EventID)
		// Подтверждаем успешную обработку
		if err := delivery.Ack(false); err != nil {
			c.logger.Error("RABBITMQ_PROCESS", "Failed to execute ack command",
				"error", err)
		}
	}
}

func (c *rabbitMQClient) Close() {
	c.logger.Info("RABBITMQ_CLOSE", "Closing RabbitMQ client")

	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			c.logger.Error("RABBITMQ_CLOSE", "Failed to close channel",
				"error", err)
		} else {
			c.logger.Info("RABBITMQ_CLOSE", "Channel closed successfully")
		}
	}

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			c.logger.Error("RABBITMQ_CLOSE", "Failed to close connection",
				"error", err)
		} else {
			c.logger.Info("RABBITMQ_CLOSE", "Connection closed successfully")
		}
	}
}
