package services

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/igorrize/go_bot/internal/app/commands"
	"github.com/igorrize/go_bot/internal/storage"
	"log"
	"strconv"
	"strings"
)

func CallbackQueryHandler(query *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI, chatId int64) {
	split := strings.Split(query.Data, ":")
	if split[0] == "pager" {
		HandleNavigationCallbackQuery(query.Message.MessageID, bot, chatId, split[1:]...)
		return
	}
}

func HandleNavigationCallbackQuery(messageId int, bot *tgbotapi.BotAPI, chatId int64, data ...string) {
	pagerType := data[0]
	maxPages, _ := strconv.Atoi(data[3])
	currentPage, _ := strconv.Atoi(data[1])
	itemsPerPage := 1
	log.Printf("max pages1"+strconv.Itoa(maxPages))

	redisData := storage.GetKey(data[2])
	searchData := strings.Split(redisData,"|")
	if pagerType == "next" {
		nextPage := currentPage + 1
		if nextPage < maxPages {
			SendSearchData(searchData, nextPage, maxPages, itemsPerPage, &messageId, chatId, bot, data[2])
		}
	}
	if pagerType == "prev" {
		previousPage := currentPage - 1
		if previousPage >= 0 {
			SendSearchData(searchData, previousPage, maxPages, itemsPerPage, &messageId, chatId, bot, data[2])
		}
	}
	if pagerType == "tldr" {
		commander := commands.NewComander(bot)
		commander.Short(&messageId, chatId, data[4:])
	}
}

func SendSearchData(data []string, currentPage, maxPages, count int, messageId *int, chatId int64, bot *tgbotapi.BotAPI, redisKey string) {

	text, keyboard := SearchDataTextMarkup(data, currentPage, count, maxPages, redisKey)

    var cfg tgbotapi.Chattable
    if messageId == nil {
        msg := tgbotapi.NewMessage(chatId, text)
        msg.ReplyMarkup = keyboard
        cfg = msg
    } else {
        msg := tgbotapi.NewEditMessageText(chatId, *messageId, text)
        msg.ReplyMarkup = &keyboard
        cfg = msg
    }

	bot.Send(cfg)
}

func SearchDataTextMarkup(data []string, currentPage, count, maxPages int, redisKey string) (text string, markup tgbotapi.InlineKeyboardMarkup) {
	text = data[currentPage]
	var rows []tgbotapi.InlineKeyboardButton

	if currentPage > 0 {
	buttonData := fmt.Sprintf("pager:prev:%d:%s:%d", currentPage-1, redisKey, maxPages)
	rows = append(rows, tgbotapi.NewInlineKeyboardButtonData("Prev Definition", buttonData))
	}

	if currentPage < maxPages-1 {
		buttonData := fmt.Sprintf("pager:next:%d:%s:%d", currentPage+1, redisKey, maxPages)
		rows = append(rows, tgbotapi.NewInlineKeyboardButtonData("Next Definition", buttonData))
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardButtonData("tldr", fmt.Sprintf("pager:next:%d:%d", currentPage, count)))
	markup = tgbotapi.NewInlineKeyboardMarkup(rows)
	return
}
