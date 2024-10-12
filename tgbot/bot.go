package tgbot

import (
	"log"
	"time"

	"MineCoreBot/config"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

var tgbot *gotgbot.Bot

func SetupBot() *ext.Updater {
	var err error

	tgbot, err = gotgbot.NewBot(config.BotToken, nil)
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Println("an error occurred while handling tgbot update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})
	updater := ext.NewUpdater(dispatcher, nil)

	dispatcher.AddHandler(handlers.NewCommand("start", StartCommand))
	dispatcher.AddHandler(handlers.NewCommand("console", ConsoleCommand))
	dispatcher.AddHandler(handlers.NewCommand("whitelist", WhiteListCommand))
	dispatcher.AddHandler(handlers.NewCommand("ban", BanCommand))
	dispatcher.AddHandler(handlers.NewCommand("reload", ReloadCommand))
	dispatcher.AddHandler(handlers.NewCommand("stop", StopCommand))
	dispatcher.AddHandler(handlers.NewCommand("help", HelpCommand))
	dispatcher.AddHandler(handlers.NewMessage(message.Text, HandleChatTopic))

	err = updater.StartPolling(tgbot, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 10,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 40,
			},
		},
	})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}
	log.Printf("%s has been started...\n", tgbot.User.Username)
	return updater
}
