include .env
export

run:
	PORT=$(PORT) \
	WEBHOOK_URL=$(WEBHOOK_URL) \
	BOT_TOKEN=$(BOT_TOKEN) \
	go run cmd/dictiobot/main.go