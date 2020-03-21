package routers

import (
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/helpers"
)

func HandlerSemesters(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()

	return semesterService.List(ctx)
}

func HandlerSemestersDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()

	params := mux.Vars(r)

	semesterID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerSemestersDetail/parseID", helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return semesterService.Detail(ctx, semesterID)
}
