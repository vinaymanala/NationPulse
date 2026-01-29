package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/vinaymanala/nationpulse-reporting-svc/internal/config"
	"github.com/vinaymanala/nationpulse-reporting-svc/internal/service"
	"github.com/vinaymanala/nationpulse-reporting-svc/internal/store"
	. "github.com/vinaymanala/nationpulse-reporting-svc/internal/types"
)

func main() {

	// Load environment variables from .env for local development
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or failed to load; relying on environment variables")
	}

	cfg := config.Load()

	// r := gin.Default()
	ctx := context.Background()

	rds := store.NewRedis(cfg)
	db := store.NewPgClient(ctx, cfg)

	configs := &Configs{
		Cfg:   cfg,
		Ctx:   ctx,
		DB:    db,
		Cache: rds,
	}
	exportService := service.NewExportService(configs)
	exportService.Serve()

	// r.POST("/produce", func(c *gin.Context) {
	// 	type Req struct {
	// 		Message string `json:"message"`
	// 	}
	// 	var body Req
	// 	if err := c.BindJSON(&body); err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	writer := internal.NewKafkaWriter()
	// 	err := writer.WriteMessages(context.Background(),
	// 		kafka.Message{
	// 			Key:   []byte(time.Now().Format(time.RFC3339)),
	// 			Value: []byte(body.Message),
	// 		},
	// 	)
	// 	internal.CloseWriter(writer)

	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to produce"})
	// 		return
	// 	}
	// 	c.JSON(http.StatusOK, gin.H{"status": "sent"})
	// })

	// go startConsumer("group1") // start consumer in background

	// log.Println("Starting gin at :8080")
	// if err := r.Run(":8080"); err != nil {
	// 	panic("Error occured while running the server:" + err.Error())
	// }
}
