package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func handleTask(service TaskService) func(res http.ResponseWriter, req *http.Request) {

	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json; charset=UTF-8")
		switch req.Method {
		case http.MethodGet:
			id, err := strconv.Atoi(req.URL.Query().Get("id"))

			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				r, _ := json.Marshal(Err{fmt.Sprint(err)})
				res.Write(r)
				return
			}

			task, err := service.getTask(id)

			if err != nil {
				res.WriteHeader(http.StatusNotFound)
				r, _ := json.Marshal(Err{fmt.Sprint(err)})
				res.Write(r)
				return
			}

			r, _ := json.Marshal(task)
			res.Write(r)
			return

		case http.MethodPut:
			var task Task
			var buf bytes.Buffer

			_, err := buf.ReadFrom(req.Body)

			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				r, _ := json.Marshal(Err{fmt.Sprint(err)})
				res.Write(r)
				return
			}

			if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
				res.WriteHeader(http.StatusBadRequest)
				r, _ := json.Marshal(Err{fmt.Sprint(err)})
				res.Write(r)
				return
			}

			if err = service.validate(&task); err != nil {
				res.WriteHeader(http.StatusBadRequest)
				r, _ := json.Marshal(Err{fmt.Sprint(err)})
				res.Write(r)
				return
			}

			if err = service.update(task); err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				r, _ := json.Marshal(Err{fmt.Sprint(err)})
				res.Write(r)
				return
			}

			res.WriteHeader(http.StatusOK)
			r, _ := json.Marshal(Ok{})
			res.Write(r)

		case http.MethodPost:
			var task Task
			var buf bytes.Buffer

			_, err := buf.ReadFrom(req.Body)

			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				r, _ := json.Marshal(Err{fmt.Sprint(err)})
				res.Write(r)
				return
			}

			if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
				res.WriteHeader(http.StatusBadRequest)
				r, _ := json.Marshal(Err{fmt.Sprint(err)})
				res.Write(r)
				return
			}

			task.ID = "0"

			if err = service.validate(&task); err != nil {
				res.WriteHeader(http.StatusBadRequest)
				r, _ := json.Marshal(Err{fmt.Sprint(err)})
				res.Write(r)
				return
			}

			taskId, err := service.add(task)

			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				r, _ := json.Marshal(Err{fmt.Sprint(err)})
				res.Write(r)
				return
			}

			res.WriteHeader(http.StatusOK)
			r, _ := json.Marshal(map[string]int{"id": taskId})
			res.Write(r)

		case http.MethodDelete:
			id, err := strconv.Atoi(req.URL.Query().Get("id"))
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				r, _ := json.Marshal(Err{fmt.Sprint(err)})
				res.Write(r)
				return
			}
			task, err := service.getTask(id)

			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				r, _ := json.Marshal(Err{fmt.Sprint(err)})
				res.Write(r)
				return
			}

			if err = service.delete(task); err != nil {
				res.WriteHeader(http.StatusBadRequest)
				r, _ := json.Marshal(Err{fmt.Sprint(err)})
				res.Write(r)
				return
			}

			res.WriteHeader(http.StatusOK)
			r, _ := json.Marshal(Ok{})
			res.Write(r)

		case http.MethodOptions:
			res.Header().Set("Allow", "GET, POST, PUT, DELETE, OPTIONS")
			res.WriteHeader(http.StatusNoContent)

		default:
			res.Header().Set("Allow", "GET, POST, PUT, DELETE, OPTIONS")
			res.WriteHeader(http.StatusMethodNotAllowed)
			r, _ := json.Marshal(Err{"метод не поддерживается"})
			res.Write(r)
		}
	}

}

func handleTasks(service TaskService) func(res http.ResponseWriter, req *http.Request) {

	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json; charset=UTF-8")
		switch req.Method {
		case http.MethodGet:
			search := req.URL.Query().Get("search")

			tasks, _ := service.getTasks(search)
			res.WriteHeader(http.StatusOK)
			r, _ := json.Marshal(map[string][]Task{"tasks": tasks})
			res.Write(r)

		case http.MethodOptions:
			res.Header().Set("Allow", "GET")
			res.WriteHeader(http.StatusNoContent)

		default:

			res.Header().Set("Allow", "GET")
			res.WriteHeader(http.StatusMethodNotAllowed)
			r, _ := json.Marshal(Err{"метод не поддерживается"})
			res.Write(r)
		}

	}
}

func handleNextDate(res http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case http.MethodGet:
		date := req.URL.Query().Get("date")
		repeat := req.URL.Query().Get("repeat")
		now, err := time.Parse(dateLayout, req.URL.Query().Get("now"))

		if err != nil {
			res.Header().Set("Content-Type", "application/json; charset=UTF-8")
			res.WriteHeader(http.StatusBadRequest)
			r, _ := json.Marshal(Err{fmt.Sprint(err)})
			res.Write(r)
			return
		}
		res.Header().Set("Content-Type", "text/plain; charset=utf-8")
		res.WriteHeader(http.StatusOK)
		r, _ := nextDate(now, date, repeat)
		res.Write([]byte(r))

	case http.MethodOptions:
		res.Header().Set("Allow", "GET")
		res.WriteHeader(http.StatusNoContent)

	default:
		res.Header().Set("Content-Type", "application/json; charset=UTF-8")
		res.Header().Set("Allow", "GET")
		res.WriteHeader(http.StatusMethodNotAllowed)
		r, _ := json.Marshal(Err{"метод не поддерживается"})
		res.Write(r)
	}

}

func handleDone(service TaskService) func(res http.ResponseWriter, req *http.Request) {

	return func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			id, err := strconv.Atoi(req.URL.Query().Get("id"))

			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				r, _ := json.Marshal(Err{fmt.Sprint(err)})
				res.Write(r)
				return
			}

			task, err := service.getTask(id)

			if err != nil {
				res.WriteHeader(http.StatusNotFound)
				r, _ := json.Marshal(Err{fmt.Sprint(err)})
				res.Write(r)
				return
			}

			if task.Repeat != "" {
				task.Date, _ = nextDate(time.Now(), task.Date, task.Repeat)
				err = service.update(task)

				if err != nil {
					res.WriteHeader(http.StatusBadRequest)
					r, _ := json.Marshal(Err{fmt.Sprint(err)})
					res.Write(r)
					return
				}

			}

			if task.Repeat == "" {
				err = service.delete(task)

				if err != nil {
					res.WriteHeader(http.StatusBadRequest)
					r, _ := json.Marshal(Err{fmt.Sprint(err)})
					res.Write(r)
					return
				}

			}

			res.Header().Set("Content-Type", "application/json; charset=UTF-8")
			res.WriteHeader(http.StatusOK)
			r, _ := json.Marshal(Ok{})
			res.Write(r)

		case http.MethodOptions:
			res.Header().Set("Allow", "GET")
			res.WriteHeader(http.StatusNoContent)

		default:
			res.Header().Set("Allow", "GET")
			res.WriteHeader(http.StatusMethodNotAllowed)
			r, _ := json.Marshal(Err{"метод не поддерживается"})
			res.Write(r)
		}
	}

}
