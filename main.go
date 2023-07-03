package main

import (
	"log"

	"strconv"

	"os"

	godotenv "github.com/joho/godotenv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const greetingsMessage = "Добрый день, для получения списка файлов /list, для чтения /read, /kroki <diagram> для формирования диаграммы"
const noSuchCommanText = "Нет такой команды, для получения списка команд используйте /start"

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No env file foiund")
	}
}

func main() {
	debug, existsDebug := os.LookupEnv("DEBUG")
	botToken, existToken := os.LookupEnv("BOT_TOKEN")
	folder, existFolder := os.LookupEnv("BOT_DIR")

	if !existToken || !existFolder || !existsDebug {
		log.Panic("Not bot token found!")
	}
	log.Print(folder)

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug, _ = strconv.ParseBool(debug)

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			switch update.Message.Command() {
			case "kroki":
				text := getMessageTextOnly(update.Message.Text, update.Message.Command())
				msg := tgbotapi.NewPhoto(update.Message.Chat.ID, getKrokiMedia(text))
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			default:
				msg := formReplyMessage(update.Message.Chat.ID, update, folder)
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			}
		}
	}
}

func formReplyMessage(chatId int64, update tgbotapi.Update, folder string) tgbotapi.MessageConfig {
	message := update.Message

	switch message.Command() {
	case "start":
		return tgbotapi.NewMessage(chatId, greetingsMessage)
	case "list":
		return tgbotapi.NewMessage(chatId, getFileList(folder))
	case "read":
		text := getMessageTextOnly(message.Text, message.Command())
		return tgbotapi.NewMessage(chatId, readFile(text, folder))
	default:
		return tgbotapi.NewMessage(chatId, noSuchCommanText)
	}
}
