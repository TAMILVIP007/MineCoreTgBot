package utils

import (
	"regexp"
	"strings"
)

func ClassifyMinecraftLog(logLine string) (string, string, string) {
	logLine = strings.TrimSpace(logLine)
	var username, userMessage string
	re := regexp.MustCompile(`\[\d{2}:\d{2}:\d{2}\] \[Server thread/INFO\]: \[Not Secure\] <([^>]+)> (.+)`)
	if matches := re.FindStringSubmatch(logLine); matches != nil {
		username = matches[1]
		userMessage = matches[2]
	}

	if userMessage != "" {
		return username, userMessage, "User Sent Message"
	}
	if strings.Contains(logLine, "joined the game") {
		return extractUsername(logLine), "", "Join Event"
	}
	if strings.Contains(logLine, "left the game") {
		return extractUsername(logLine), "", "Left Event"
	}
	if strings.Contains(logLine, "was slain by") || strings.Contains(logLine, "died:") {
		return extractUsername(logLine), "", "Death Message"
	}
	if strings.Contains(logLine, "has made the advancement") {
		return extractUsername(logLine), "", "Achievement"
	}
	return extractUsername(logLine), "", "Other"
}

func extractUsername(logLine string) string {
	re := regexp.MustCompile(`\[\d{2}:\d{2}:\d{2}\] \[Server thread/INFO\]: ([^ ]+)`)
	if matches := re.FindStringSubmatch(logLine); matches != nil {
		return matches[1]
	}
	return ""
}