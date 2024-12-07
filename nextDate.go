package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	// check repeat is empty
	if repeat == "" {
		//fmt.Println("repeat field is empty")
		return "", errors.New("repeat field is empty")
	}

	// check parse date
	dateTime, err := time.Parse("20060102", date)
	if err != nil {
		//fmt.Println("date is incorrect format")
		return "", errors.New("date is incorrect format. Plz use 20060102")
	}

	// split repeat for different args (y or d + day nums )
	repeat_list := strings.Split(repeat, " ")
	/*
		for _, v := range repeat_list {
			fmt.Println(v)
		}
	*/

	// add 1 y
	if len(repeat_list) == 1 && repeat_list[0] == "y" {
		fmt.Println("add year for ", dateTime)
		dateTime = dateTime.AddDate(1, 0, 0)
		duration := dateTime.Sub(now)
		for duration < 0 {
			dateTime = dateTime.AddDate(1, 0, 0)
			duration = dateTime.Sub(now)
		}
		fmt.Println(dateTime)
		return dateTime.Format("20060102"), err

	} else if len(repeat_list) == 2 && repeat_list[0] == "d" {
		days, err := strconv.Atoi(repeat_list[1])
		if err == nil && days < 400 && days > 0 {
			dateTime = dateTime.AddDate(0, 0, days)
			duration := dateTime.Sub(now)
			fmt.Println("add some days")
			for duration < 0 {
				dateTime = dateTime.AddDate(0, 0, days)
				duration = dateTime.Sub(now)
			}
			return dateTime.Format("20060102"), err
		} else {
			fmt.Println("days not between 0 and 400 or not integer")
			return "", errors.New("days not between 0 and 400 or not integer")
		}
	}

	return "", errors.New("repeat is incorrect format. Plz use 'y' or 'd 10'")
}

func getNextDate(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	fmt.Println(now, date, repeat)

	timeNow, _ := time.Parse("20060102", now)
	resp, _ := NextDate(timeNow, date, repeat)
	fmt.Println(resp)
	// в заголовок записываем тип контента, у нас это данные в формате JSON
	//w.Header().Set("Content-Type", "application/json")
	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	w.Write([]byte(resp))
}
