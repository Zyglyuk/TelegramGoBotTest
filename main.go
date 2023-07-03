package main

import (
	"log"

	"strings"

	"strconv"

	"os"

	godotenv "github.com/joho/godotenv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const greetingsMessage = "Добрый день, для получения списка файлов /list, для чтения /read"
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
			msg := formReplyMessage(update.Message.Chat.ID, update, folder)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
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
		text, _ := strings.CutPrefix(message.Text, "/read")
		return tgbotapi.NewMessage(chatId, readFile(text, folder))
	default:
		return tgbotapi.NewMessage(chatId, noSuchCommanText)
	}
}

func getFileList(folder string) string {
	files, err := os.ReadDir(folder)
	if err != nil {
		log.Fatal(err)
	}
	message := ""
	for _, file := range files {
		message = message + file.Name() + ";  "
	}
	return message
}

func readFile(text string, folder string) string {
	file := strings.Trim(text, " ")
	if strings.HasPrefix(file, "@") {
		file = strings.Split(file, " ")[1]
	}
	fileContent, err := os.ReadFile(folder + "/" + file)
	if err != nil {
		return err.Error()
	}
	return string(fileContent)
}
