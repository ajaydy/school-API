package routers

import (
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/api"
	"school/helpers"
)

func HandlerIntake(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()

	return intakeService.List(ctx)
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

//func HandlerAddIntake(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
//
//	ctx := r.Context()
//
//	var param api.IntakeParamAdd
//
//	err := helpers.ParseBodyRequestData(ctx, r, &param)
//	if err != nil {
//		return nil, helpers.ErrorWrap(err, "handler", "HandlerAddIntake/ParseBodyRequestData",
//			helpers.BadRequestMessage, http.StatusBadRequest)
//
//	}
//
//	return intakeService.Add(ctx, param)
//}
//
//func HandlerUpdateIntake(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
//
//	ctx := r.Context()
//
//	params := mux.Vars(r)
//
//	intakeID, err := uuid.FromString(params["id"])
//
//	if err != nil {
//		return nil, helpers.ErrorWrap(err, "handler", "HandlersUpdateIntake/parseID",
//			helpers.BadRequestMessage, http.StatusBadRequest)
//	}
//
//	var param api.IntakeParamUpdate
//
//	err = helpers.ParseBodyRequestData(ctx, r, &param)
//	if err != nil {
//
//		return nil, helpers.ErrorWrap(err, "handler", "HandlerUpdateIntake/ParseBodyRequestData",
//			helpers.BadRequestMessage, http.StatusBadRequest)
//
//	}
//
//	param.ID = intakeID
//
//	return intakeService.Update(ctx, param)
//}
