#!/bin/bash


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


echo "Checking and freeing ports..."
free_port 5432
free_port 6060
free_port 8080
