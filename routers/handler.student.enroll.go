package routers

import (
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/api"
	"school/helpers"
)

func HandlerStudentEnrollDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()

	params := mux.Vars(r)

	studentEnrollID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerStudentEnrollDetail/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	param := api.StudentEnrollDetailParam{ID: studentEnrollID}

	return studentEnrollService.Detail(ctx, param)
}

func HandlerStudentEnrollAdd(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.StudentEnrollAddParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerAddStudent/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	return studentEnrollService.Add(ctx, param)
}

func HandlerTimetable(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()
	filter, err := helpers.ParseFilter(ctx, r)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerTimetable/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return studentEnrollService.List(ctx, filter)
}

func HandlerStudentEnrollBySession(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	filter, err := helpers.ParseFilter(ctx, r)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerStudentEnrollBySession/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return studentEnrollService.StudentEnrollListBySession(ctx, filter)
}

func HandlerStudentEnrollDelete(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	studentEnrollID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerStudentEnrollDelete/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.StudentEnrollDeleteParam

	param.ID = studentEnrollID

	return studentEnrollService.Delete(ctx, param)
}
