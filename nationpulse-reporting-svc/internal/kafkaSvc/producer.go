package kafkaSvc

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/vinaymanala/nationpulse-reporting-svc/internal/config"
)

type Producer struct {
	cfg config.Config
	ctx context.Context
}

func NewProducer(cfg config.Config, ctx context.Context) *Producer {
	return &Producer{
		cfg: cfg,
		ctx: ctx,
	}
}

func (p *Producer) NewWriter(topic string, writerStr string) *kafka.Writer {
	writeLog := log.New(os.Stdout, writerStr+": ", 0)
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:      p.cfg.KafkaBrokers,
		Topic:        topic,
		BatchTimeout: 2 * time.Second,
		// BatchSize:    10,
		MaxAttempts:  10,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
		Balancer:     &kafka.LeastBytes{},
		Logger:       writeLog,
	})
}

func (p *Producer) WriteMessage(w *kafka.Writer, data []byte) error {
	writeCtx, cancel := context.WithTimeout(p.ctx, 30*time.Second)

	err := w.WriteMessages(writeCtx, kafka.Message{
		Key:   []byte("Export generated"),
		Value: data,
	})
	cancel()

	if err != nil {
		return fmt.Errorf("error writing message: %w", err)
	}
	fmt.Println("Write Successfull============")
	return nil
}
