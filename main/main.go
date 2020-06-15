package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tBot "gopkg.in/tucnak/telebot.v2"

	"github.com/stanf0rd/utekai/database"
	"github.com/stanf0rd/utekai/generator"
	"github.com/stanf0rd/utekai/sheets"
)

var bot *tBot.Bot
var texts = map[string]string{
	"anonimous_ask":           "Хотите ли вы отправлять ответы анонимно? ",
	"greeting_text":           "Бот будет просить вас остановиться несколько раз в день в течение недели. После активации сообщения у вас есть 30 секунд до вопроса. Остановитесь, прислушайтесь к себе и миру.\n\nДалее приходит вопрос. Ответьте на него одним сообщением - текстом или фото с подписью. Пожалуйста, не торопитесь, дайте себе время на погружение. \n\nВаши ответы будут использоваться для создания произведения.",
	"notify_suggest":          "Для большего вовлечения в процесс включите уведомления.",
	"anonimous_confirmed":     "Ваши ответы будут обезличены.",
	"non_anonimous_confirmed": "Художница будет видеть, кому принадлежат ответы.",
	"stop_request":            "Остановись",
	"stopped_button":          "Останавливаюсь (30c)",
	"after_answer":            "Благодарю",
	"admin_hello":             "Здраствуй, Нео",
}
var readyButton = createButton(texts["stopped_button"])
var admins = strings.Split(os.Getenv("ADMIN_TELEGRAM_IDS"), ",")

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

	bot.Handle("/what", what)
	bot.Handle("/who", who)
	bot.Handle("/how", how)

	bot.Handle("/broadcast", broadcast)
	bot.Handle(tBot.OnText, getAnswer)
	bot.Handle(tBot.OnPhoto, getAnswer)
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
			isAdmin := stringInSlice(strconv.Itoa(callback.Sender.ID), admins)

			u := database.User{
				TelegramID: callback.Sender.ID,
				Anonymous:  anonymous,
				Admin:      isAdmin,
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
					if err := sheets.AddUserToSheet(u); err != nil {
						log.Printf("Cannot update user list: %v", err)
					}
				} else {
					log.Fatalf("Cannot save user %d, error: %v", u.TelegramID, err)
				}
			}

			if isAdmin {
				bot.Send(u, texts["admin_hello"])
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
	nonAnonymous := createButton("Нет")
	anonymous := createButton("Да")

	bot.Handle(&anonymous, createRegistrar(true))
	bot.Handle(&nonAnonymous, createRegistrar(false))

	return func(message *tBot.Message) {
		bot.Send(message.Sender, texts["anonimous_ask"], &tBot.ReplyMarkup{
			InlineKeyboard: [][]tBot.InlineButton{{anonymous, nonAnonymous}},
		})
	}
}

func info(message *tBot.Message) {
	bot.Send(message.Sender, `

	`, tBot.NoPreview)
}

func what(message *tBot.Message) {
	bot.Send(message.Sender, `
Этот бот – часть серии перформансов художницы Перебатовой Елизаветы. Пауза определяется авторкой как встреча с собой и момент максимального сосредоточения внимания. Настолько конкретного, что остальные процессы становятся почти невозможными. Направляя внимание остановкой и вопросами, человек способен пойти по непривычному сценарию и что-то обнаружить в себе. Вокруг себя. В себе через то, как ощущается всё вокруг.

Работа с остановкой содержит отсылки к суфийской философии и упражнениям, согласно которым человеческое сознание может преодолеть свои рамки именно в момент прерывания действия. Бот предлагает человеку остановиться и ответить на вопрос, а затем собирает данные для следующего произведения. Тем самым он создает базу ответов, через которую зритель может войти в контакт с собой через чужой опыт и соотношение с ним.

Позже здесь появится ссылка на произведение, созданное на основе ответов.
	`, tBot.NoPreview)
}

func who(message *tBot.Message) {
	bot.Send(message.Sender, `
Контакты художницы для обратной связи:
@Lisavetata
https://www.facebook.com/elisaveta.perebatova
lisavetalis@yandex.ru

Разработка бота:
@dpaliy
https://github.com/stanf0rd
	`, tBot.NoPreview)
}

func how(message *tBot.Message) {
	bot.Send(message.Sender, texts["greeting_text"], tBot.NoPreview)
}
