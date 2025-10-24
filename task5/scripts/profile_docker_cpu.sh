#!/bin/bash
echo "=== Docker CPU Profiling ==="

# Запускаем сервисы
docker compose up -d
echo "Waiting for services to start..."
sleep 10

echo "Starting load test..."
go run scripts/test_load.go &

echo "Collecting CPU profile (30 seconds)..."
curl -s "http://localhost:6060/debug/pprof/profile?seconds=30" > pprof/cpu_docker.pprof

echo "Stopping services..."
docker compose down

echo "CPU profile saved to cpu_docker.pprof"
echo "Analyze with: go tool pprof cpu_docker.pprof"