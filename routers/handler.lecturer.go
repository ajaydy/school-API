package routers

import (
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/api"
	"school/helpers"
)

func HandlerLecturers(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()
	filter, err := helpers.ParseFilter(ctx, r)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerLecturers/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return lecturerService.List(ctx, filter)
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

func HandlerAddLecturers(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.LecturerParamAdd

	err := helpers.ParsePOSTRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerAddLecturers/ParsePOSTRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	return lecturerService.Add(ctx, param)
}

func HandlerUpdateLecturers(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	lecturerID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlersUpdateLecturers/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.LecturerParamUpdate

	err = helpers.ParsePOSTRequestData(ctx, r, &param)
	if err != nil {

		return nil, helpers.ErrorWrap(err, "handler", "HandlerUpdateLecturers/ParsePOSTRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	param.ID = lecturerID

	return lecturerService.Update(ctx, param)
}
