#!/bin/bash
echo "=== Local Profiling ==="

# Функция для проверки порта
check_port() {
    lsof -i :$1 >/dev/null 2>&1
}

# Функция для освобождения порта
free_port() {
    local port=$1
    local pid=$(lsof -i :$port -t 2>/dev/null)
    if [ -n "$pid" ]; then
        echo "Freeing port $port (PID: $pid)"
        kill -9 $pid 2>/dev/null
        sleep 2
    fi
}

echo "0. Checking and freeing ports..."
free_port 6060
free_port 8080

# Проверяем что порты свободны
if check_port 6060; then
    echo "Error: Port 6060 is still in use"
    echo "Please run: kill -9 \$(lsof -i :6060 -t)"
    exit 1
fi

if check_port 8080; then
    echo "Error: Port 8080 is still in use" 
    echo "Please run: kill -9 \$(lsof -i :8080 -t)"
    exit 1
fi

echo "1. Starting server..."
docker compose up postgres &
sleep 5

go run cmd/main.go &
SERVER_PID=$!
sleep 5

# Проверяем что сервер запустился
if ! ps -p $SERVER_PID > /dev/null; then
    echo "Error: Server failed to start"
    docker compose down
    exit 1
fi

sleep 3

echo "2. Starting load test..."
go run scripts/test_load.go &
LOAD_PID=$!
sleep 2

echo "3. Collecting profiles..."
echo "   CPU profile (30s)..."
curl -s "http://localhost:8080/debug/pprof/profile?seconds=30" > pprof/local_cpu.pprof

echo "   Memory profile..."
curl -s "http://localhost:8080/debug/pprof/heap" > pprof/local_heap.pprof

echo "   Trace (5s)..."
curl -s "http://localhost:8080/debug/pprof/trace?seconds=5" > pprof/local_trace.out

echo "4. Stopping services..."
kill $SERVER_PID $LOAD_PID 2>/dev/null
wait $SERVER_PID 2>/dev/null
docker compose down

echo "Profiles saved to pprof/"