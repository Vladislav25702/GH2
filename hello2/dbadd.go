package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    
    "log"
    "net/http"
 		"time"

    _ "modernc.org/sqlite"
)


// TextMessage структура JSON полученных текстовых сообщений
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
	Text      string `json:"text"`
	ChatID    int `json:"chat_id"` //Simplified struct for database
}

var lastUpdateId uint32 = 0
var telegramBotAPI = "https://api.telegram.org/bot6971995963:AAH8uBq41VNlrm0BtKmk54s4_mN5ZcaksG0/getUpdates"

func SaveMessages() (bool, TextMessage, error) {
	url := fmt.Sprintf("%s?offset=%d", telegramBotAPI, lastUpdateId+1)
	resp, err := http.Get(url)
	if err != nil {
		return false, TextMessage{}, fmt.Errorf("error sending GET request: %w", err)
	}
	defer resp.Body.Close() //Important to close the response body

	if resp.StatusCode != http.StatusOK {
		return false, TextMessage{}, fmt.Errorf("Telegram API returned status code: %d", resp.StatusCode)
	}

	var P TextMessage
	if err := json.NewDecoder(resp.Body).Decode(&P); err != nil {
		return false, TextMessage{}, fmt.Errorf("error decoding JSON response: %w", err)
	}

	if len(P.Result) == 0 {
		return false, TextMessage{}, nil // No new messages
	}
	lastUpdateId = uint32(P.Result[len(P.Result)-1].UpdateID)
	return true, P, nil
}

func sendInputRequest(resp TextMessage) error {
	url := "http://localhost:8080/api/input_messages"
	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("error serializing data: %w", err)
	}
	postResponse, err := http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("error sending POST request: %w", err)
	}
	defer postResponse.Body.Close()
	if postResponse.StatusCode != http.StatusOK {
		return fmt.Errorf("POST request failed with status code: %d", postResponse.StatusCode)
	}
	return nil
}

func main() {
	for {
		NewMessages, telegramResponse, err := SaveMessages()
		if err != nil {
			log.Printf("Error processing messages: %v", err)
		} else if NewMessages {
			log.Println("Отправить...")
			if err := sendInputRequest(telegramResponse); err != nil {
				log.Printf("Error sending input request: %v", err)
			}
		} else {
			log.Println("Нет новых сообщений...")
		}
		time.Sleep(15 * time.Second)
	}
}