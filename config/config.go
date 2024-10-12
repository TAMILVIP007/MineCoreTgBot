package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	RconnPort        string
	RconnHost        string
	RconnPassword    string
	BotToken         string
	ChatId           int64
	ChatTopic        int64
	DeathMsgTopic    int64
	EventTopic       int64
	LogFilePath      string
	OwnerId		  int64
	AchievementTopic int64
	DefaultTopic    int64
)

func InitConfig() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return err
	}
	RconnPort = os.Getenv("RCONN_PORT")
	RconnHost = os.Getenv("RCONN_HOST")
	RconnPassword = os.Getenv("RCONN_PASSWORD")
	BotToken = os.Getenv("BOT_TOKEN")
	ChatTopic = ConvertInt64(os.Getenv("CHAT_TOPIC"))
	DeathMsgTopic = ConvertInt64(os.Getenv("DEATH_MSG_TOPIC"))
	ChatId = ConvertInt64(os.Getenv("CHAT_ID"))
	EventTopic = ConvertInt64(os.Getenv("EVENT_TOPIC"))
	LogFilePath = os.Getenv("LOG_FILE_PATH")
	AchievementTopic = ConvertInt64(os.Getenv("ACHIEVEMENT_TOPIC"))
	DefaultTopic = ConvertInt64(os.Getenv("DEFAULT_TOPIC"))
	OwnerId = ConvertInt64(os.Getenv("OWNER_ID"))
	return nil
}

func ConvertInt64(number string) int64 {
	value, err := strconv.ParseInt(number, 10, 64)
	if err != nil {
		return 0
	}
	return value
}
