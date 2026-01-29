package internal

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

// var (
// 	topic          = "message-log"
// 	broker1Address = "localhost:9092"
// 	broker2Address = "localhost:9094"
// 	broker3Address = "localhost:9095"
// )

func Produce(ctx context.Context) {

	i := 0

	l := log.New(os.Stdout, "kafka writer: ", 0)
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker1Address, broker2Address, broker3Address},
		// Topic:        topic,
		MaxAttempts:  10,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		Balancer:     &kafka.LeastBytes{},
		Logger:       l,
	})
	defer w.Close()

	for {

		writeCtx, cancel := context.WithTimeout(ctx, 30*time.Second)

		err := w.WriteMessages(writeCtx, kafka.Message{
			Key:   []byte(strconv.Itoa(i)),
			Value: []byte("this is message:" + strconv.Itoa(i)),
		})
		cancel()

		if err != nil {
			fmt.Printf("Error writing message %d: %v\n", i, err)
			// Don't panic, wait and retry
			time.Sleep(2 * time.Second)
			continue
		}

		fmt.Println("writes:", i)
		i++
		time.Sleep(time.Second)
	}
}

func Consume(ctx context.Context) {
	l := log.New(os.Stdout, "kafka reader: ", 0)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker1Address, broker2Address, broker3Address},
		// Topic:          topic,
		GroupID:        "my-group",
		StartOffset:    kafka.LastOffset,
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
		SessionTimeout: 20 * time.Second,
		Logger:         l,
	})
	defer r.Close()
	for {
		readCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)

		msg, err := r.ReadMessage(readCtx)
		cancel()

		if err != nil {
			fmt.Printf("Error reading message: %v\n", err)
			// Don't panic, wait and retry
			time.Sleep(2 * time.Second)
			continue
		}

		fmt.Printf("received - key: %s, value: %s\n", string(msg.Key), string(msg.Value))
	}
}
