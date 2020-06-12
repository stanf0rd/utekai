package main

import (
	"fmt"
	"log"
	"os"
	"time"

	tBot "gopkg.in/tucnak/telebot.v2"

	"github.com/stanf0rd/utekai/database"
	"github.com/stanf0rd/utekai/generator"
	"github.com/stanf0rd/utekai/sheets"
)

var bot *tBot.Bot
var texts = map[string]string{
	"anonimous_ask":           "Оставить художнице возможность связаться с вами?",
	"greeting_text":           "Бот будет отправлять что-то иногда. Отвечай на вопросы. Можно фото, текст, чо попросят.",
	"notify_suggest":          "Для большего вовлечения в процесс включи уведомления.",
	"anonimous_confirmed":     "Ваши ответы будут обезличены.",
	"non_anonimous_confirmed": "Художница увидит ваш ник и получит возможность связаться с вами.",
	"stop_request":            "Остановись, чувачелло",
	"stopped_button":          "Остановился",
}
var readyButton = createButton(texts["stopped_button"])

func main() {
	var err error
	bot, err = initBot()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Telegram-bot initialized")
	}

	printAllUsers()
	printAllPauses()
	printAllQuestions()

	bot.Handle("/start", createStart())
	bot.Handle("/pause", func(message *tBot.Message) {
		user, err := database.GetUserByTelegramID(message.Sender.ID)
		if err != nil {
			log.Fatal(err)
		}
		pause(*user)
	})
	bot.Handle("/print", func(message *tBot.Message) {
		// printAllUsers()
		// printAllQuestions()
	})
	bot.Handle(&readyButton, ask)

	log.Println("Starting bot...")
	bot.Start()
}

func createButton(text string) tBot.InlineButton {
	return tBot.InlineButton{
		Unique: generator.String(10),
		Text:   text,
	}
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

func createStart() func(*tBot.Message) {
	createRegistrar := func(anonymous bool) func(*tBot.Callback) {
		text := fmt.Sprintf("%s\n\n%s", texts["greeting_text"], texts["notify_suggest"])
		if anonymous {
			text = fmt.Sprintf("%s\n\n%s", text, texts["anonimous_confirmed"])
		} else {
			text = fmt.Sprintf("%s\n\n%s", text, texts["non_anonimous_confirmed"])
		}

		register := func(callback *tBot.Callback) {
			u := database.User{
				TelegramID: callback.Sender.ID,
				Anonymous:  anonymous,
			}

			exists, err := u.Exists()
			if err != nil {
				log.Fatalf("Cannot check user %d existance, error: %v", u.TelegramID, err)
			}

			if exists {
				err := u.UpdateAnonymity()

				if err == nil {
					log.Printf("User #%d set his anonymity to %t", u.ID, u.Anonymous)
				} else {
					log.Fatalf("Cannot update user %d, error: %v", u.TelegramID, err)
				}
			} else {
				err := u.Save()

				if err == nil {
					log.Printf("User %d was saved, his ID was set to %d", u.TelegramID, u.ID)
					sheets.AddUserToSheet(u)
				} else {
					log.Fatalf("Cannot save user %d, error: %v", u.TelegramID, err)
				}
			}

			_, err = bot.Edit(callback.Message, text)
			if err != nil {
				log.Fatalln("Cannot edit message")
			}
		}

		return register
	}

	return createAnonymityAsk(createRegistrar)
}

func createAnonymityAsk(
	createRegistrar func(anonimous bool) func(*tBot.Callback),
) func(*tBot.Message) {
	nonAnonymous := createButton("Да")
	anonymous := createButton("Нет")

	bot.Handle(&anonymous, createRegistrar(true))
	bot.Handle(&nonAnonymous, createRegistrar(false))

	return func(message *tBot.Message) {
		bot.Send(message.Sender, texts["anonimous_ask"], &tBot.ReplyMarkup{
			InlineKeyboard: [][]tBot.InlineButton{{nonAnonymous, anonymous}},
		})
	}
}
