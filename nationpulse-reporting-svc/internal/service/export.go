package service

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/vinaymanala/nationpulse-reporting-svc/internal/csvSvc"
	"github.com/vinaymanala/nationpulse-reporting-svc/internal/kafkaSvc"
	"github.com/vinaymanala/nationpulse-reporting-svc/internal/store"
	. "github.com/vinaymanala/nationpulse-reporting-svc/internal/types"
)

type ExportService struct {
	configs *Configs
}

func NewExportService(configs *Configs) *ExportService {
	return &ExportService{
		configs: configs,
	}
}

func (es *ExportService) Serve() {

	c := kafkaSvc.NewConsumer(es.configs.Cfg, es.configs.Ctx)
	p := kafkaSvc.NewProducer(es.configs.Cfg, es.configs.Ctx)

	r := c.NewReader("my-group", "message-log")

	// w := p.NewWriter("message-log", "writer1")
	w := p.NewWriter("message-send", "writer2")

	csvService := csvSvc.NewCsvService(es.configs.Cfg, es.configs.Ctx)

	processMessage := func(messageData []byte) error {
		log.Printf("Processing message: %s\n", string(messageData))

		var req store.ExportRequest

		if err := json.Unmarshal(messageData, &req); err != nil {
			log.Printf("Error unmarshaling byte data: %s", err)
			return err
		}

		// fetch from cache
		data, err := GetFromCache(es.configs.Cache, es.configs.Ctx, req.ExportID)

		if data == nil && err == nil {
			log.Printf("No key exist in cache. Trying DB.")
		}
		if err != nil {
			log.Printf("Error received: %s", err)
			return err
		}
		// fmt.Println("DATA STATUSCODE ========>", data.StatusCode)
		if data != nil && data.StatusCode == 1 {
			fmt.Println("DATA FROM CACHE ========>", data)
			resp, err := json.Marshal(data)
			if err != nil {
				log.Printf("Error marshaling data: %s", err)
				return err
			}
			if err = PublishMessage(p, w, resp); err != nil {
				log.Printf("Error publishing message: %s", err)
				return err
			}
			log.Println("Report fetched from cache succesfully.")
			return nil
		}
		log.Printf("Fetching from DB...")
		// fetch from db , generate csv, set cache
		dataFromDB, errDB := GetFromDB(es.configs.DB, es.configs.Ctx, req)
		if errDB != nil {
			return err
		}
		// respDB, err := json.Marshal(dataFromDB)
		// if err != nil {
		// 	log.Printf("Error marshaling data: %s", err)
		// 	return err
		// }
		// if err = PublishMessage(p, w, respDB); err != nil {
		// 	log.Printf("Error publishing message: %s", err)
		// 	return err
		// }

		req.Filters.Records = dataFromDB

		// generate csv
		generatedCsvInBytes, err := csvService.GenerateCSV(req)
		if err != nil {
			return fmt.Errorf("Error generating csv: %w\n", err)
		}

		fmt.Printf("CSV generated in bytes successfully.")

		// set in cache
		response := &store.ExportResponse{
			Status:     "completed",
			StatusCode: 1,
			Data:       generatedCsvInBytes,
		}

		inBytes, err := json.Marshal(response)

		if err != nil {
			log.Printf("Error while marshalling the set cache data %s", err)
			return err
		}
		fmt.Println("DATA FROM  ========>", inBytes)
		if err := es.configs.Cache.SetInCache(es.configs.Ctx, req.ExportID, inBytes); err != nil {
			log.Printf("Error occured while setting data in cache %s", err)
			return err
		}

		// send message
		if err := PublishMessage(p, w, inBytes); err != nil {
			log.Printf("Error publishing message: %s", err)
			return err
		}
		return nil
	}

	c.ReadMessagesWithCallback(r, es.configs, processMessage)

}
