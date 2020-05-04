package main

import (
	"fmt"
	"log"
	"os"
	"time"

	tBot "gopkg.in/tucnak/telebot.v2"

	"github.com/stanf0rd/utekai/sheets"
)

var bot *tBot.Bot
var texts = map[string]string{
	"hello":                   "Привет! Мы клёвый бот, зовёмся паузой, Лиза классная, небо голубое.\n\n",
	"anonimous_ask":           "Хошь анонимности?",
	"anonimous_confirmed":     "Всё анонимно до чёртиков.",
	"non_anonimous_confirmed": "Вася, записывай.",
	"anonimous_notif":         "Ответы будут обезличены.",
	"non_anonimous_notif":     "Художница увидит ваш ник.",
	"notify_suggest":          "Включи уведомления, а то чо ты как лох. Чтобы бот и в могиле достучался.",
}

func main() {
	var err error
	bot, err = initBot()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Telegram-bot initialized")
	}

	sheets.Authorize()
	log.Println("Authorized to Google services")

	bot.Handle("/start", setupStartCases())

	log.Println("Starting bot...")
	bot.Start()
}

func initBot() (*tBot.Bot, error) {
	botToken := os.Getenv("BOT_TOKEN")
	apiURL := os.Getenv("TELEGRAM_API")

	return tBot.NewBot(tBot.Settings{
		Token:  botToken,
		URL:    apiURL,
		Poller: &tBot.LongPoller{Timeout: 10 * time.Second},
	})
}

func setupStartCases() func(*tBot.Message) {
	anonimous := tBot.InlineButton{
		Unique: "anonimous",
		Text:   "Yes",
	}
	nonAnonimous := tBot.InlineButton{
		Unique: "non_anonimous",
		Text:   "No",
	}

	bot.Handle(&anonimous, createAnonymityChoiceProcessor(true))
	bot.Handle(&nonAnonimous, createAnonymityChoiceProcessor(false))

	return func(message *tBot.Message) {
		log.Print(fmt.Printf("/hello: sender %v", message.Sender))

		bot.Send(message.Sender, texts["hello"]+texts["anonimous_ask"], &tBot.ReplyMarkup{
			InlineKeyboard: [][]tBot.InlineButton{{anonimous, nonAnonimous}},
		})

		go sheets.Check()
	}
}

func createAnonymityChoiceProcessor(anonimous bool) func(*tBot.Callback) {
	var notifText, confirmText string
	if anonimous {
		notifText = texts["anonimous_notif"]
		confirmText = texts["anonimous_confirmed"]
	} else {
		notifText = texts["non_anonimous_notif"]
		confirmText = texts["non_anonimous_confirmed"]
	}

	return func(callback *tBot.Callback) {
		defer bot.Respond(callback, &tBot.CallbackResponse{
			Text: notifText,
		})

		_, err := bot.Edit(callback.Message, texts["hello"]+confirmText)
		if err != nil {
			log.Fatalln("Cannot edit hello message", err)
		}
		_, err = bot.Send(callback.Sender, texts["notify_suggest"])
		if err != nil {
			log.Fatalln("Cannot send notify suggestion ", err)
		}
	}
}
