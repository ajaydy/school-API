package routers

import (
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/helpers"
)

func HandlerStudentEnrollDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()

	params := mux.Vars(r)

	studentEnrollID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerStudentEnrollDetail/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return studentService.Detail(ctx, studentEnrollID)
}
