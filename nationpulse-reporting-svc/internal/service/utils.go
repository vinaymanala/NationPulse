package service

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/vinaymanala/nationpulse-reporting-svc/internal/kafkaSvc"
	"github.com/vinaymanala/nationpulse-reporting-svc/internal/store"
)

func GetFromCache(redis *store.Redis, ctx context.Context, exportID string) (*store.ExportResponse, error) {
	// check in cache
	cacheData, err := redis.CheckInCache(ctx, exportID)
	if cacheData == nil && err == nil {
		return nil, nil
	}
	fmt.Println("DATA FROM CACHE ========>", cacheData)
	if err != nil || cacheData == nil {
		log.Printf("Error received: ExportID: %s, err: %s", exportID, err)
		return nil, err
	}
	return cacheData, nil
}

func GetFromDB(pg *store.PgClient, ctx context.Context, req store.ExportRequest) ([][]string, error) {
	data, err := pg.FetchFromDB(ctx, req)
	if err != nil {
		log.Printf("Error occured fetching from DB: %s", err)
		return nil, err
	}
	return data, nil
}

func PublishMessage(p *kafkaSvc.Producer, w *kafka.Writer, resp []byte) error {
	if err := p.WriteMessage(w, resp); err != nil {
		return fmt.Errorf("error sending generated csv: %w", err)
	}
	return nil
}
