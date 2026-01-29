package csvSvc

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"log"

	"github.com/vinaymanala/nationpulse-reporting-svc/internal/config"
	"github.com/vinaymanala/nationpulse-reporting-svc/internal/store"
)

type CsvSvc struct {
	cfg config.Config
	ctx context.Context
}

func NewCsvService(cfg config.Config, ctx context.Context) *CsvSvc {
	return &CsvSvc{
		cfg: cfg,
		ctx: ctx,
	}
}

func (c *CsvSvc) GenerateCSV(req store.ExportRequest) ([]byte, error) {

	headers := req.Filters.Headers
	records := req.Filters.Records

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	if len(records) > 0 {
		if err := w.Write(headers); err != nil {
			log.Printf("Error writing headers: %s", err)
			return nil, err
		}
	}

	for _, record := range records {
		row := make([]string, len(headers))
		for i := range headers {
			if i < len(record) {
				row[i] = record[i]
			} else {
				row[i] = ""
			}
		}
		fmt.Println(record)
		// row = append(row, record...)
		if err := w.Write(row); err != nil {
			log.Printf("Error writing a record: %s", err)
			return nil, err
		}
	}
	w.Flush()

	// fmt.Println(headers, records)
	return buf.Bytes(), nil
}
