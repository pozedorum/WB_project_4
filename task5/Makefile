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
	rm -f url-shortener pprof/*
	rm -f *.pprof
	bash scripts/free_ports.sh

migrate:
	docker compose exec postgres psql -U postgres -d delayed_notifier -f /docker-entrypoint-initdb.d/init.sql


# Бенчмарки
bench:
	go test -bench=. -benchmem -benchtime=5s ./benchmarks/...

# Запуск всех тестов профилирования
profile-local:  clean
	bash scripts/profile_local.sh
analyze-local: analyze-cpu analyze-mem #analyze-trace

profile-docker: profile-docker-cpu profile-docker-mem profile-docker-trace
analyze-docker:
	bash scripts/analyze_docker_pprof.sh

# Docker CPU профилирование
profile-docker-cpu:
	chmod +x scripts/profile_docker_cpu.sh
	./scripts/profile_docker_cpu.sh

# Docker Memory профилирование
profile-docker-mem:
	chmod +x scripts/profile_docker_memory.sh
	./scripts/profile_docker_memory.sh

# Docker Trace профилирование  
profile-docker-trace:
	chmod +x scripts/profile_docker_trace.sh
	./scripts/profile_docker_trace.sh



# Нагрузочное тестирование (требует запущенного Docker)
load-test-docker:
	docker compose up -d
	sleep 5
	go run scripts/test_load.go
	docker compose down

# Анализ собранных профилей
analyze-cpu:
	go tool pprof -source_path . pprof/local_cpu.pprof

analyze-mem:
	go tool pprof -source_path . pprof/local_heap.pprof

analyze-trace:
	go tool trace pprof/local_trace.out