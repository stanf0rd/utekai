package main

import (
	"log"
	"os"
	"strconv"

	"github.com/stanf0rd/utekai/database"
	"github.com/stanf0rd/utekai/generator"
)

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

func pause(u database.User) {
	questionID := chooseQuestion(u)
	if questionID == 0 {
		log.Printf("User #%d already asked max count of questions", u.ID)
		return
	}

	log.Printf("Question ID chosen: %v", questionID)

	q, err := database.GetQuestionByID(questionID)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Question chosen: %v", q.Body)

	p := database.Pause{
		User:     u.ID,
		Question: q.ID,
	}

	if err := p.Save(); err != nil {
		log.Fatalf("Cannot save pause: %v", err)
	}
}
