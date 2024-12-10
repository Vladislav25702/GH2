package main

import (
	  "bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type FrontendMessage struct {
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}



const telegramAPI = "https://api.telegram.org/bot6971995963:AAH8uBq41VNlrm0BtKmk54s4_mN5ZcaksG0/sendMessage"


// Function to send a message to Telegram
func sendToTelegram(chatId int, text string) error {
	telegramMsg := struct {
		ChatID int    `json:"chat_id"`
		Text   string `json:"text"`
	}{
		ChatID: chatId,
		Text:   text,
	}

	jsonData, err := json.Marshal(telegramMsg)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	resp, err := http.Post(telegramAPI, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("error sending POST request to Telegram: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return fmt.Errorf("Telegram API returned status: %d - %s", resp.StatusCode, bodyString)
	}
	return nil
}

// API handler to receive messages from the frontend
func sendMessageHandler(w http.ResponseWriter, r *http.Request) {


	

	
	var msg FrontendMessage
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if err := sendToTelegram(msg.ChatID, msg.Text); err != nil {
		http.Error(w, fmt.Sprintf("Error sending message to Telegram: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}





func main() {
	
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stdout) //Send logs to standard output

	http.HandleFunc("/api/send", sendMessageHandler)
	log.Println("Telegram API service listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))

	
	
}

