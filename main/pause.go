package main

import (
	"log"
	"os"
	"strconv"

	"github.com/stanf0rd/utekai/database"
	"github.com/stanf0rd/utekai/generator"
	tBot "gopkg.in/tucnak/telebot.v2"
)

func pause(u database.User) {
	activePause, err := database.GetActivePauseByUserID(u.ID)
	if err != nil {
		log.Fatal(err)
	}

	if activePause != nil {
		activePause.Status = "failed"
		err := activePause.UpdateStatus()
		if err != nil {
			log.Fatal(err)
		}

		_, err = bot.Edit(activePause, texts["stop_request"])
		if err != nil {
			log.Println(err)
		}
	}

	questionID := chooseQuestion(u)
	if questionID == 0 {
		log.Printf("User #%d already asked max count of questions", u.ID)
		return
	}

	q, err := database.GetQuestionByID(questionID)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Question chosen: %v", q.Body)

	message := request(u)

	p := database.Pause{
		User:         u.ID,
		Question:     q.ID,
		AskMessageID: message.ID,
		ChatID:       message.Chat.ID,
	}

	if err := p.Save(); err != nil {
		log.Fatalf("Cannot save pause: %v", err)
	}

	printAllPauses()
}

func request(u database.User) *tBot.Message {
	message, err := bot.Send(u, texts["stop_request"], &tBot.ReplyMarkup{
		InlineKeyboard: [][]tBot.InlineButton{{readyButton}},
	})

	if err != nil {
		log.Fatalf("Unable to send pause request: %v", err)
	}

	return message
}

func ask(callback *tBot.Callback) {
	_, err := bot.Edit(callback.Message, callback.Message.Text)
	if err != nil {
		log.Fatalln("Cannot edit message")
	}

	u, err := database.GetUserByTelegramID(callback.Sender.ID)
	if err != nil {
		log.Fatalf("Cannot get user to ask a question: %v", err)
	}

	pause, err := database.GetActivePauseByUserID(u.ID)
	if err != nil {
		log.Fatalf("Cannot get active pause: %v", err)
	}
	if pause == nil {
		log.Fatalf("Cannot get active pause after request accept")
	}

	pause.Status = "asked"
	if err := pause.UpdateStatus(); err != nil {
		log.Fatalf("Cannot save pause: %v", err)
	}

	q, err := pause.GetQuestion()
	if err != nil {
		log.Fatalf("Cannot receive question by pause: %v", err)
	}

	bot.Send(pause, q.Body)
	printAllPauses()
}

func chooseQuestion(user database.User) int {
	pausesPerUser, _ := strconv.Atoi(os.Getenv("PAUSES_PER_USER"))

	askedQuestions, err := database.GetAskedQuestionsIDs(user.ID)
	askedQuestionsCount := len(askedQuestions)

	if askedQuestionsCount == pausesPerUser {
		return 0
	}

	var questions []int
	if askedQuestionsCount == 0 {
		questions, err = database.GetQuestionIDsByOrder("first")
	} else if askedQuestionsCount == pausesPerUser-1 {
		questions, err = database.GetQuestionIDsByOrder("last")
	} else {
		questions, err = database.GetQuestionIDsByOrder("random")
		if err == nil {
			questions = filterIntArray(questions, askedQuestions)
		}
	}

	if err != nil {
		log.Fatalln("Cannot get questions from database")
	}

	chosen := generator.GetRandomFromArray(questions)
	return chosen
}
