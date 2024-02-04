package services

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	maxPages, _ := strconv.Atoi(data[1])
	currentPage, _ := strconv.Atoi(data[2])
	itemsPerPage, _ := strconv.Atoi(data[3])
	if pagerType == "next" {
		nextPage := currentPage + 1
		if nextPage < maxPages {
			SendSearchData(data[4:], nextPage, itemsPerPage, maxPages, &messageId, chatId, bot)
		}
	}
	if pagerType == "prev" {
		previousPage := currentPage - 1
		if previousPage >= 0 {
			SendSearchData(data[4:], previousPage, itemsPerPage, maxPages, &messageId, chatId, bot)
		}
	}
}

func SendSearchData(data []string, currentPage, maxPages, count int, messageId *int, chatId int64, bot *tgbotapi.BotAPI) {

	text, keyboard := SearchDataTextMarkup(data, currentPage, count, maxPages)

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

func SearchDataTextMarkup(data []string, currentPage, count, maxPages  int) (text string, markup tgbotapi.InlineKeyboardMarkup) {
	text = strings.Join(data[currentPage*count:currentPage*count+count], "\n")

	var rows []tgbotapi.InlineKeyboardButton
	if currentPage > 0 {
		rows = append(rows, tgbotapi.NewInlineKeyboardButtonData("Previous", fmt.Sprintf("pager:prev:%d:%d", currentPage, count)))
	}
	if currentPage < maxPages-1 {
		rows = append(rows, tgbotapi.NewInlineKeyboardButtonData("Next", fmt.Sprintf("pager:next:%d:%d", currentPage, count)))
	}

	markup = tgbotapi.NewInlineKeyboardMarkup(rows)
	return
}
