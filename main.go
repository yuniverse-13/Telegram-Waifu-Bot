package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/yuniverse-13/Telegram-Waifu-Bot/internal/bot"
	"github.com/yuniverse-13/Telegram-Waifu-Bot/internal/characters"
	"github.com/yuniverse-13/Telegram-Waifu-Bot/internal/ratings"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	log.Println("Bot starting...")
	
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatalln("Bot token is empty. Set environment variable TELEGRAM_BOT_TOKEN.")
	}
	
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("FATAL: DATABASE_URL environment variable is not set.")
	}
	log.Printf("Using DATABASE_URL: %s\n", connStr)
	
	gormDB, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	
	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying sql.DB from GORM: %v", err)
	}
	defer sqlDB.Close()
	
	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to the database!")
	
	charRepo   := characters.NewCharacterRepository(gormDB)
	ratingRepo := ratings.NewRepository(gormDB)

	myBot, err := bot.NewBot(botToken, charRepo, ratingRepo)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	if err := myBot.Start(); err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}
	log.Println("The bot has finished its work.")
}