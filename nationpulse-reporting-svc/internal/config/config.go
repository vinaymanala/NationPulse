package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	// Kafka
	KafkaBrokers       []string
	KafkaConsumerTopic string
	KafkaProducerTopic string
	// Port
	Port int
	//Redis
	RedisAddr string
	RedisPass string
	RedisDB   int
	//Postgres
	PostgresHost string
	PostgresPass string
	PostgresName string
	PostgresUser string
	PostgresAddr string
}

func Load() Config {
	defaultPort := 8080
	defaultRedisDB := 0

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "localhost:9092"
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6378"
	}
	redisPass := os.Getenv("REDIS_PASS")
	redisDB := os.Getenv("REDIS_DB")
	if redisDB == "" {
		redisDB = "0"
	}

	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = defaultPort
	}

	redisDBInt, err := strconv.Atoi(redisDB)
	if err != nil {
		redisDBInt = defaultRedisDB
	}

	pgHost := os.Getenv("PG_DB_HOST")
	if pgHost == "" {
		pgHost = "postgres-db"
	}

	pgName := os.Getenv("PG_DB_NAME")
	if pgName == "" {
		pgName = "nationPulseDB"
	}

	pgUser := os.Getenv("PG_DB_USER")
	if pgUser == "" {
		pgUser = "postgres"
	}

	pgPass := os.Getenv("PG_DB_PASS")
	if pgPass == "" {
		pgPass = "postgres"
	}

	pgAddr := os.Getenv("PG_DB_ADDR")
	if pgAddr == "" {
		pgAddr = "localhost:5432"
	}
	kafkaReaderTopic := os.Getenv("KAFKA_READER_TOPIC")
	if kafkaReaderTopic == "" {
		kafkaReaderTopic = "message-send"
	}
	kafkaWriterTopic := os.Getenv("KAFKA_WRITER_TOPIC")
	if kafkaWriterTopic == "" {
		kafkaWriterTopic = "message-log"
	}

	return Config{
		KafkaBrokers:       strings.Split(kafkaBrokers, ","),
		KafkaConsumerTopic: kafkaReaderTopic,
		KafkaProducerTopic: kafkaWriterTopic,
		Port:               port,
		RedisAddr:          redisAddr,
		RedisPass:          redisPass,
		RedisDB:            redisDBInt,
		PostgresHost:       pgHost,
		PostgresName:       pgName,
		PostgresUser:       pgUser,
		PostgresPass:       pgPass,
		PostgresAddr:       pgAddr,
	}
}
