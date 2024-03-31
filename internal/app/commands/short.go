package commands
import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/igorrize/go_bot/internal/clients"
	"github.com/igorrize/go_bot/internal/storage"
	"log"
	"os"
	"strings"
	"sync"
)

func  (c *Commander) Short (messageId *int, chatId int64, redisKey string) {
	apiKey := os.Getenv("API_KEY")
	host := os.Getenv("HOST")
	log.Printf("starting shortener")
	shortClient := ud_client.NewShortClient(apiKey, host)
	summary, err := shortClient.ShortDefinitions(grabData(redisKey))

	log.Printf(summary)

	if err != nil {
		log.Printf("error from client")
	}
	log.Printf("continue search")
	
	msg := tgbotapi.NewEditMessageText(chatId, *messageId, "TLDR - " + summary)
	c.bot.Send(msg)
}

type dataForShortnes struct {
	UDdata []string
	WordData string
}

func grabData(redisKey string) (text string) {
	var waitgroup sync.WaitGroup
	waitgroup.Add(2)
	data := dataForShortnes{}

	go func() {
		getRedisData(&data, redisKey)
		waitgroup.Done()
	}()
	go func() {
		getWordData(&data, redisKey)
		waitgroup.Done()
	}()
	return strings.Join(data.UDdata, " ") + data.WordData
}

func getRedisData(dataForShortnes *dataForShortnes, redisKey string) {
	redisData := storage.GetKey(redisKey)
	dataForShortnes.UDdata = strings.Split(redisData,"|")
}

type WordApiResponse struct {
	Word        string `json:"word"`
	Definitions []struct {
		Definition   string `json:"definition"`
		PartOfSpeech string `json:"partOfSpeech"`
	} `json:"definitions"`
}

func getWordData(dataForShortnes *dataForShortnes, redisKey string) {
	wordApiKey := "3ba6283321msh80a967930e01301p13ea9fjsna2afafa73bdc"
	wordHost := "wordsapiv1.p.rapidapi.com"
	wordClient := ud_client.NewWordApiClient(wordApiKey, wordHost)
	term := strings.Split(redisKey,"_")[0]

	definition := wordClient.WordDefinition(term)
	dataForShortnes.WordData = definition
}