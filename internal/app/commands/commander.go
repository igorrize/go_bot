package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)
	
type Commander struct {
	bot *tgbotapi.BotAPI
}

func NewComander(bot *tgbotapi.BotAPI) *Commander {
	return &Commander{
		bot,
		}
}