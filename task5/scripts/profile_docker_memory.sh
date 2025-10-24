#!/bin/bash
echo "=== Docker Memory Profiling ==="

docker compose up -d
sleep 10

echo "Collecting heap profile..."
curl -s "http://localhost:6060/debug/pprof/heap" > pprof/heap_docker.pprof

echo "Running load test..."
go run scripts/test_load.go

echo "Collecting heap after load..."
curl -s "http://localhost:6060/debug/pprof/heap" > pprof/heap_after_load_docker.pprof

docker compose down

echo "Memory profiles saved"
echo "Analyze with: go tool pprof heap_docker.pprof"