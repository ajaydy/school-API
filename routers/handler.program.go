package routers

import (
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/api"
	"school/helpers"
)

func HandlerProgramDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()

	params := mux.Vars(r)

	programID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerProgramDetail/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	param := api.ProgramDetailParam{ID: programID}

	return programService.Detail(ctx, param)
}

func HandlerAddProgram(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.ProgramAddParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerAddProgram/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	return programService.Add(ctx, param)
}
