#!/bin/bash
set -e

API_URL=${API_URL:-"http://localhost:14880"}

check_zombies() {
    zombies=$(ps aux | awk '$8 ~ /^Z/ && $0 !~ /awk/ && $0 !~ /check.sh/ {print}')
    echo "Checking for zombie processes..."
    
    if [ -n "$zombies" ]; then
        echo "FAIL: Zombie processes still exist"
        echo "Found zombie processes:"
        echo "$zombies"
        curl -X POST -H "Content-Type: application/json" \
             -d '{"completed":false,"message":"Zombie processes detected"}' \
             "${API_URL}/lab/status"
        exit 1
    else
        echo "SUCCESS: No zombie processes found"
        curl -X POST -H "Content-Type: application/json" \
             -d '{"completed":true,"message":"All zombie processes eliminated"}' \
             "${API_URL}/lab/status"
        exit 0
    fi
}

check_zombies
