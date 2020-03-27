package helpers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/schema"
	"html"
	"net/http"
	"reflect"
	"strings"
)

var decoder = schema.NewDecoder()

type (
	FilterOption struct {
		Limit  int    `json:"limit" schema="limit"`
		Offset int    `json:"offset" schema="offset"`
		Search string `json:"search" schema="search"`
		Dir    string `json:"dir" schema="dir"`
	}

	Filter struct {
		FilterOption `json:"filter,omitempty"`
	}
)

func ParsePOSTRequestData(ctx context.Context, r *http.Request, data interface{}) error {
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {

	}

	value := reflect.ValueOf(data).Elem()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if field.Type() != reflect.TypeOf("") {
			continue
		}
		str := field.Interface().(string)
		field.SetString(html.EscapeString(str))

	}
	valid, err := govalidator.ValidateStruct(data)

	if err != nil {
		return err
	}

	if !valid {
		return errors.New("Invalid data")
	}

	return nil

}

func ParseFilter(ctx context.Context, r *http.Request) (Filter, error) {
	var filter Filter
	err := decoder.Decode(&filter, r.URL.Query())
	if err != nil {
		return filter, nil
	}

	if strings.ToLower(filter.Dir) != "asc" && strings.ToLower(filter.Dir) != "desc" {
		filter.Dir = "ASC"
	}

	return filter, nil
}
