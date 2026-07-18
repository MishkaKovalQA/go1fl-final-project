package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeat rule is empty")
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", err
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("repeat rule is empty")
	}

	switch parts[0] {
	case "y":
		if len(parts) != 1 {
			return "", errors.New("invalid yearly repeat rule")
		}

		for {
			date = date.AddDate(1, 0, 0)
			if date.After(now) {
				return date.Format(dateFormat), nil
			}
		}

	case "d":
		if len(parts) != 2 {
			return "", errors.New("invalid daily repeat rule")
		}

		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", err
		}

		if days < 1 || days > 400 {
			return "", errors.New("day interval must be from 1 to 400")
		}

		for {
			date = date.AddDate(0, 0, days)
			if date.After(now) {
				return date.Format(dateFormat), nil
			}
		}

	default:
		return "", errors.New("unsupported repeat rule")
	}
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	nowValue := r.FormValue("now")
	if nowValue != "" {
		parsedNow, err := time.Parse(dateFormat, nowValue)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		now = parsedNow
	} else {
		now, _ = time.Parse(dateFormat, now.Format(dateFormat))
	}

	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	next, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, _ = w.Write([]byte(next))
}
