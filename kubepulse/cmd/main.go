package main

import (
	"fmt"
	"os"
	"time"

	"github.com/kubepulse/internal"
)

func main() {

	// List all service to health check
	services := map[string]string{
		"BFF-Service":   os.Getenv("BFF_ADDR"),
		"Reporting-Svc": os.Getenv("REPORTING_ADDR"),
		"Ingestion-Svc": os.Getenv("INGESTION_ADDR"),
		"Cronjob-Svc":   os.Getenv("CRONJOB_ADDR"),
	}

	for {
		fmt.Println("\n--- Pulse check:", time.Now().Format(time.Kitchen), "---")
		for name, addr := range services {
			go func() {
				status := internal.CheckHealth(addr)
				fmt.Printf("[%s]  Status:%s\n", name, status)
			}()
		}
		time.Sleep(10 * time.Second)
	}
}
