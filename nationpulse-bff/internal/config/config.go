package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	AccessSecret     string
	RefreshSecret    string
	KafkaBrokers     []string
	KafkaReaderTopic string
	KafkaWriterTopic string
	Port             int
	RedisAddr        string
	RedisPass        string
	RedisDB          int
	PostgresHost     string
	PostgresPass     string
	PostgresName     string
	PostgresUser     string
	PostgresAddr     string
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
		redisAddr = "localhost:6379"
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
	accessSecret := os.Getenv("ACCESS_SECRET")
	if accessSecret == "" {
		accessSecret = "epAWkXwJr03d6MN9pXUAxGj+dB6U+8HcYxwTGxVU1gg="
	}
	refreshSecret := os.Getenv("REFRESH_SECRET")
	if refreshSecret == "" {
		refreshSecret = `YN5Vgber0y5IX5CUq\ZlIQC8avwpsX67+Q\8OcvJ23c=`
	}

	return Config{
		AccessSecret:     accessSecret,
		RefreshSecret:    refreshSecret,
		KafkaBrokers:     strings.Split(kafkaBrokers, ","),
		KafkaReaderTopic: kafkaReaderTopic,
		KafkaWriterTopic: kafkaWriterTopic,
		Port:             port,
		RedisAddr:        redisAddr,
		RedisPass:        redisPass,
		RedisDB:          redisDBInt,
		PostgresHost:     pgHost,
		PostgresName:     pgName,
		PostgresUser:     pgUser,
		PostgresPass:     pgPass,
		PostgresAddr:     pgAddr,
	}
}
