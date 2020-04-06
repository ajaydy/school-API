package util

import (
	"errors"
	"math"
	"math/rand"
	"strconv"
	"time"
)

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func GetGender(gender int) (string, error) {
	switch gender {
	case 0:
		return "Male", nil
	case 1:
		return "Female", nil

	default:
		return "", errors.New("Invalid")
	}
}

func GetAge(dateOfBirth time.Time) float64 {
	return math.Floor(time.Now().Sub(dateOfBirth).Hours() / 24 / 365)
}

func GetYearCode() (int, error) {
	currentTime := time.Now()
	year := currentTime.Year()

	a := strconv.Itoa(year)

	yearCode := a[2:4]

	code, err := strconv.Atoi(yearCode)

	if err != nil {
		return 0, err
	}

	return code, nil

}

func GetGrade(marks int) string {

	if marks <= 100 && marks >= 75 {
		return "A"
	} else if marks < 75 && marks >= 60 {
		return "B"
	} else if marks < 60 && marks >= 47 {
		return "C"
	} else if marks < 47 && marks >= 40 {
		return "D"
	} else {
		return "F"
	}
}

func GetDay(day int) string {

	switch day {
	case 1:
		return "Monday"
	case 2:
		return "Tuesday"
	case 3:
		return "Wednesday"
	case 4:
		return "Thursday"
	case 5:
		return "Friday"
	default:
		return "Invalid"
	}
}

func GetTrimester(month int) int {

	switch month {
	case 4:
		return 1
	case 7:
		return 2
	case 11:
		return 3
	default:
		return 0
	}

}
