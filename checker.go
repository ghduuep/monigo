package main

import (
	"net/http"
	"time"
)

func Check(url string) (string, string, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)

	if err != nil {
		return "DOWN", "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return "UP", resp.Status, nil
	}

	return "DOWN", resp.Status, nil
}
