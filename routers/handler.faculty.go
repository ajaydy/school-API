package routers

import (
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/api"
	"school/helpers"
)

func HandlerFacultyDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()

	params := mux.Vars(r)

	facultyID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerFacultyDetail/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	param := api.FacultyDetailParam{ID: facultyID}

	return facultyService.Detail(ctx, param)
}
