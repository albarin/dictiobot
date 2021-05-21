package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/albarin/dictiobot/pkg/words"

	"github.com/newrelic/go-agent/v3/newrelic"
	"gopkg.in/tucnak/telebot.v2"
)

func main() {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("dictiobot"),
		newrelic.ConfigLicense(os.Getenv("NEW_RELIC_TOKEN")),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	if err != nil {
		log.Fatalf("error initializing newrelic: %s", err)
		return
	}

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
		app,
	)

	bot.Handle(telebot.OnQuery, onQueryHandler(bot, api, app))

	bot.Start()

	app.Shutdown(10 * time.Second)
}

func onQueryHandler(bot *telebot.Bot, api *words.API, app *newrelic.Application) func(q *telebot.Query) {
	return func(q *telebot.Query) {
		word := q.Text
		app.RecordCustomEvent("onQuery", map[string]interface{}{
			"word": word,
		})

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
	}
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
