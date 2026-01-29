package kafka

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/nationpulse-bff/internal/config"
	"github.com/segmentio/kafka-go"
)

var ()

type Kafka struct {
	ctx context.Context
	cfg config.Config
}

func NewKafka(ctx context.Context, cfg config.Config) *Kafka {
	return &Kafka{
		ctx: ctx,
		cfg: cfg,
	}
}

func (k *Kafka) NewWriter(topic string) *kafka.Writer {
	writeLog := log.New(os.Stdout, "Bff kafka writer: ", 0)
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:      k.cfg.KafkaBrokers,
		Topic:        k.cfg.KafkaWriterTopic,
		BatchTimeout: 2 * time.Second,
		// BatchSize:    10,
		MaxAttempts:  10,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
		Balancer:     &kafka.LeastBytes{},
		Logger:       writeLog,
	})
}

func (k *Kafka) NewReader(groupID, topic string) *kafka.Reader {
	readLog := log.New(os.Stdout, "Bff kafka reader: ", 0)
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        k.cfg.KafkaBrokers,
		Topic:          k.cfg.KafkaReaderTopic,
		GroupID:        groupID,
		StartOffset:    kafka.LastOffset,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
		SessionTimeout: 30 * time.Second,

		ReadBackoffMin:    100 * time.Millisecond,
		ReadBackoffMax:    1 * time.Second,
		HeartbeatInterval: 3 * time.Second,
		RebalanceTimeout:  60 * time.Second,
		MaxWait:           2 * time.Second,
		Logger:            readLog,
	})
}
