package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

const (
	repllyInPrivateChannel = true
)

func getEnvVariable(key string) string {
	// загружаем .env файл
	err := godotenv.Load("key.env")

	if err != nil {
		log.Fatal("Ошибка при загрузке.env файла")
	}

	return os.Getenv(key)
}

func main() {
	bot, err := tgbotapi.NewBotAPI(getEnvVariable("key"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Авторизован: %s", bot.Self.UserName)

	// установка long-polling request
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // проверяем наличие нового сообщения
			text := update.Message.Text      // текст сообщения
			chatID := update.Message.Chat.ID // ID чата
			userID := update.Message.From.ID // ID пользователя
			var replyMessage string

			log.Printf("[%s][%d] %s", update.Message.From.UserName, userID, text)

			// проверяем текст сообщения и запсиываем в переменную
			switch {
			case text == "Тест": // запрос пользователя
				replyMessage = TestMessage
				break
			case text == "Список": // запрос пользователя
				replyMessage = ListMessage
				break
			}

			if len(replyMessage) > 0 {
				if repllyInPrivateChannel {
					chatID = userID
					msg := tgbotapi.NewMessage(chatID, replyMessage)

					bot.Send(msg)
				} else {
					// отправляем ответ
					msg := tgbotapi.NewMessage(chatID, replyMessage) // создаём новое сообщение
					msg.ReplyToMessageID = update.Message.MessageID  // указываем сообщение, на которое нужно ответить

					bot.Send(msg)
				}
			}
		}
	}
}
