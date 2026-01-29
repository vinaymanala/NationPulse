package internal

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"time"
)

type sampleJSON struct {
	ID   int
	Name string
	Date time.Time
}

func GenerateCSV() ([]byte, error) {

	records := []sampleJSON{
		{
			ID:   1,
			Name: "Jack",
			Date: time.Now(),
		},
		{
			ID:   2,
			Name: "John",
			Date: time.Now(),
		},
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	headers := []string{"ID", "Name", "Date"}

	if len(records) > 0 {
		if err := w.Write(headers); err != nil {
			fmt.Println("Error writing header to csv")
			return nil, err
		}
	}

	for _, record := range records {
		row := make([]string, len(headers))
		row = append(row, strconv.Itoa(record.ID), record.Name, record.Date.String())
		if err := w.Write(row); err != nil {
			fmt.Println("Error writing data to csv")
			return nil, err
		}
	}
	w.Flush()
	return buf.Bytes(), nil
}
