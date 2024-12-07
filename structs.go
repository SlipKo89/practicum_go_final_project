package main

import (
	_ "modernc.org/sqlite"
)

type empty struct {
}

type Task struct {
	ID      string `json:"id",omitempty`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment",omitempty`
	Repeat  string `json:"repeat",omitempty`
}

// create struct variable for create json
type TaskRes struct {
	ID string `json:"id"`
}

type AddTaskError struct {
	Error string `json:"error"`
}
