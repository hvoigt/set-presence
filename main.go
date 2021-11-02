package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func setStatus(text string, emoji string) []byte {
	token := os.Getenv("SLACK_TOKEN")
	if len(token) == 0 {
		log.Fatalf("Please set SLACK_TOKEN environment variable")
	}

	payload, err := json.Marshal(map[string]string{
		"status_text":       text,
		"status_emoji":      emoji,
		"status_expiration": "0",
	})
	if err != nil {
		log.Fatalln(err)
	}

	profile, err := json.Marshal(map[string]string{
		"profile": string(payload),
	})
	if err != nil {
		log.Fatalln(err)
	}

	url := "https://slack.com/api/users.profile.set"

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(profile))
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	reply, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return reply
}

func setPresence(away bool) []byte {
	token := os.Getenv("SLACK_TOKEN")
	if len(token) == 0 {
		log.Fatalf("Please set SLACK_TOKEN environment variable")
	}

	presence := "auto"
	if away {
		presence = "away"
	}

	payload, err := json.Marshal(map[string]string{
		"presence": presence,
	})
	if err != nil {
		log.Fatalln(err)
	}

	url := "https://slack.com/api/users.setPresence"

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	reply, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return reply
}

func main() {
	weekday := time.Now().Weekday()
	if weekday == time.Monday {
		setStatus("Out of office", ":away:")
		setPresence(true)
	} else {
		setStatus("", "")
		setPresence(false)
	}
}
