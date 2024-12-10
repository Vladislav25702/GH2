package main

import (
	
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"io"

	_ "modernc.org/sqlite"
)

type TextMessage struct {
	Ok     bool `json:"ok"`
	Result []struct {
		UpdateID int `json:"update_id"`
		Message  struct {
			MessageID int `json:"message_id"`
			From      struct {
				ID           int    `json:"id"`
				IsBot        bool   `json:"is_bot"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
			} `json:"from"`
			Chat struct {
				ID        int    `json:"id"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Type      string `json:"type"`
			} `json:"chat"`
			Date int    `json:"date"`
			Text string `json:"text"`
		} `json:"message"`
	} `json:"result"`
}

type message struct {
	MessageID int `json:"message_id"`
	From      struct {
		ID           int    `json:"id"`
		IsBot        bool   `json:"is_bot"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		LanguageCode string `json:"language_code"`
	} `json:"from"`
	Chat struct {
		ID        int    `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Type      string `json:"type"`
	} `json:"chat"`
	Date int    `json:"date"`
	Text string `json:"text"`
}

type ChatMessage struct {
	MessageID int    `json:"message_id"`
	ChatID    int    `json:"chat_id"`
	Text      string `json:"text"`
}


func getMessagesFromDatabaseAPI(url string) ([]ChatMessage, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request to database API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("database API returned non-200 status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var messages []ChatMessage
	err = json.Unmarshal(body, &messages)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON response: %w", err)
	}

	return messages, nil
}

func getMessagesHandler(w http.ResponseWriter, r *http.Request) {
	databaseAPIURL := "http://localhost:8080/api/getmessages" // Address of the Database API service
	messages, err := getMessagesFromDatabaseAPI(databaseAPIURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(messages)
}


func main() {
	log.SetPrefix("GetMessagesService: ")
	f, err := os.OpenFile("get_messages_logs.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	http.HandleFunc("/api/messages", getMessagesHandler)

	log.Println("GetMessages service listening on :8080")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Failed to start GetMessages service: %v", err)
	}
}