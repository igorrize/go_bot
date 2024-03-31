package commands

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/igorrize/go_bot/internal/clients"
	"github.com/igorrize/go_bot/internal/storage"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func (c *Commander) Search (inputMessage *tgbotapi.Message) {
	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(3*time.Second))
	defer cancel()
	apiKey := os.Getenv("API_KEY")
	host := os.Getenv("HOST")
	udClient := ud_client.NewUDClient(apiKey, host)
	term := inputMessage.CommandArguments()
	definition, err := udClient.DefineTerm(ctx, term)
	if err != nil {
		log.Printf("error from client")
	}
	redisKey := strconv.FormatInt(inputMessage.Chat.ID, 10) + "_" + term
	err = storage.SetKey(redisKey, strings.Join(definition, "|"))
	if err != nil {
		log.Fatal(err)
	}
	maxPages := len(definition)
	msg := tgbotapi.NewMessage(inputMessage.Chat.ID, "Definition of your meme is - " + definition[0])
	buttonData := fmt.Sprintf("pager:next:1:%s:%d", redisKey, maxPages)
	nextButton := tgbotapi.NewInlineKeyboardButtonData("Next Definition", buttonData)
	keyboardRow := tgbotapi.NewInlineKeyboardRow(nextButton)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(keyboardRow)
	msg.ReplyMarkup = keyboard
    c.bot.Send(msg)
}