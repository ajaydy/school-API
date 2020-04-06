package routers

import (
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/api"
	"school/helpers"
)

func HandlerProgramList(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	filter, err := helpers.ParseFilter(ctx, r)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerProgramList/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return programService.List(ctx, filter)

}

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

func HandlerProgramAdd(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.ProgramAddParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerProgramAdd/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	return programService.Add(ctx, param)
}

func HandlerProgramUpdate(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	programID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerProgramUpdate/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.ProgramUpdateParam

	err = helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {

		return nil, helpers.ErrorWrap(err, "handler", "HandlerProgramUpdate/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	param.ID = programID

	return programService.Update(ctx, param)
}

func HandlerProgramDelete(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	programID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerProgramDelete/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.ProgramDeleteParam

	param.ID = programID

	return programService.Delete(ctx, param)
}
