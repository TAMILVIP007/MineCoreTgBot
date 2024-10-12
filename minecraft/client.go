package minecraft

import (
	"MineCoreBot/config"
	"fmt"
	"html"
	"strings"

	"github.com/gorcon/rcon"
)

func InitClient() error {
	_, err := SendCommand("list")
	if err != nil {
		return err
	}
	return nil
}

func SendCommand(command string) (string, error) {
	commanD := strings.TrimSpace(command)
	conn, err := rcon.Dial(config.RconnHost+":"+config.RconnPort, config.RconnPassword)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	resp, err := conn.Execute(commanD)
	if err != nil {
		return "", err
	}
	return html.EscapeString(resp), nil
}

func SendTgToMc(message string, name string) (string, error) {
	command := fmt.Sprintf(
		"tellraw @a [{\"text\":\"[FromTg: \",\"color\":\"blue\"},{\"text\":\"%s\",\"color\":\"red\"},{\"text\":\"]\\n%s\",\"color\":\"white\"}]",
		name, message,
	)
	resp, err := SendCommand(command)
	if err != nil {
		return "", err
	}
	return resp, nil
}
