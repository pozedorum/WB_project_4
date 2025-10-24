
### Первый бенчмарк
BenchmarkShortenHandler-8         401496             14543 ns/op            8127 B/op              40 allocs/op
BenchmarkRedirectHandler-8        408236             14077 ns/op            8756 B/op              29 allocs/op
BenchmarkCreateShortURL-8        1923670              2614 ns/op             428 B/op               8 allocs/op
BenchmarkGenerateShortURL-8     10018942               584.0 ns/op           144 B/op               5 allocs/op
PASS
ok      github.com/pozedorum/WB_project_4/task5/benchmarks      26.651s


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