package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
)

type Goal struct {
	Name           string `json:"name"`
	StartTimestamp string `json:"startTimestamp"`
	StartAmount    string `json:"startAmount"`
	EndTimestamp   string `json:"endTimestamp"`
	EndAmount      string `json:"endAmount"`
	CurrentAmount  string `json:"currentAmount"`
}

type UpdateField struct {
	Name   string `json:"name"`
	Column string `json:"column"`
	Value  string `json:"value"`
}

func main() {
	logName := "server.log"
	logFile, err := os.OpenFile(logName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.SetOutput(logFile)

	filename := "data.sqlite"
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS goals (name TEXT, startTimestamp INTEGER, startAmount REAL, endTimestamp INTEGER, endAmount REAL, currentAmount REAL)")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Server is running on port 8000")
	http.ListenAndServe("0.0.0.0:8000", nil)
}
