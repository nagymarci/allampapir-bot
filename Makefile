build:
	go build ./...
	chmod 644 bot
	zip bot.zip bot