package main

import (
	"database/sql"
	"fmt"
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
	
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	sslMode := os.Getenv("DB_SSLMODE")
	
	if dbHost == "" { dbHost = "localhost" }
	if dbPort == "" { dbPort = "5432" }
	if dbUser == "" { dbUser = "waifu_bot_user" }
	if dbPassword == "" { log.Fatalln("DB_PASSWORD is not set") }
	if dbName == "" { dbName = "waifu_bot_db" }
	if sslMode == "" { sslMode = "disable" }
		
	connStr := fmt.Sprintf("host=%s port=%s user=%s password='%s' dbname=%s sslmode=%s", dbHost, dbPort, dbUser, dbPassword, dbName, sslMode)
	
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
	//characters.SeedInitialData(charRepo)

	myBot, err := bot.NewBot(botToken, charRepo)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	if err := myBot.Start(); err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}

	log.Println("The bot has finished its work.")
}