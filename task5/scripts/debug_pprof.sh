#!/bin/bash
echo "=== Debug PProf Setup ==="

# Собираем и запускаем
docker compose up -d --build
echo "Waiting for startup..."
sleep 10

echo "1. Checking if app is running..."
docker compose ps

echo "2. Checking app logs..."
docker compose logs app --tail=20

echo "3. Testing pprof endpoints on port 6060..."

# Проверяем доступность
echo "Testing /debug/pprof/ on port 6060:"
curl -s -o test_index.html -w "HTTP Status: %{http_code}\n" http://localhost:6060/debug/pprof/

echo "Testing /debug/pprof/heap:"
curl -s -o test_heap.pprof -w "HTTP Status: %{http_code}, Size: %{size_download} bytes\n" http://localhost:6060/debug/pprof/heap

echo "Testing /debug/pprof/profile (5 seconds):"
timeout 10 curl -s -o test_cpu.pprof http://localhost:6060/debug/pprof/profile?seconds=5
echo "CPU profile size: $(wc -c < test_cpu.pprof) bytes"

echo "4. File contents check:"
ls -la test_*

docker compose down

echo "=== Debug Complete ==="