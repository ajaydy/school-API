package routers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"school/helpers"
)

type (
	HandlerFunc func(http.ResponseWriter, *http.Request) (interface{}, *helpers.Error)
)

func (fn HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var errs []string
	r.ParseForm()
	data, err := fn(w, r)
	if err != nil {
		errs = append(errs, err.Error())
		w.WriteHeader(err.StatusCode)
	}
	resp := helpers.Response{
		Data: data,
		BaseResponse: helpers.BaseResponse{
			Errors: errs,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		return
	}
}

func InitHandlers() *mux.Router {
	r := mux.NewRouter()

	http.Handle("/", r)

	apiV1 := r.PathPrefix("/api/v1").Subrouter()
	apiV1.Handle("/students", HandlerFunc(HandlerStudents)).Methods(http.MethodGet)
	apiV1.Handle("/students/{id}", HandlerFunc(HandlerStudentsDetail)).Methods(http.MethodGet)
	apiV1.Handle("/lecturers", HandlerFunc(HandlerLecturers)).Methods(http.MethodGet)
	apiV1.Handle("/lecturers/{id}", HandlerFunc(HandlerLecturersDetail)).Methods(http.MethodGet)
	apiV1.Handle("/subjects", HandlerFunc(HandlerSubjects)).Methods(http.MethodGet)
	apiV1.Handle("/subjects/{id}", HandlerFunc(HandlerSubjectsDetail)).Methods(http.MethodGet)

	return r
}
