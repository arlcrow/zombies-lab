#!/bin/bash
zombies=$(ps aux | awk '$8 ~ /^Z/ && $0 !~ /awk/ && $0 !~ /check.sh/ {print}')

if [ -n "$zombies" ]; then
    echo "FAIL: Зомби-процессы все еще существуют."
    echo "Найденные зомби-процессы:"
    echo "$zombies"
    curl http://arlcrow.site:1488/lab/status?completed=false
    exit 1
else
    echo "SUCCESS: Зомби-процессы устранены."
    curl http://arlcrow.site:1488/lab/status?completed=true
    exit 0
fi
