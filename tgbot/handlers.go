package tgbot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"MineCoreBot/config"
	"MineCoreBot/minecraft"
	"MineCoreBot/utils"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/hpcloud/tail"
)

func StartCommand(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b, "Hello! I'm MineCoreBot, here to help with your Minecraft server.", nil)
	return err
}

func executeCommand(b *gotgbot.Bot, ctx *ext.Context, cmd string) error {
	if ctx.EffectiveSender.Id() != config.OwnerId {
		return replyWithMessage(b, ctx, "You are not authorized to execute this command.")
	}
	if len(ctx.Args()) == 0 {
		return replyWithMessage(b, ctx, "You need to provide a command to execute.")
	}
	command := strings.Join(ctx.Args()[1:], " ")
	resp, err := minecraft.SendCommand(fmt.Sprintf("%s %s", cmd, command))
	if err != nil {
		return replyWithMessage(b, ctx, "An error occurred while executing the command.")
	}
	return replyWithResponse(b, ctx, resp)
}

func replyWithMessage(b *gotgbot.Bot, ctx *ext.Context, msg string) error {
	_, err := ctx.EffectiveMessage.Reply(b, msg, nil)
	return err
}

func replyWithResponse(b *gotgbot.Bot, ctx *ext.Context, resp string) error {
	_, err := ctx.EffectiveMessage.Reply(b, fmt.Sprintf("<b>Response from server:</b>\n<code>%s</code>", resp), &gotgbot.SendMessageOpts{ParseMode: "HTML"})
	return err
}

func ConsoleCommand(b *gotgbot.Bot, ctx *ext.Context) error {
	return executeCommand(b, ctx, "")
}

func WhiteListCommand(b *gotgbot.Bot, ctx *ext.Context) error {
	return executeCommand(b, ctx, "whitelist")
}

func BanCommand(b *gotgbot.Bot, ctx *ext.Context) error {
	return executeCommand(b, ctx, "kick")
}

func ReloadCommand(b *gotgbot.Bot, ctx *ext.Context) error {
	return executeCommand(b, ctx, "reload")
}

func StopCommand(b *gotgbot.Bot, ctx *ext.Context) error {
	return executeCommand(b, ctx, "stop")
}

func HelpCommand(b *gotgbot.Bot, ctx *ext.Context) error {
	helpMsg := `<b>Available commands:</b>
		<code>/start</code> - Start the bot
		<code>/console &lt;command&gt;</code> - Execute a command in the Minecraft server console
		<code>/whitelist &lt;command&gt;</code> - Execute a command in the Minecraft server whitelist
		<code>/ban &lt;command&gt;</code> - Execute a command in the Minecraft server ban
		<code>/reload &lt;command&gt;</code> - Execute a command in the Minecraft server reload
		<code>/stop &lt;command&gt;</code> - Execute a command in the Minecraft server stop
		<code>/help</code> - Show this help message`
	_, err := ctx.EffectiveMessage.Reply(b, helpMsg, &gotgbot.SendMessageOpts{ParseMode: "HTML"})
	return err
}

func SendLogs(msg string, topicid int64) {
	if _, err := tgbot.SendMessage(config.ChatId, msg, &gotgbot.SendMessageOpts{DisableNotification: true, MessageThreadId: topicid, ParseMode: "HTML"}); err != nil {
		log.Println("Error sending message to chat:", err)
	}
}

var (
	lastDefaultMsgId  int64
	currentDefaultMsg string
)

func HandleLogs() {
	t, err := tail.TailFile(config.LogFilePath, tail.Config{Follow: true, ReOpen: true})
	if err != nil {
		log.Println("Error opening log file:", err)
		return
	}

	for line := range t.Lines {
		if line.Err != nil {
			log.Println("Error reading log file:", line.Err)
			continue
		}
		processLogLine(line.Text)
	}
}

func processLogLine(logLine string) {
	username, userMessage, logType := utils.ClassifyMinecraftLog(logLine)
	switch logType {
	case "User Sent Message":
		SendLogs(fmt.Sprintf("<b>%s</b>: %s", username, userMessage), config.ChatTopic)
	case "Join Event":
		SendLogs(fmt.Sprintf("<b>%s</b> joined the game", username), config.EventTopic)
	case "Left Event":
		SendLogs(fmt.Sprintf("<b>%s</b> left the game", username), config.EventTopic)
	case "Death Message":
		SendLogs(fmt.Sprintf("<b>%s</b> died", username), config.DeathMsgTopic)
	case "Achievement":
		SendLogs(fmt.Sprintf("<b>%s</b> has made the advancement", username), config.AchievementTopic)
	default:
		sendOrEditDefaultLog(logLine)
	}
}

func sendOrEditDefaultLog(newLog string) {
	formattedNewLog := fmt.Sprintf("<code>%s</code>", newLog)

	if len(currentDefaultMsg)+len(formattedNewLog) > 3000 || lastDefaultMsgId == 0 {
		sendDefaultLogMessage(formattedNewLog)
	} else {
		updateDefaultLogMessage(formattedNewLog)
	}

}

func updateDefaultLogMessage(formattedNewLog string) {
	currentDefaultMsg += "\n" + formattedNewLog
	attempts := 0
	const maxAttempts = 5

	for attempts < maxAttempts {
		if _, _, err := tgbot.EditMessageText(currentDefaultMsg, &gotgbot.EditMessageTextOpts{
			ParseMode: "HTML",
			MessageId: lastDefaultMsgId,
			ChatId:    config.ChatId,
		}); err != nil {
			if te, ok := err.(*gotgbot.TelegramError); ok {
				if te.Code == 429 {
					retryAfter := te.ResponseParams.RetryAfter
					if retryAfter > 0 {
						log.Printf("Rate limit exceeded. Retrying after %d seconds...\n", retryAfter)
						time.Sleep(time.Duration(retryAfter) * time.Second)
					} else {
						attempts++
						time.Sleep(2 * time.Second)
					}
					continue
				}
			}
			log.Println("Error editing message in default topic:", err)
			break
		}
		break
	}
}

func sendDefaultLogMessage(formattedNewLog string) {
	msg, err := tgbot.SendMessage(config.ChatId, formattedNewLog, &gotgbot.SendMessageOpts{
		DisableNotification: true,
		MessageThreadId:     config.DefaultTopic,
		ParseMode:           "HTML",
	})
	if err != nil {
		log.Println("Error sending new message to default topic:", err)
	} else {
		lastDefaultMsgId = msg.MessageId
		currentDefaultMsg = formattedNewLog
	}
}

func HandleChatTopic(b *gotgbot.Bot, ctx *ext.Context) error {
	chatId := ctx.EffectiveChat.Id
	topicId := ctx.EffectiveMessage.MessageThreadId
	if chatId != config.ChatId || topicId != config.ChatTopic {
		return nil
	}
	senderName := ctx.EffectiveMessage.From.FirstName
	resp, err := minecraft.SendTgToMc(ctx.EffectiveMessage.Text, senderName)
	if err != nil {
		return replyWithMessage(b, ctx, "An error occurred while sending the message to the server.")
	}
	if resp != "" {
		return replyWithMessage(b, ctx, resp)
	}
	return nil
}
