package kafkaSvc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/vinaymanala/nationpulse-reporting-svc/internal/config"
	"github.com/vinaymanala/nationpulse-reporting-svc/internal/types"
)

type Consumer struct {
	cfg config.Config
	ctx context.Context
}

func NewConsumer(cfg config.Config, ctx context.Context) *Consumer {
	return &Consumer{
		cfg: cfg,
		ctx: ctx,
	}
}

func (c *Consumer) NewReader(groupID, topic string) *kafka.Reader {
	readLog := log.New(os.Stdout, "Reader: ", 0)
	fmt.Println("Kafka Brokers", c.cfg.KafkaBrokers)
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        c.cfg.KafkaBrokers,
		Topic:          topic,
		GroupID:        groupID,
		StartOffset:    kafka.LastOffset,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
		SessionTimeout: 20 * time.Second,
		// ReadTimeout:       10 * time.Second, // Increase from default
		ReadBackoffMin:    100 * time.Millisecond,
		ReadBackoffMax:    1 * time.Second,
		HeartbeatInterval: 3 * time.Second,
		RebalanceTimeout:  60 * time.Second,
		MaxWait:           2 * time.Second,
		Logger:            readLog,
	})
}

func (c *Consumer) ReadMessagesWithCallback(r *kafka.Reader, configs *types.Configs, callbackFn func([]byte) error) {

	defer r.Close()
	// backoff := 1

	for {
		select {
		case <-c.ctx.Done():
			fmt.Println("Consumer shutting down...")
			return
		default:
		}

		// Use 30-60 seconds, NOT 1 second
		readCtx, cancel := context.WithTimeout(c.ctx, 2*time.Minute)
		msg, err := r.ReadMessage(readCtx)
		cancel()

		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				continue
			}
			fmt.Printf("Error reading message: %v\n", err)
			// time.Sleep(time.Duration(backoff) * time.Second)
			// backoff = min(backoff*2, 30)
			continue
		}

		// backoff = 1
		err = callbackFn(msg.Value)
		if err != nil {
			fmt.Printf("Error processing message: %v\n", err)
			continue
		}

		fmt.Printf("received - key: %s, value: %s\n", string(msg.Key), string(msg.Value))

		if err := r.CommitMessages(c.ctx, msg); err != nil {
			fmt.Println("Error committing", err)
		}
	}

	// return &msg, nil
}
