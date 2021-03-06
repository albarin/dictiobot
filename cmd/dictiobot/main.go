package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/albarin/dictiobot/pkg/words"

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

	api := words.New(
		"https://wordsapiv1.p.rapidapi.com/words",
		os.Getenv("WORDSAPI_TOKEN"),
		http.DefaultClient,
	)

	bot.Handle(telebot.OnQuery, func(q *telebot.Query) {
		word := q.Text

		definitions, err := api.Word(word)
		if err != nil {
			log.Printf("error searching for word: %s", err)
		}

		results := make(telebot.Results, len(definitions))
		for i, def := range definitions {
			results[i] = createResult(word, def, i)
		}

		err = bot.Answer(q, &telebot.QueryResponse{Results: results})
		if err != nil {
			log.Println(err)
		}
	})

	bot.Start()
}

func createResult(word string, def words.Result, i int) *telebot.ArticleResult {
	result := &telebot.ArticleResult{
		Title:       word,
		Description: def.Definition,
	}

	result.SetContent(&telebot.InputTextMessageContent{
		Text:      def.Format(word),
		ParseMode: telebot.ModeMarkdownV2,
	})

	result.SetResultID(strconv.Itoa(i))
	return result
}
