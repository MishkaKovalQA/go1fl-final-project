package api

import (
	"errors"
	"fmt"
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
		return "", errors.New("Не указано правило повторения")
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", errors.New("Неверный формат даты")
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("Не указано правило повторения")
	}

	switch parts[0] {
	case "y":
		if len(parts) != 1 {
			return "", errors.New("Неверное правило повторения по годам")
		}

		for {
			date = date.AddDate(1, 0, 0)

			if date.After(now) {
				return date.Format(dateFormat), nil
			}
		}

	case "d":
		if len(parts) != 2 {
			return "", errors.New("Неверное правило повторения по дням")
		}

		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf(
				"Некорректный интервал дней %q: %w",
				parts[1],
				err,
			)
		}

		if days < 1 || days > 400 {
			return "", errors.New("Интервал дней должен быть от 1 до 400")
		}

		for {
			date = date.AddDate(0, 0, days)

			if date.After(now) {
				return date.Format(dateFormat), nil
			}
		}

	case "w":
		if len(parts) != 2 {
			return "", errors.New("Неверное правило повторения по неделям")
		}

		weekdays, err := parseInts(parts[1])
		if err != nil {
			return "", err
		}

		for _, weekday := range weekdays {
			if weekday < 1 || weekday > 7 {
				return "", errors.New("День недели должен быть от 1 до 7")
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

		return "", errors.New(
			"Не удалось найти следующую дату по недельному правилу",
		)

	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return "", errors.New("Неверное правило повторения по месяцам")
		}

		days, err := parseInts(parts[1])
		if err != nil {
			return "", err
		}

		for _, day := range days {
			if day < -2 || day == 0 || day > 31 {
				return "", errors.New("Неверный день месяца")
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
					return "", errors.New("Месяц должен быть от 1 до 12")
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

		return "", errors.New(
			"Не удалось найти следующую дату по месячному правилу",
		)

	default:
		return "", errors.New("Неподдерживаемое правило повторения")
	}
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	nowValue := r.FormValue("now")
	if nowValue != "" {
		parsedNow, err := time.Parse(dateFormat, nowValue)
		if err != nil {
			http.Error(
				w,
				"Параметр now должен быть в формате ГГГГММДД",
				http.StatusBadRequest,
			)
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
			return nil, fmt.Errorf(
				"Некорректное число %q: %w",
				field,
				err,
			)
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
