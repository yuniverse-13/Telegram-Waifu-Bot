package main

import (
	"log"
	"os"

	"github.com/yuniverse-13/Telegram-Waifu-Bot/internal/bot"
)

func main() {
	log.Println("Bot starting...")

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatalln("Bot token is empty. Set environment variable TELEGRAM_BOT_TOKEN.")
	}

	myBot, err := bot.NewBot(botToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	if err := myBot.Start(); err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}

	log.Println("The bot has finished its work.")
}
