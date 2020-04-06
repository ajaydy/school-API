package routers

import (
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/api"
	"school/helpers"
)

func HandlerIntakeList(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	filter, err := helpers.ParseFilter(ctx, r)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerIntakeList/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return intakeService.List(ctx, filter)

}

func HandlerIntakeDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()

	params := mux.Vars(r)

	intakeID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerIntakeDetail/parseID", helpers.BadRequestMessage, http.StatusBadRequest)
	}

	param := api.IntakeDetailParam{ID: intakeID}

	return intakeService.Detail(ctx, param)
}

func HandlerIntakeAdd(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.IntakeAddParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerIntakeAdd/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	return intakeService.Add(ctx, param)
}

func HandlerIntakeUpdate(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	intakeID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerIntakeUpdate/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.IntakeUpdateParam

	err = helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {

		return nil, helpers.ErrorWrap(err, "handler", "HandlerIntakeUpdate/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	param.ID = intakeID

	return intakeService.Update(ctx, param)
}

func HandlerIntakeDelete(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	intakeID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerIntakeDelete/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.IntakeDeleteParam

	param.ID = intakeID

	return intakeService.Delete(ctx, param)
}
