package commands

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type Commander struct {
	bot *tgbotapi.BotAPI
}

func NewComander(bot *tgbotapi.BotAPI) *Commander {
	return &Commander{
		bot,
	}
}

func (c *Commander) HandleCommand(message *tgbotapi.Message) {
	switch message.Command() {
	case "help":
		c.Help(message)
	case "search":
		c.Search(message)
	default:
		msg := tgbotapi.NewMessage(message.Chat.ID, "I don't know that command")
		_, err := c.bot.Send(msg)
		if err != nil {
			log.Println("Something went wrong: ", err)
		}
	}
}
