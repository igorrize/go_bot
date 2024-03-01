package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/igorrize/go_bot/internal/app/commands"
	"github.com/igorrize/go_bot/internal/services"
	"github.com/igorrize/go_bot/internal/storage"
	"log"
	"os"
)

func main() {

	token := os.Getenv("TELEGA_TOKEN")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.UpdateConfig{
		Timeout: 60,
		}

		updates := bot.GetUpdatesChan(u)
		if err != nil {
			log.Panic(err)
		}

		redisClient := storage.InitRedisClient()

		defer storage.CloseRedisClient(redisClient)

		commander := commands.NewComander(bot)

		for update := range updates {
			if update.CallbackQuery != nil {
				services.CallbackQueryHandler(update.CallbackQuery, bot, update.CallbackQuery.Message.Chat.ID)
				continue
			} else if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "help":
					commander.Help(update.Message)
					case "search":
						commander.Search(update.Message)
						default:
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know that command")
							_, err := bot.Send(msg)
							if err != nil {
								log.Fatal("Something went wrong")
							}
				}
			}
		}


}

