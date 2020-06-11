package main

import (
	"log"

	"github.com/stanf0rd/utekai/database"
	"github.com/stanf0rd/utekai/sheets"
)

func printAllUsers() {
	users, err := database.GetAllUsers()
	if err != nil {
		log.Fatal(err)
	}

	sheets.PrintUsers(users)
}

func printAllQuestions() {
	questions, err := database.GetAllQuestions()
	if err != nil {
		log.Fatal(err)
	}

	sheets.PrintQuestions(questions)
}
