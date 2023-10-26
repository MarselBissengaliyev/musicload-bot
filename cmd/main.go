package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/MarselBissengaliyev/realtime-chat/bot"
	config "github.com/MarselBissengaliyev/realtime-chat/configs"
)

func main() {
	config, err := config.InitConfig("configs")

	if err != nil {
		log.Fatal("Error loading .env file", err.Error())
		return
	}

	bot, err := bot.NewTelegramBot(
		config.TelegramApiToken,
		config.MaxDownloadTime,
		config.MaxVideoDuration,
		config.BotUsername,
	)

	if err != nil {
		log.Printf("Error initializing bot: %s\n", err.Error())
		return
	}

	sigint := make(chan os.Signal)

	signal.Notify(sigint, syscall.SIGTERM)
	signal.Notify(sigint, syscall.SIGINT)

	go func() {
		<-sigint

		fmt.Println("Gracefully stopping application")
		bot.Stop()

		os.Exit(1)
	}()

	if err := bot.Run(config.Debug); err != nil {
		fmt.Printf("Error occured while running main event loop: %s", err.Error())
	}
}
