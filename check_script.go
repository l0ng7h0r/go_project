package main

import (
	"database/sql"
	"fmt"
	"log"
	
	"github.com/l0ng7h0r/golang/pkg/config"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()
	db, err := sql.Open("postgres", cfg.DBDsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Try inserting a dummy refresh token for an existing user
	var userID string
	err = db.QueryRow("SELECT id FROM users LIMIT 1").Scan(&userID)
	if err != nil {
		fmt.Println("No users found or error:", err)
		return
	}

	fmt.Println("Testing insert into refresh_tokens for user:", userID)
	_, err = db.Exec(`INSERT INTO refresh_tokens(user_id, token, expires_at) VALUES ($1,$2,NOW() + INTERVAL '7 days')`, userID, "dummy_token_123")
	if err != nil {
		fmt.Println("ERROR inserting refresh token:", err)
	} else {
		fmt.Println("SUCCESS inserting refresh token!")
	}
}
