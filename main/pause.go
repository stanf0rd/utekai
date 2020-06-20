package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/stanf0rd/utekai/database"
	"github.com/stanf0rd/utekai/generator"
	tBot "gopkg.in/tucnak/telebot.v2"
)

func broadcast(message *tBot.Message) {
	user, err := database.GetUserByTelegramID(message.Sender.ID)
	if err != nil {
		log.Printf("Cannot get user by telegramID %v", err)
		return
	}

	if !user.Admin {
		log.Printf(
			"Non-admin user #%d tried to access broadcast feature", message.Sender.ID,
		)
		return
	}

	users, err := database.GetAllUsers()
	if err != nil {
		log.Fatalf("Cannot get users from DB: %v", err)
	}

	for _, user := range users {
		pause(user)
		time.Sleep(10 * time.Second)
	}
}

func pause(u database.User) {
	activePause, err := database.GetActivePauseByUserID(u.ID)
	if err != nil {
		log.Fatal(err)
	}

	if activePause != nil {
		err := activePause.UpdateStatus("failed")
		if err != nil {
			log.Fatal(err)
		}

		err = bot.Delete(activePause)
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

	message, err := request(u)
	if err != nil {
		log.Printf("Unable to send pause request: %v", err)
		return
	}

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

func request(u database.User) (*tBot.Message, error) {
	message, err := bot.Send(u, texts["stop_request"], &tBot.ReplyMarkup{
		InlineKeyboard: [][]tBot.InlineButton{{readyButton}},
	})

	if err != nil {
		return nil, fmt.Errorf("Unable to send pause request: %v", err)
	}

	return message, nil
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

	q, err := pause.GetQuestion()
	if err != nil {
		log.Fatalf("Cannot receive question by pause: %v", err)
	}

	time.Sleep(30 * time.Second)
	bot.Send(pause, q.Body)
	if err := pause.UpdateStatus("asked"); err != nil {
		log.Fatalf("Cannot update pause status: %v", err)
	}
	printAllPauses()
}

func getAnswer(message *tBot.Message) {
	u, err := database.GetUserByTelegramID(message.Sender.ID)
	if err != nil {
		log.Printf("User written a message not found: %v", err)
		return
	}

	pause, err := database.GetActivePauseByUserID(u.ID)
	if err != nil {
		log.Printf("Cannot get active pause: %v", err)
		return
	}
	if pause == nil {
		log.Printf("Cannot get active accepted pause to write answer")
		return
	}
	if pause.Status != "asked" {
		log.Printf("Pause is not accepted, cannot write answer")
		return
	}

	if err := pause.UpdateStatus("done"); err != nil {
		log.Printf("Cannot update pause status: %v", err)
		return
	}

	var answer string
	if message.Photo != nil {
		answer = message.Caption
		filePath, err := getPhoto(message.Photo)
		if err != nil {
			log.Printf("Cannot get photo: %v", err)
		} else if err := pause.AddPhoto(filePath); err != nil {
			log.Printf("Cannot add photo to pause: %v", err)
		}
	} else {
		answer = message.Text
	}

	if err := pause.AddAnswer(answer); err != nil {
		log.Printf("Cannot save pause: %v", err)
		return
	}

	bot.Send(pause, texts["after_answer"])

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

func getPhoto(photo *tBot.Photo) (filePath string, err error) {
	filePath = fmt.Sprintf("/images/%s.jpeg", generator.String(10))

	readCloser, err := bot.GetFile(photo.MediaFile())
	if err != nil {
		return "", fmt.Errorf("Cannot get file reader from bot: %v", err)
	}

	bytes, err := ioutil.ReadAll(readCloser)
	if err != nil {
		return "", fmt.Errorf("Cannot read file: %v", err)
	}

	err = ioutil.WriteFile(filePath, bytes, 0777)
	if err != nil {
		return "", fmt.Errorf("Cannot write file to fs: %v", err)
	}

	return filePath, nil
}
