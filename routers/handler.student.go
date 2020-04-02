package routers

import (
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/api"
	"school/helpers"
)

func HandlerStudentList(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()
	filter, err := helpers.ParseFilter(ctx, r)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerStudents/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return studentService.List(ctx, filter)
}

func HandlerStudentDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()

	params := mux.Vars(r)

	studentID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerStudentDetail/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	param := api.StudentDetailParam{ID: studentID}

	return studentService.Detail(ctx, param)
}

func HandlerStudentAdd(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.StudentParamAdd

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerAddStudent/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	return studentService.Add(ctx, param)
}

func HandlerStudentUpdate(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	studentID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlersUpdateStudents/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.StudentParamUpdate

	err = helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {

		return nil, helpers.ErrorWrap(err, "handler", "HandlerUpdateStudents/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	param.ID = studentID

	return studentService.Update(ctx, param)
}

//
//func HandlerRegisterStudents(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
//
//	ctx := r.Context()
//
//	var param api.StudentParamRegister
//
//	err := helpers.ParseBodyRequestData(ctx, r, &param)
//	if err != nil {
//		return nil, helpers.ErrorWrap(err, "handler", "HandlerRegisterStudents/ParseBodyRequestData",
//			helpers.BadRequestMessage, http.StatusBadRequest)
//
//	}
//
//	return studentService.Register(ctx, param)
//}
//
func HandlerStudentLogin(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.StudentParamLogin

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerLoginStudents/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}
	return studentService.Login(ctx, param)
}
