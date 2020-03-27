package routers

import (
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/api"
	"school/helpers"
)

func HandlerClassroomDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()

	params := mux.Vars(r)

	classroomID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerClassroomDetail/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	param := api.ClassroomDetailParam{ID: classroomID}

	return classroomService.Detail(ctx, param)
}
