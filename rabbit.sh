#!/bin/bash

start_or_run () {
    docker inspect seekr_rabbitmq > /dev/null 2>&1

    if [ $? -eq 0 ]; then
        echo "Starting SeekR RabbitMQ container..."
        docker start seekr_rabbitmq
    else
        echo "SeekR RabbitMQ container not found, creating a new one..."
        docker run -d --name seekr_rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.13-management
    fi
}

case "$1" in
    start)
        start_or_run
        ;;
    stop)
        echo "Stopping SeekR RabbitMQ container..."
        docker stop seekr_rabbitmq
        ;;
    logs)
        echo "Fetching logs for SeekR RabbitMQ container..."
        docker logs -f seekr_rabbitmq
        ;;
    *)
        echo "Usage: $0 {start|stop|logs}"
        exit 1
esac