package models

import (
	"errors"
	"math"
	"time"
)

func Gender(gender int) (string, error) {
	switch gender {
	case 0:
		return "Male", nil
	case 1:
		return "Female", nil

	default:
		return "", errors.New("Invalid")
	}
}

func Age(dateOfBirth time.Time) float64 {
	return math.Floor(time.Now().Sub(dateOfBirth).Hours() / 24 / 365)
}
