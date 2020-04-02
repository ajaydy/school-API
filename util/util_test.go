package util

import (
	"fmt"
	"log"
	"testing"
)

func TestGetYearCode(t *testing.T) {

	code, err := GetYearCode()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(code)

}
