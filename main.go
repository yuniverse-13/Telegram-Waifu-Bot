package main

import (
	"log"
	"os"
	
	"github.com/youniverse-13/Telegram-Waifu-Bot/internal/bot"
)

func main() {
	log.Println("Bot starting...")
	
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatalln("Bot token is empty")
	}
	
	myBot, err := bot.NewBot(botToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}
	
	if err := myBot.Start(); err != nil {
		log.Fatalf("Bot encountered an error: %v", err)
	}
	
	log.Println("Bot finished.")
}