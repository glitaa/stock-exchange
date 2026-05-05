#!/bin/bash

set -e

PORT=${1:-8080}

export APP_PORT=$PORT

echo "Starting Stock Exchange on port: $PORT..."
echo "Use the -d flag if you want to run containers in the background."

docker compose up --build