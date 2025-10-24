#!/bin/bash
echo "=== Docker Trace Profiling ==="

docker compose up -d
sleep 10

echo "Collecting trace (5 seconds)..."
curl -s "http://localhost:6060/debug/pprof/trace?seconds=5" > pprof/trace_docker.out

docker compose down

echo "Trace saved to trace_docker.out"
echo "Analyze with: go tool trace trace_docker.out"