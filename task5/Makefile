.PHONY: build run test clean migrate

build:
	docker compose build

run:
	docker compose up

lint:
	golangci-lint run ./...
	go vet ./...

rebuild:
	docker compose down -v
	docker compose up --build
	
run-detached:
	docker compose up -d

stop:
	docker compose down

logs:
	docker compose logs -f app

clean:
	docker compose down -v
	rm -f url-shortener

migrate:
	docker compose exec postgres psql -U postgres -d delayed_notifier -f /docker-entrypoint-initdb.d/init.sql


# Бенчмарки
bench:
	go test -bench=. -benchmem -benchtime=5s ./benchmarks/...

# CPU профилирование
profile-cpu:
	go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30

# Memory profiling
profile-mem:
	go tool pprof http://localhost:8080/debug/pprof/heap

# Trace
trace:
	curl http://localhost:8080/debug/pprof/trace?seconds=5 > trace.out
	go tool trace trace.out

# Load testing
load-test:
	go run scripts/load_test.go

# Сравнение производительности
bench-compare:
	go test -bench=BenchmarkRedirect -count=5 ./benchmarks/... > old.txt
	# вносим изменения
	go test -bench=BenchmarkRedirect -count=5 ./benchmetrics/... > new.txt
	benchstat old.txt new.txt