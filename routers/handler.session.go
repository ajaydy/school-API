package routers

import (
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/api"
	"school/helpers"
)

func HandlerSessionDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()

	params := mux.Vars(r)

	sessionID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerSessionDetail/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	param := api.SessionDetailParam{ID: sessionID}

	return sessionService.Detail(ctx, param)
}

func HandlerSessionList(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()

	filter, err := helpers.ParseFilter(ctx, r)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerSessionList/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return sessionService.List(ctx, filter)
}

func HandlerSessionListByLecturer(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	filter, err := helpers.ParseFilter(ctx, r)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", " HandlerSessionByLecturer/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return sessionService.ListByLecturer(ctx, filter)
}

func HandlerSessionAdd(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.SessionAddParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerSessionAdd/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	return sessionService.Add(ctx, param)
}

func HandlerSessionUpdate(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	sessionID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerSessionUpdate/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.SessionUpdateParam

	err = helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {

		return nil, helpers.ErrorWrap(err, "handler", "HandlerResultUpdate/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	param.ID = sessionID

	return sessionService.Update(ctx, param)
}

func HandlerSessionDelete(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	sessionID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerSessionDelete/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.SessionDeleteParam

	param.ID = sessionID

	return sessionService.Delete(ctx, param)
}
