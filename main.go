package main

import (
	"log"
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

	go monitorLoop(&gfarma)
	go monitorLoop(&cafefacil)
	go monitorLoop(&titantelecom)

	select {}
}

func monitorLoop(website *Website) {
	for {
		newStatus, err := Check(website.URL)

		if err != nil {
			log.Printf("[%s] Erro: %v", website.URL, err)
		}

		if website.LastStatus != "UNKNOWN" && website.LastStatus != newStatus {
			log.Printf("MUDANÇA DE STATUS: %s - %s", website.URL, newStatus)

			go func() {
				err := sendEmailNotification(website.URL, newStatus)
				if err != nil {
					log.Printf("Erro ao enviar e-mail de aviso: %v", err)
				}
			}()
		} else {
			log.Printf("[%s] Status: %s", website.URL, newStatus)
		}

		website.LastStatus = newStatus
		time.Sleep(website.Interval)
	}
}
