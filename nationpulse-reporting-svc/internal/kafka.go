package internal

import (
	"time"

	"github.com/segmentio/kafka-go"
)

var (
	KafKaTopic     = "orders"
	broker1Address = "localhost:9092"
	broker2Address = "localhost:9094"
	broker3Address = "localhost:9095"
)

// Producer Writer
func NewKafkaWriter(topic string) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{broker1Address, broker2Address, broker3Address},
		Topic:        topic,
		MaxAttempts:  10,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		Balancer:     &kafka.LeastBytes{},
		// Logger:       l,
	})
}

// Consumer Reader
func NewKafkaReader(groupId, topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{broker1Address, broker2Address, broker3Address},
		Topic:          topic,
		GroupID:        groupId,
		StartOffset:    kafka.LastOffset,
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
		SessionTimeout: 20 * time.Second,
		// Logger:         l,
	})
}

// Close writer
func CloseWriter(w *kafka.Writer) error { return w.Close() }
