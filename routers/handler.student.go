package routers

import (
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/helpers"
)

func HandlerStudents(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()

	return studentService.List(ctx)
}

func HandlerStudentsDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()

	params := mux.Vars(r)

	studentID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerStudentsDetail/parseID", helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return studentService.Detail(ctx, studentID)
}
