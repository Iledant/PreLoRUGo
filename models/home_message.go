package models

import "database/sql"

import "fmt"

// HomeMessage model
type HomeMessage struct {
	Title NullString `json:"Title"`
	Body  NullString `json:"Body"`
}

// Get fetches the home message from database
func (h *HomeMessage) Get(db *sql.DB) error {
	err := db.QueryRow(`SELECT title,body FROM home_message`).Scan(&h.Title, &h.Body)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

// Set insert or update the send message
func (h *HomeMessage) Set(db *sql.DB) error {
	if _, err := db.Exec(`DELETE from home_message`); err != nil {
		return fmt.Errorf("delete %v")
	}
	_, err := db.Exec(`INSERT INTO home_message (title,body) VALUES($1,$2)`,
		h.Title, h.Body)
	return err
}
