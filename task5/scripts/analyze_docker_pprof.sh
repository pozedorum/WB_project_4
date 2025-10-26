#!/bin/bash
echo "=== CPU - All Project Functions ==="
go tool pprof -top pprof/cpu_docker.pprof 2>/dev/null | grep "WB_project_4" | head -20

echo -e "\n=== Memory - All Project Functions ==="
go tool pprof -top -sample_index=inuse_space pprof/heap_after_load_docker.pprof 2>/dev/null | grep "WB_project_4" | head -20

echo -e "\n=== Memory Growth - Project Functions ==="
go tool pprof -top -base pprof/heap_docker.pprof pprof/heap_after_load_docker.pprof 2>/dev/null | grep "WB_project_4" | head -15