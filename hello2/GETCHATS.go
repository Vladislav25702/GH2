package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type ChatMessage struct {
	MessageID int    `json:"message_id"`
	ChatID    int    `json:"chat_id"`
	Text      string `json:"text"`
}

type Chat struct {
	ChatID    int    `json:"chat_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Type      string `json:"type"`
}


func getChatsFromDatabaseAPI(url string) ([]Chat, error) {
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

	var chats []Chat
	err = json.Unmarshal(body, &chats)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON response: %w", err)
	}

	return chats, nil
}



func getChatsHandler(w http.ResponseWriter, r *http.Request) {
	databaseAPIURL := "http://localhost:8080/api/getchats"
	chats, err := getChatsFromDatabaseAPI(databaseAPIURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(chats)
}

func main() {
log.SetPrefix("GetMessagesService: ")
	f, err := os.OpenFile("get_messages_logs.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	http.HandleFunc("/api/chats", getChatsHandler)

	log.Println("GetMessages service listening on :8080")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Failed to start GetMessages service: %v", err)
	}



	
	log.Println("GetMessages service listening on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Failed to start GetMessages service: %v", err)
	}
}