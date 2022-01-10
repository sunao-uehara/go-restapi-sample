#!/usr/bin/env bash

op=$1

ABSPATH=$(cd "$(dirname "$0")"; pwd)

start()
{
  echo "start redis"
  docker rm myRedis
  docker pull arm64v8/redis:5.0.14
  docker run -p 6379:6379 -d --name myRedis -v $ABSPATH/data/redis/home/data:/data arm64v8/redis:5.0.14

  echo "start mysql"
  docker rm myMysql
  docker pull --platform linux/x86_64 mysql:5.7 && \
  docker run -d -p 3306:3306 --name myMysql  -v $ABSPATH/data/mysql/home/data:/var/lib/mysql --platform linux/x86_64 -e MYSQL_ROOT_USER=root -e MYSQL_ALLOW_EMPTY_PASSWORD=yes mysql:5.7
}

stop()
{
  echo "stop and remove container"
  docker stop myRedis
  docker stop myMysql
}

restart()
{
  stop
  start
}

DOCKER_P=`docker ps`

if [ -z "$DOCKER_P" ]; then
  echo "Docker is not running"
  exit 0
fi

if [ "$op" = "start" ]; then
  start
elif [ "$op" = "stop" ]; then
  stop
elif [ "$op" = "restart" ]; then
  restart
else
  echo "invalid command: please use [stop|start|restart]"
  exit 0
fi

