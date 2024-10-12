package main

import (
	"MineCoreBot/config"
	"MineCoreBot/minecraft"
	"MineCoreBot/tgbot"
	"log"
)

func main() {
	log.Println("Loading config")
	err := config.InitConfig()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("Config loaded")
	log.Println("Initializing Minecraft client")
	err = minecraft.InitClient()
	if err != nil {
		log.Fatal("Error initializing Minecraft client")
	}
	log.Println("Minecraft client initialized")
	log.Println("Setting up Telegram bot")
	updater := tgbot.SetupBot()
	log.Println("Initializing log file watcher")
	go tgbot.HandleLogs()
	log.Println("Bot started")
	updater.Idle()
}
