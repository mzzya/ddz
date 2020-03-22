#! /bin/bash

# 关闭
docker-compose stop
# 清理
docker rm $(docker-compose ps -q)
# 启动
docker-compose up -d --remove-orphans

#docker network create --subnet 172.88.88.0/24 simple_api
# 测试kafka队列正常输出
# docker exec -it kafka1 bash
## /opt/kafka/bin/kafka-console-consumer.sh --topic simple_api --bootstrap-server kafka1:9092
## kafka-topics.sh --describe --topic simple_api --zookeeper zookeeper:2181
