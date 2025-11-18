package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Website struct {
	URL        string
	Interval   time.Duration
	LastStatus string
}

func main() {
	gfarma := Website{
		URL:        "https://gfarma.com",
		Interval:   1 * time.Minute,
		LastStatus: "UNKNOWN",
	}

	cafefacil := Website{
		URL:        "https://cafefacil.com.br",
		Interval:   1 * time.Minute,
		LastStatus: "UNKNOWN",
	}

	titantelecom := Website{
		URL:        "https://titantelecom.com.br",
		Interval:   5 * time.Minute,
		LastStatus: "UNKNOWN",
	}
	go startMonitoring(gfarma)
	go startMonitoring(cafefacil)
	go startMonitoring(titantelecom)

	select {}
}

func checkWebsite(url string) (string, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)

	if err != nil {
		return "DOWN", err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return "UP", nil
	}

	return "DOWN", nil
}

func startMonitoring(website Website) {
	fmt.Printf("Iniciando monitoramento para %s\n", website.URL)

	for {
		status, err := checkWebsite(website.URL)
		if err != nil {
			log.Printf("Erro ao verificar %s: %v", website.URL, err)
		}

		if status != website.LastStatus && website.LastStatus != "UNKNOWN" {
			log.Printf("Status alterado para %s para o site %s", status, website.URL)
		}
		website.LastStatus = status
		time.Sleep(website.Interval)
	}
}
