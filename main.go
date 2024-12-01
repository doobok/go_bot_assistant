package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	// Загружаем переменные окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки файла .env")
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	debug := os.Getenv("DEBUG")
	port := os.Getenv("PORT_NUMBER")

	fmt.Println("Токен:", token)
	fmt.Println("Режим DEBUG:", debug)

	// Создаём новый объект бота
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Настроим Webhook
	webhookURL := "https://5377-95-158-14-21.ngrok-free.app/webhook/"
	webhookConfig, err := tgbotapi.NewWebhook(webhookURL)
	if err != nil {
		log.Panic("Ошибка при создании webhook:", err)
	}

	_, err = bot.Request(webhookConfig)
	if err != nil {
		log.Panic("Ошибка при установке Webhook:", err)
	}

	// Обработка запросов на Webhook
	http.HandleFunc("/webhook/", func(w http.ResponseWriter, r *http.Request) {
		var update tgbotapi.Update
		// Используем json.NewDecoder для декодирования JSON в объект Update
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&update); err != nil {
			http.Error(w, "Ошибка при обработке запроса", http.StatusInternalServerError)
			return
		}

		// Обрабатываем обновления
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сообщение получено локально!")
			_, err := bot.Send(msg)
			if err != nil {
				log.Println("Ошибка при отправке сообщения:", err)
			}
		}
	})

	// Запускаем HTTP сервер для обработки запросов на Webhook
	log.Println("Listening on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
