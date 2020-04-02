package routers

import (
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/api"
	"school/helpers"
)

func HandlerClassBySession(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()
	filter, err := helpers.ParseFilter(ctx, r)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerClassBySession/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return classService.ListBySession(ctx, filter)
}

func HandlerClassDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	classID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerClassDetail/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	param := api.ClassDetailParam{ID: classID}

	return classService.Detail(ctx, param)
}

func HandlerAddClass(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.ClassAddParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerAddClass/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	return classService.Add(ctx, param)
}
