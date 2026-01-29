docker exec -it kafka-1 kafka-topics.sh --create --topic my-replicated-topic --partitions 3 --replication-factor 3 --bootstrap-server localhost:9092

/opt/kafka/bin/kafka-topics.sh
