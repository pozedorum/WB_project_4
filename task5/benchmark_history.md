
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

После закидывания логов в нейронку, она скомпоновала их в следующий результат:

Выявленные проблемные места:
1. Медленные SQL-запросы (57.36% времени)

GetOriginalURLIfExists - 1.12s (28.94%) - SELECT по short_code

RegisterClick - 1.10s (28.42%) - INSERT + UPDATE в транзакции

2. Блокирующие операции

Redirect хендлер - 1.39s (35.92%) блокирует ответ

Синхронные вызовы к БД в основном потоке