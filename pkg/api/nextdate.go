package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	dateFormat  = "20060102"
	maxScanDays = 2000
)

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

	case "w":
		if len(parts) != 2 {
			return "", errors.New("invalid weekly repeat rule")
		}

		weekdays, err := parseInts(parts[1])
		if err != nil {
			return "", err
		}

		for _, weekday := range weekdays {
			if weekday < 1 || weekday > 7 {
				return "", errors.New("weekday must be from 1 to 7")
			}
		}

		if date.Before(now) {
			date = now
		}

		for i := 0; i < maxScanDays; i++ {
			weekday := int(date.Weekday())
			if weekday == 0 {
				weekday = 7
			}

			if date.After(now) && contains(weekdays, weekday) {
				return date.Format(dateFormat), nil
			}

			date = date.AddDate(0, 0, 1)
		}

		return "", errors.New("next weekly date not found")

	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return "", errors.New("invalid monthly repeat rule")
		}

		days, err := parseInts(parts[1])
		if err != nil {
			return "", err
		}

		for _, day := range days {
			if day < -2 || day == 0 || day > 31 {
				return "", errors.New("invalid day of month")
			}
		}

		var months []int

		if len(parts) == 3 {
			months, err = parseInts(parts[2])
			if err != nil {
				return "", err
			}

			for _, month := range months {
				if month < 1 || month > 12 {
					return "", errors.New("month must be from 1 to 12")
				}
			}
		}

		if date.Before(now) {
			date = now
		}

		for i := 0; i < maxScanDays; i++ {
			if date.After(now) {
				if len(months) == 0 || contains(months, int(date.Month())) {
					currentDay := date.Day()
					lastDay := lastDayOfMonth(date)

					for _, allowedDay := range days {
						switch allowedDay {
						case -1:
							if currentDay == lastDay {
								return date.Format(dateFormat), nil
							}

						case -2:
							if currentDay == lastDay-1 {
								return date.Format(dateFormat), nil
							}

						default:
							if currentDay == allowedDay {
								return date.Format(dateFormat), nil
							}
						}
					}
				}
			}

			date = date.AddDate(0, 0, 1)
		}

		return "", errors.New("next monthly date not found")

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

	next, err := NextDate(
		now,
		r.FormValue("date"),
		r.FormValue("repeat"),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, _ = w.Write([]byte(next))
}

func parseInts(value string) ([]int, error) {
	fields := strings.Split(value, ",")
	nums := make([]int, 0, len(fields))

	for _, field := range fields {
		num, err := strconv.Atoi(field)
		if err != nil {
			return nil, err
		}

		nums = append(nums, num)
	}

	return nums, nil
}

func contains(nums []int, target int) bool {
	for _, num := range nums {
		if num == target {
			return true
		}
	}

	return false
}

func lastDayOfMonth(date time.Time) int {
	return time.Date(
		date.Year(),
		date.Month()+1,
		0,
		0,
		0,
		0,
		0,
		date.Location(),
	).Day()
}
