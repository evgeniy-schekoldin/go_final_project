package main

import (
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	separator  = ","
	dateLayout = "20060102"

	unknownError          = "неизвестная ошибка"
	invalidParameterError = "параметр задан неверно"
)

type TaskService struct {
	store TaskStore
}

func NewTaskService(store TaskStore) TaskService {
	return TaskService{store: store}
}

func (s TaskService) add(task Task) (int, error) {

	id, err := s.store.add(task)

	if err != nil {
		return 0, err
	}

	task.ID = strconv.Itoa(id)
	return id, nil
}

func (s TaskService) getTasks(search string) ([]Task, error) {

	tasks, err := s.store.getTasks(search)

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s TaskService) getTask(id int) (Task, error) {

	var task Task
	task, err := s.store.getTask(id)

	if err != nil {
		return task, err
	}

	return task, nil
}

func (s TaskService) update(task Task) error {

	id, _ := strconv.Atoi(task.ID)

	if _, err := s.store.getTask(id); err != nil {
		return err
	}

	err := s.store.update(task)

	if err != nil {
		return err
	}

	return nil
}

func (s TaskService) delete(task Task) error {

	id, _ := strconv.Atoi(task.ID)

	if _, err := s.store.getTask(id); err != nil {
		return err
	}

	err := s.store.delete(id)

	if err != nil {
		return err
	}

	return nil
}

func (s TaskService) validate(task *Task) error {

	if _, err := strconv.Atoi(task.ID); err != nil {
		return err
	}

	if task.Title == "" {
		return fmt.Errorf("не указан заголовок задачи")
	}

	if task.Date != "" {

		taskDate, err := time.Parse(dateLayout, task.Date)

		if err != nil {
			return err
		}

		if taskDate.Before(time.Now().Truncate(24*time.Hour)) && task.Repeat != "" {
			task.Date, err = nextDate(time.Now(), task.Date, task.Repeat)
		}

		if err != nil {
			return err
		}

		if taskDate.Before(time.Now().Truncate(24*time.Hour)) && task.Repeat == "" {
			task.Date = time.Now().Format(dateLayout)
		}

	}

	if task.Date == "" {
		task.Date = time.Now().Format(dateLayout)
	}

	return nil
}

func nextDate(now time.Time, date string, repeat string) (string, error) {

	param := strings.Split(repeat, " ")
	nextDate, err := time.Parse(dateLayout, date)

	if err != nil {
		log.Print(fmt.Sprint(err))
		return "", err
	}

	switch {

	case param[0] == "d":

		if len(param) < 2 {
			err := fmt.Errorf(invalidParameterError)
			log.Print(fmt.Sprint(err))
			return "", err
		}

		d, err := strconv.Atoi(param[1])

		if err != nil {
			log.Print(fmt.Sprint(err))
			return "", err
		}

		if d > 400 {
			err := fmt.Errorf(invalidParameterError)
			log.Print(fmt.Sprint(err))
			return "", err
		}

		for ok := false; !ok; ok = nextDate.After(now) {
			nextDate = nextDate.AddDate(0, 0, d)
		}

		return fmt.Sprint(nextDate.Format(dateLayout)), nil

	case param[0] == "y":

		for ok := false; !ok; ok = nextDate.After(now) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}

		return fmt.Sprint(nextDate.Format(dateLayout)), nil

	case param[0] == "w":

		var weekdays []int

		if len(param) < 2 {
			err := fmt.Errorf(invalidParameterError)
			log.Print(fmt.Sprint(err))
			return "", err
		}

		param[1] = strings.Replace(param[1], "7", "0", -1)

		for _, d := range strings.Split(param[1], separator) {
			weekday, err := strconv.Atoi(d)

			if err != nil {
				log.Print(fmt.Sprint(err))
				return "", err
			}

			if weekday > 7 {
				err := fmt.Errorf(invalidParameterError)
				log.Print(fmt.Sprint(err))
				return "", err
			}

			weekdays = append(weekdays, weekday)
		}

		for ok := false; !ok; ok = nextDate.After(now) && slices.Contains(weekdays, int(nextDate.Weekday())) {
			nextDate = nextDate.AddDate(0, 0, 1)
		}

		return fmt.Sprint(nextDate.Format(dateLayout)), nil

	case param[0] == "m":

		var days []int
		var months []int

		for _, d := range strings.Split(param[1], separator) {
			day, err := strconv.Atoi(d)

			if err != nil {
				log.Print(fmt.Sprint(err))
				return "", err
			}

			if day < -2 || day > 31 {
				err := fmt.Errorf(invalidParameterError)
				log.Print(fmt.Sprint(err))
				return "", err
			}

			days = append(days, day)
		}

		if len(param) > 2 {

			for _, m := range strings.Split(param[2], separator) {
				month, err := strconv.Atoi(m)

				if err != nil {
					log.Print(fmt.Sprint(err))
					return "", err
				}

				if month > 12 {
					err := fmt.Errorf(invalidParameterError)
					log.Print(fmt.Sprint(err))
					return "", err
				}

				months = append(months, month)
			}

		}

		if len(months) == 0 {

			for i := 0; i < 12; i++ {
				months = append(months, i+1)
			}

		}

		for {
			nextDate = nextDate.AddDate(0, 0, 1)

			if nextDate.After(now) && (slices.Contains(days, int(nextDate.Day())) && slices.Contains(months, int(nextDate.Month()))) {
				break
			}

			if nextDate.After(now) && slices.Contains(days, -1) && nextDate == nextDate.AddDate(0, 1, -nextDate.Day()) {
				break
			}

			if nextDate.After(now) && slices.Contains(days, -2) && nextDate == nextDate.AddDate(0, 1, -nextDate.Day()-1) {
				break
			}

		}

		return fmt.Sprint(nextDate.Format(dateLayout)), nil

	}

	return "", fmt.Errorf(unknownError)

}
