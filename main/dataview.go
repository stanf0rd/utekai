package main

import (
	"log"

	"github.com/stanf0rd/utekai/database"
	"github.com/stanf0rd/utekai/sheets"
)

func printAllUsers() {
	users, err := database.GetAllUsers()
	if err != nil {
		log.Printf("Cannot get all users: %v", err)
		return
	}

	sheets.PrintUsers(users)
	if err := sheets.PrintUsers(users); err != nil {
		log.Printf("Cannot print users: %v", err)
	}
}

func printAllQuestions() {
	questions, err := database.GetAllQuestions()
	if err != nil {
		log.Printf("Cannot get all questions: %v", err)
		return
	}

	if err := sheets.PrintQuestions(questions); err != nil {
		log.Printf("Cannot print questions: %v", err)
	}
}

func printAllPauses() {
	pauses, err := database.GetAllPauses()
	if err != nil {
		log.Printf("Cannot get all pauses: %v", err)
		return
	}

	if err := sheets.PrintPauses(pauses); err != nil {
		log.Printf("Cannot print pauses: %v", err)
	}
}
