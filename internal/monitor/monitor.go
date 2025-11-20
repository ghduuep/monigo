package monitor

import (
	"log"
	"net/http"
	"time"

	"github.com/ghduuep/pingly/models"
)

func CheckSite(url string) (string, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return resp.Status, err
	}
	defer resp.Body.Close()

	return resp.Status, nil
}

func startMonitoring(sites []*models.Site) {
	for _, site := range sites {
		status, err := CheckSite(site.URL)
		if err != nil {
			log.Printf("Cannot connect to website: %v", err)
		}
		if status != site.Status && status != "UNKNOWN" {
			log.Printf("Status changed for %s: %s -> %s", site.URL, site.Status, status)
			// Here you would typically update the site's status in the database
		}
	}
}
