#!/bin/bash
docker container stop todo_test
docker rm todo_test
docker container stop mysql
docker rm mysql

docker network create -d bridge todo_net
docker run --rm -d --network todo_net --name mysql -e MYSQL_ROOT_PASSWORD=test -p 3306:3306 mysql:latest
while ! docker exec -i mysql mysql -uroot -ptest <<< "CREATE DATABASE todo_test;"
do
	echo "Connecting to MySQL (errors are normal)..."
	sleep 5
done
APP_PORT=1818 docker run --rm -it --network todo_net --name todo_test -p 1818:1818 todo_test
docker ps
