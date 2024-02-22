package commands

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/igorrize/go_bot/internal/clients"
	"log"
	"os"
)

func (c *Commander) Search (inputMessage *tgbotapi.Message) {
	apiKey := os.Getenv("API_KEY")
	host := os.Getenv("HOST")
	log.Printf("starting search")
	udClient := ud_client.NewUDClient(apiKey, host)
	term := inputMessage.CommandArguments()
	definition, err := udClient.DefineTerm(term)
	if err != nil {
		log.Printf("error from client")
	}
	log.Printf("continue search")
	maxPages := len(definition) 
	msg := tgbotapi.NewMessage(inputMessage.Chat.ID, "Definition of your meme is - " + definition[0])
	buttonData := fmt.Sprintf("pager:next:1:%s:%d", definition, maxPages)
	log.Printf("buttondata||" + buttonData)
	nextButton := tgbotapi.NewInlineKeyboardButtonData("Next Definition", buttonData)
	keyboardRow := tgbotapi.NewInlineKeyboardRow(nextButton)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(keyboardRow)
	msg.ReplyMarkup = keyboard
    c.bot.Send(msg)
}

