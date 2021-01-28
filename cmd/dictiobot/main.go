package main

import (
	"log"
	"os"

	"gopkg.in/tucnak/telebot.v2"
)

func main() {
	webhook := &telebot.Webhook{
		Listen: ":" + os.Getenv("PORT"),
		Endpoint: &telebot.WebhookEndpoint{
			PublicURL: os.Getenv("WEBHOOK_URL") + "/bot" + os.Getenv("BOT_TOKEN"),
		},
	}

	settings := telebot.Settings{
		Token:  os.Getenv("BOT_TOKEN"),
		Poller: webhook,
	}

	bot, err := telebot.NewBot(settings)
	if err != nil {
		log.Fatalf("error initializing bot: %s", err)
		return
	}

	bot.Handle(telebot.OnText, func(msg *telebot.Message) {
		_, err := bot.Send(msg.Sender, "I'm a 🤖!")
		if err != nil {
			log.Printf("error sending message: %s", err)
			return
		}
	})

	bot.Start()
}
