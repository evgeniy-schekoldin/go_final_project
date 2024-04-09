package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	_ "modernc.org/sqlite"
)

const (
	webDir = "./web"
	dbFile = "scheduler.db"
)

func main() {

	_, err := os.Executable()

	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	var install bool

	if _, err = os.Stat(dbFile); err != nil {
		install = true
		log.Print(err)
	}

	db, err := sql.Open("sqlite", dbFile)

	if err != nil {
		log.Fatal(err)
	}

	if install {

		schema, err := os.Open("schema.sql")

		if err != nil {
			log.Fatal(err)
		}

		query, err := io.ReadAll(schema)

		if err != nil {
			log.Fatal(err)
		}

		if _, err = db.Exec(string(query)); err != nil {
			log.Fatal(err)
		}

	}

	store := NewTaskStore(db)
	service := NewTaskService(store)

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	mux.HandleFunc("/api/nextdate", handleNextDate)
	mux.HandleFunc("/api/task", handleTask(service))
	mux.HandleFunc("/api/task/done", handleDone(service))
	mux.HandleFunc("/api/tasks", handleTasks(service))

	if err := http.ListenAndServe(":7540", mux); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}

}
