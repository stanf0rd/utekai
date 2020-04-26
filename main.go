package main

import (
	"fmt"
	"log"
	"os"
	"time"

	tBot "gopkg.in/tucnak/telebot.v2"

	"github.com/stanf0rd/utekai/sheets"
)

func main() {
	fmt.Println("BOT_TOKEN:", os.Getenv("BOT_TOKEN"))

	botToken := os.Getenv("BOT_TOKEN")
	apiURL := os.Getenv("TELEGRAM_API")

	b, err := tBot.NewBot(tBot.Settings{
		Token:  botToken,
		URL:    apiURL,
		Poller: &tBot.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	sheets.Authorize()

	b.Handle("/hello", func(m *tBot.Message) {
		b.Send(m.Sender, "hello world")
		log.Print(fmt.Printf("/hello: sender %v", m.Sender))

		sheets.Check()
	})

	b.Start()
}
