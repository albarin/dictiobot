include .env
export

run:
	PORT=$(PORT) \
	WEBHOOK_URL=$(WEBHOOK_URL) \
	BOT_TOKEN=$(BOT_TOKEN) \
	WORDSAPI_TOKEN=$(WORDSAPI_TOKEN) \
	go run cmd/dictiobot/main.go