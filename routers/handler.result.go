package routers

import (
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/api"
	"school/helpers"
)

func HandlerResultDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()

	params := mux.Vars(r)

	resultID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerResultDetail/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	param := api.ResultDetailParam{ID: resultID}

	return resultService.Detail(ctx, param)
}

func HandlerResultList(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()

	filter, err := helpers.ParseFilter(ctx, r)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerResult/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return resultService.List(ctx, filter)
}

//
//func HandlerResultAdd(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
//
//	ctx := r.Context()
//
//	var param api.ResultAddParam
//
//	err := helpers.ParseBodyRequestData(ctx, r, &param)
//	if err != nil {
//		return nil, helpers.ErrorWrap(err, "handler", "HandlerResultAdd/ParseBodyRequestData",
//			helpers.BadRequestMessage, http.StatusBadRequest)
//
//	}
//
//	return resultService.Add(ctx, param)
//}

func HandlerResultUpdate(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	resultID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerResultUpdate/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.ResultUpdateParam

	err = helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {

		return nil, helpers.ErrorWrap(err, "handler", "HandlerResultUpdate/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	param.ID = resultID

	return resultService.Update(ctx, param)
}

func HandlerResultDelete(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	resultID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerResultDelete/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.ResultDeleteParam

	param.ID = resultID

	return resultService.Delete(ctx, param)
}

func HandlerResultListByStudentEnroll(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	studentEnrollID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerResultListByStudentEnroll/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.ResultListByStudentEnrollParam

	err = helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {

		return nil, helpers.ErrorWrap(err, "handler", "HandlerResultListByStudentEnroll/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}
	param.StudentEnrollID = studentEnrollID

	filter, err := helpers.ParseFilter(ctx, r)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerResultListByStudentEnroll/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return resultService.ListByStudentEnroll(ctx, filter, param)
}

func HandlerResultListByOneStudent(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	filter, err := helpers.ParseFilter(ctx, r)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerResultListByStudentEnroll/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return resultService.ListByOneStudent(ctx, filter)
}
