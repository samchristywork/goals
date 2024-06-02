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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	http.HandleFunc("/goals", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT name, startTimestamp, startAmount, endTimestamp, endAmount, currentAmount FROM goals")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer rows.Close()

		var goals []Goal
		for rows.Next() {
			var goal Goal
			rows.Scan(&goal.Name, &goal.StartTimestamp, &goal.StartAmount, &goal.EndTimestamp, &goal.EndAmount, &goal.CurrentAmount)
			goals = append(goals, goal)
		}

		goalsJSON, err := json.Marshal(goals)
		if err != nil {
			fmt.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(goalsJSON)
	})

	http.HandleFunc("/add-goal", func(w http.ResponseWriter, r *http.Request) {
		var goal Goal
		err := json.NewDecoder(r.Body).Decode(&goal)
		if err != nil {
			fmt.Println(err)
			return
		}

		_, err = db.Exec("INSERT INTO goals (name, startTimestamp, startAmount, endTimestamp, endAmount, currentAmount) VALUES (?, ?, ?, ?, ?, ?)", goal.Name, goal.StartTimestamp, goal.StartAmount, goal.EndTimestamp, goal.EndAmount, goal.StartAmount)
		if err != nil {
			fmt.Println(err)
			return
		}

		log.Printf("Added %s", goal.Name)

		w.WriteHeader(http.StatusCreated)
	})

	http.HandleFunc("/update-field", func(w http.ResponseWriter, r *http.Request) {
		var field UpdateField
		err := json.NewDecoder(r.Body).Decode(&field)
		if err != nil {
			fmt.Println(err)
			return
		}

		sql := fmt.Sprintf("UPDATE goals SET %s = ? WHERE name = ?", field.Column)
		_, err = db.Exec(sql, field.Value, field.Name)
		if err != nil {
			fmt.Println(err)
			return
		}

		log.Printf("Updated %s in %s to %s", field.Column, field.Name, field.Value)

		w.WriteHeader(http.StatusNoContent)
	})

	http.HandleFunc("/delete-goal", func(w http.ResponseWriter, r *http.Request) {
		var goal Goal
		err := json.NewDecoder(r.Body).Decode(&goal)
		if err != nil {
			fmt.Println(err)
			return
		}

		_, err = db.Exec("DELETE FROM goals WHERE name = ?", goal.Name)
		if err != nil {
			fmt.Println(err)
			return
		}

		log.Printf("Deleted %s", goal.Name)

		w.WriteHeader(http.StatusNoContent)
	})

	fmt.Println("Server is running on port 8000")
	http.ListenAndServe("0.0.0.0:8000", nil)
}
