package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	// check repeat is empty
	if repeat == "" {
		return "", errors.New("repeat field is empty")
	}

	// check parse date
	dateTime, err := time.Parse(time_format, date)
	if err != nil {
		return "", errors.New("date is incorrect format. Plz use 20060102")
	}

	// split repeat for different args (y or d + day nums )
	repeat_list := strings.Split(repeat, " ")

	// add 1 y
	if len(repeat_list) == 1 && repeat_list[0] == "y" {
		dateTime = dateTime.AddDate(1, 0, 0)
		duration := dateTime.Sub(now)
		for duration < 0 {
			dateTime = dateTime.AddDate(1, 0, 0)
			duration = dateTime.Sub(now)
		}
		return dateTime.Format(time_format), err

	} else if len(repeat_list) == 2 && repeat_list[0] == "d" {
		days, err := strconv.Atoi(repeat_list[1])
		if err == nil && days < 400 && days > 0 {
			dateTime = dateTime.AddDate(0, 0, days)
			duration := dateTime.Sub(now)
			for duration < 0 {
				dateTime = dateTime.AddDate(0, 0, days)
				duration = dateTime.Sub(now)
			}
			return dateTime.Format(time_format), err
		} else {
			return "", errors.New("days not between 0 and 400 or not integer")
		}
	}

	return "", errors.New("repeat is incorrect format. Plz use 'y' or 'd 10'")
}

func getNextDate(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	timeNow, _ := time.Parse(time_format, now)
	resp, _ := NextDate(timeNow, date, repeat)
	// в заголовок записываем тип контента, у нас это данные в формате JSON
	//w.Header().Set("Content-Type", "application/json")
	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	w.Write([]byte(resp))
}
