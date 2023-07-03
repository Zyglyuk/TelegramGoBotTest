package main

import (
	"io"
	"net/http"
	"os"

	"log"

	"strings"

	bytes "bytes"
	zlib "compress/zlib"
	base64 "encoding/base64"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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

func readFile(file string, folder string) string {
	fileContent, err := os.ReadFile(folder + "/" + file)
	if err != nil {
		return err.Error()
	}
	return string(fileContent)
}

func getMessageTextOnly(text string, command string) string {
	trimmedText, _ := strings.CutPrefix(text, "/"+command)
	trimmedText = strings.Trim(trimmedText, " ")
	if strings.HasPrefix(trimmedText, "@") {
		trimmedText, _ = strings.CutPrefix(trimmedText, "@ZyglyukGoBot")
		trimmedText = strings.Trim(trimmedText, " ")
	}
	return trimmedText
}

func getKrokiMedia(text string) tgbotapi.RequestFileData {
	diagramType := strings.Split(text, " ")[0]
	enceodedDiagram, _ := encode(text)

	url := "https://kroki.io/" + diagramType + "/png/" + enceodedDiagram

	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
		log.Println("Error retrieving data")
	}

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	return tgbotapi.FileBytes{
		Name:  "diagram.svg",
		Bytes: body,
	}
}

func encode(input string) (string, error) {
	var buffer bytes.Buffer
	writer, err := zlib.NewWriterLevel(&buffer, 9)
	if err != nil {
		log.Println("failde to create a writer")
		return "", err
	}
	_, err = writer.Write([]byte(input))
	writer.Close()
	if err != nil {
		log.Println("fail to create the payload")
		return "", err
	}
	result := base64.URLEncoding.EncodeToString(buffer.Bytes())
	return result, nil
}
