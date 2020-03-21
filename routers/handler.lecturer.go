package routers

import (
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/helpers"
)

func HandlerLecturers(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()

	return lecturerService.List(ctx)
}

func HandlerLecturersDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()

	params := mux.Vars(r)

	lecturerID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerLecturersDetail/parseID", helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return lecturerService.Detail(ctx, lecturerID)
}
