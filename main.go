package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const factURL = "https://uselessfacts.jsph.pl/api/v2/facts/random?language=en"

func fetchFact() (string, error) {
	var body struct {
		Text string `json:"text"`
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(factURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", err
	}
	if body.Text == "" {
		return "", fmt.Errorf("no text field in response")
	}
	return body.Text, nil
}

func main() {
	a, err := fetchFact()
	if err != nil {
		return
	}
	fmt.Println(a)
}
