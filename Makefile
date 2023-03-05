build:
	go build ./cmd/bot
	chmod 644 ./bot
	zip bot.zip bot