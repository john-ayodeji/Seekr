#!/bin/bash

CONTAINER_NAME=seekr_postgres
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=seekr_db
POSTGRES_PORT=5432
VOLUME_NAME=seekr_postgres_data

start_or_run () {
    docker inspect $CONTAINER_NAME > /dev/null 2>&1

    if [ $? -eq 0 ]; then
        echo "Starting SeekR Postgres container..."
        docker start $CONTAINER_NAME
    else
        echo "SeekR Postgres container not found, creating a new one..."
        docker run -d \
          --name $CONTAINER_NAME \
          -e POSTGRES_USER=$POSTGRES_USER \
          -e POSTGRES_PASSWORD=$POSTGRES_PASSWORD \
          -e POSTGRES_DB=$POSTGRES_DB \
          -p $POSTGRES_PORT:5432 \
          -v $VOLUME_NAME:/var/lib/postgresql/data \
          postgres:16
    fi
}

case "$1" in
    start)
        start_or_run
        ;;
    stop)
        echo "Stopping SeekR Postgres container..."
        docker stop $CONTAINER_NAME
        ;;
    logs)
        echo "Fetching logs for SeekR Postgres container..."
        docker logs -f $CONTAINER_NAME
        ;;
    *)
        echo "Usage: $0 {start|stop|logs}"
        exit 1
esac
