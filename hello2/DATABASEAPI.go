package main

import (
	"database/sql"
	  "encoding/json"
	"log"
	"net/http"
	"os"
    "fmt"

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

type Chat struct {
	ChatID    int    `json:"chat_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Type      string `json:"type"`
}

var db *sql.DB

func initDB() error {
	var err error
	db, err = sql.Open("sqlite", "telegram.db")
	if err != nil {
		log.Fatal("Error initializing database: ", err)
	}

	return db.Ping()
}

func insertMessage(msg message) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback on error

	// Check if chat already exists
	var count int
	err = tx.QueryRow("SELECT COUNT(*) FROM chats WHERE id = ?", msg.Chat.ID).Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking for existing chat: %w", err)
	}

	// Only insert into chats if it doesn't already exist
	if count == 0 {
		_, err = tx.Exec("INSERT INTO chats (id, first_name, last_name) VALUES (?, ?, ?)", msg.Chat.ID, msg.Chat.FirstName, msg.Chat.LastName)
		if err != nil {
			return fmt.Errorf("error inserting chat: %w", err)
		}
	}

	_, err = tx.Exec("INSERT INTO telegram (id, text, chat_id) VALUES (?, ?, ?)", msg.MessageID, msg.Text, msg.Chat.ID)
	if err != nil {
		return fmt.Errorf("error inserting message: %w", err)
	}

	return tx.Commit()
}


func getAllMessagesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, text, chat_id FROM telegram")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []ChatMessage
	for rows.Next() {
		var msg ChatMessage
		if err := rows.Scan(&msg.MessageID, &msg.Text, &msg.ChatID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		messages = append(messages, msg)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(messages)
}

func getChatsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, first_name, last_name FROM chats")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var chats []Chat
	for rows.Next() {
		var chat Chat
		if err := rows.Scan(&chat.ChatID, &chat.FirstName, &chat.LastName); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		chats = append(chats, chat)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(chats)
}

func inputMessagesHandler(rw http.ResponseWriter, rq *http.Request) {
	var p TextMessage
	decoder := json.NewDecoder(rq.Body)
	if err := decoder.Decode(&p); err != nil {
		http.Error(rw, "Invalid JSON", http.StatusBadRequest)
		return
	}
	for index, element := range p.Result {

		log.Printf("%d) Inserting data with update_id = %d", index, element.UpdateID)
		insertMessage(element.Message)
	}
}





func main() {
    	log.SetPrefix("DBMS: ")
	f, err := os.OpenFile("logs.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)

	err = initDB()
	if err != nil {
		log.Fatal("Error while Ping Database: ", err)
	}

	http.HandleFunc("/api/input_messages", inputMessagesHandler)
	http.HandleFunc("/api/messages", getAllMessagesHandler) 
	http.HandleFunc("/api/getchats", getChatsHandler) 

	

	http.ListenAndServe(":8080", nil)
}