package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/yuniverse-13/Telegram-Waifu-Bot/internal/bot"
	"github.com/yuniverse-13/Telegram-Waifu-Bot/internal/characters"
)

func main() {
	log.Println("Bot starting...")

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatalln("Bot token is empty. Set environment variable TELEGRAM_BOT_TOKEN.")
	}
	
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("FATAL: DATABASE_URL environment variable is not set.")
	}
	log.Printf("Using DATABASE_URL: %s\n", connStr)
	
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}
	defer db.Close()
	
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to the database!")
	
	charRepo := characters.NewCharacterRepository(db)

	myBot, err := bot.NewBot(botToken, charRepo)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	if err := myBot.Start(); err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}

	log.Println("The bot has finished its work.")
}