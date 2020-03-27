package routers

import (
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/api"
	"school/helpers"
)

func HandlerStudents(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()

	filter, err := helpers.ParseFilter(ctx, r)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerStudents/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return studentService.List(ctx, filter)
}

func HandlerStudentsDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()

	params := mux.Vars(r)

	studentID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerStudentsDetail/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return studentService.Detail(ctx, studentID)
}

func HandlerAddStudents(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.StudentParamAdd

	err := helpers.ParsePOSTRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerAddStudents/ParsePOSTRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	return studentService.Add(ctx, param)
}

func HandlerUpdateStudents(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	studentID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlersUpdateStudents/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.StudentParamUpdate

	err = helpers.ParsePOSTRequestData(ctx, r, &param)
	if err != nil {

		return nil, helpers.ErrorWrap(err, "handler", "HandlerUpdateStudents/ParsePOSTRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	param.ID = studentID

	return studentService.Update(ctx, param)
}

func HandlerRegisterStudents(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.StudentParamRegister

	err := helpers.ParsePOSTRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerRegisterStudents/ParsePOSTRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	return studentService.Register(ctx, param)
}

func HandlerLoginStudents(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.StudentParamLogin

	err := helpers.ParsePOSTRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerLoginStudents/ParsePOSTRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}
	return studentService.Login(ctx, param)
}

func HandlerStudentsWithScore(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()

	filter, err := helpers.ParseFilter(ctx, r)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerStudentsWithScore/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return studentService.ListWithScore(ctx, filter)
}

func HandlerStudentsWithScoreDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()

	params := mux.Vars(r)

	studentID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerStudentsWithScoreDetail/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return studentService.DetailWithScore(ctx, studentID)
}
