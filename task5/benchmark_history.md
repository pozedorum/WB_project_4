
### Первый бенчмарк
BenchmarkShortenHandler-8         401496             14543 ns/op            8127 B/op              40 allocs/op
BenchmarkRedirectHandler-8        408236             14077 ns/op            8756 B/op              29 allocs/op
BenchmarkCreateShortURL-8        1923670              2614 ns/op             428 B/op               8 allocs/op
BenchmarkGenerateShortURL-8     10018942               584.0 ns/op           144 B/op               5 allocs/op
PASS
ok      github.com/pozedorum/WB_project_4/task5/benchmarks      26.651s

=== Load Test: URL Shortening ===
Duration: 936.875309ms
Requests: 1000
Success: 1000
Errors: 0
RPS: 1067.38
Avg Latency: 81.97 ms

=== Load Test: URL Redirects ===
Duration: 2.056715961s
Requests: 2000
Success: 2000
Errors: 0
RPS: 972.42
Avg Latency: 47.76 ms

## Анализ функций через pprof
bash scripts/profile_local.sh
make analyze-all

После закидывания логов в нейронку(их слишком много, чтобы вставить в отчёт), она скомпоновала их в следующий результат:

Выявленные проблемные места:
1. Медленные SQL-запросы (57.36% времени)

GetOriginalURLIfExists - 1.12s (28.94%) - SELECT по short_code

RegisterClick - 1.10s (28.42%) - INSERT + UPDATE в транзакции

2. Блокирующие операции

Redirect хендлер - 1.39s (35.92%) блокирует ответ   

Синхронные вызовы к БД в основном потоке


## Предпринятые изменения

1. Переписывание функции repository.RegisterClick, чтобы избавиться от транзакции и выполнять все действия в один запрос

2. Замена индексов в базе данных на конкурентные варианты (поскольку теперь транзакций нет, можно испольовать этот тип индексов)
```sql
CREATE INDEX CONCURRENTLY idx_short_urls_short_code ON short_urls(short_code);
```

3. Меняем настройки подключения к бд, увеличивая размер пула соединений
```go
opts := &dbpg.Options{
		MaxOpenConns:    50,
		MaxIdleConns:    20,
		ConnMaxLifetime: time.Hour,
	}
```

## Следующий анализ pprof

BenchmarkShortenHandler-8         397604             13759 ns/op            8127 B/op         40 allocs/op
BenchmarkRedirectHandler-8        412764             12605 ns/op            7562 B/op         29 allocs/op
BenchmarkCreateShortURL-8        1894822              2672 ns/op             430 B/op          8 allocs/op
BenchmarkGenerateShortURL-8     10362616               615.1 ns/op           144 B/op          5 allocs/op

// Показания локальных тестов
2. Starting load test...
=== Load Test: URL Shortening ===
Duration: 1.09754523s
Requests: 1000
Success: 1000
Errors: 0
RPS: 911.12
Avg Latency: 100.24 ms

=== Load Test: URL Redirects ===
3. Collecting profiles...
   CPU profile (30s)...
Duration: 1.989989013s
Requests: 2000
Success: 2000
Errors: 0
RPS: 1005.03
Avg Latency: 46.32 ms
   Memory profile...
   Trace (5s)...
4. Stopping services...

// Показания тестов в докере
=== Load Test: URL Shortening ===
Duration: 888.715183ms
Requests: 1000
Success: 1000
Errors: 0
RPS: 1125.22
Avg Latency: 78.78 ms

=== Load Test: URL Redirects ===
Duration: 1.894530207s
Requests: 2000
Success: 2000
Errors: 0
RPS: 1055.67
Avg Latency: 43.68 ms

## Предпринятые изменения

4. Изменил принцип сбора аналитики на пакетную (сначала собирается пакет аналитики, а затем вставляется в базу), благодаря чему уменьшилось количество обращений к базе данных.

// Показания локальных тестов
=== Load Test: URL Shortening ===
Duration: 958.473024ms
Requests: 1000
Success: 1000
Errors: 0
RPS: 1043.33
Avg Latency: 88.41 ms

=== Load Test: URL Redirects ===
3. Collecting profiles...
   CPU profile (30s)...
Duration: 1.745955958s
Requests: 2000
Success: 2000
Errors: 0
RPS: 1145.50
Avg Latency: 40.88 ms

// Показания тестов в докере
=== Load Test: URL Shortening ===
Duration: 827.700044ms
Requests: 1000
Success: 1000
Errors: 0
RPS: 1208.17
Avg Latency: 75.36 ms

=== Load Test: URL Redirects ===
Duration: 2.045792718s
Requests: 2000
Success: 2000
Errors: 0
RPS: 977.62
Avg Latency: 46.58 ms



