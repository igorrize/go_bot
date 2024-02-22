package commands
import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/igorrize/go_bot/internal/clients"
	"log"
	"os"
)

func  (c *Commander) Short (messageId *int, chatId int64, defs []string) {
	apiKey := os.Getenv("API_KEY")
	host := os.Getenv("HOST")
	log.Printf("starting shortener")
	shortClient := ud_client.NewShortClient(apiKey, host)
	summary, err := shortClient.ShortDefinitions(defs)
	if err != nil {
		log.Printf("error from client")
	}
	log.Printf("continue search")
	
	msg := tgbotapi.NewEditMessageText(chatId, *messageId, "TLDR - " + summary)
	c.bot.Send(msg)
}