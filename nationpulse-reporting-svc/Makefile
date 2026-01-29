.PHONY: zookeeper kafka topic up down

zookeeper:
# Runs zookeeper in a container
	docker run -d --name zookeeper \
	-p 2181:2181 \
	-e ALLOW_ANONYMOUS_LOGIN=yes \
	zookeeper


kafka:
# Runs kafka in a container
	docker run -d --name kafka -p 9092:9092 \
  	-e KAFKA_ZOOKEEPER_CONNECT=localhost:2181 \
  	-e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092 \
  	apache/kafka

topic:
# Creates a kafka topic with configurations
	docker exec kafka kafka-topics.sh --create --topic orders --bootstrap-server localhost:9092 --partitions 3 --replication-factor 1

up: zookeeper kafka topic

down: 
	docker rm -f kafka zookeeper