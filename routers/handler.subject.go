package routers

import (
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/helpers"
)

func HandlerSubjects(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()

	return subjectService.List(ctx)
}

func HandlerSubjectsDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()

	params := mux.Vars(r)

	subjectID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerSubjectsDetail/parseID", helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return subjectService.Detail(ctx, subjectID)
}
