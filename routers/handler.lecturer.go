package routers

import (
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/api"
	"school/helpers"
)

//func HandlerLecturers(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
//	ctx := r.Context()
//	filter, err := helpers.ParseFilter(ctx, r)
//	if err != nil {
//		return nil, helpers.ErrorWrap(err, "handler", "HandlerLecturers/parseFilter",
//			helpers.BadRequestMessage, http.StatusBadRequest)
//	}
//	return lecturerService.List(ctx, filter)
//}

func HandlerLecturerDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()

	params := mux.Vars(r)

	lecturerID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerLecturerDetail/parseID", helpers.BadRequestMessage, http.StatusBadRequest)
	}

	param := api.LecturerDetailParam{ID: lecturerID}

	return lecturerService.Detail(ctx, param)
}

func HandlerLecturerLogin(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.LecturerParamLogin

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerLoginLecturer/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}
	return lecturerService.Login(ctx, param)
}

func HandlerSessionByLecturer(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	filter, err := helpers.ParseFilter(ctx, r)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", " HandlerSessionByLecturer/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return lecturerService.SessionListByLecturer(ctx, filter)
}

func HandlerStudentEnrollBySession(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	filter, err := helpers.ParseFilter(ctx, r)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerStudentEnrollBySession/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return lecturerService.StudentEnrollListBySession(ctx, filter)
}

func HandlerLecturerAddResult(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.LecturerAddResultParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerLecturerAddResult/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	return lecturerService.LecturerAddResult(ctx, param)
}

func HandlerLecturerUpdateResult(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.LecturerUpdateResultParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {

		return nil, helpers.ErrorWrap(err, "handler", "HandlerLecturerUpdateResult/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	return lecturerService.LecturerUpdateResult(ctx, param)
}

func HandlerAddLecturer(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.LecturerParamAdd

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerAddLecturer/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	return lecturerService.Add(ctx, param)
}

//
//func HandlerUpdateLecturers(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
//
//	ctx := r.Context()
//
//	params := mux.Vars(r)
//
//	lecturerID, err := uuid.FromString(params["id"])
//
//	if err != nil {
//		return nil, helpers.ErrorWrap(err, "handler", "HandlersUpdateLecturers/parseID",
//			helpers.BadRequestMessage, http.StatusBadRequest)
//	}
//
//	var param api.LecturerParamUpdate
//
//	err = helpers.ParseBodyRequestData(ctx, r, &param)
//	if err != nil {
//
//		return nil, helpers.ErrorWrap(err, "handler", "HandlerUpdateLecturers/ParseBodyRequestData",
//			helpers.BadRequestMessage, http.StatusBadRequest)
//
//	}
//
//	param.ID = lecturerID
//
//	return lecturerService.Update(ctx, param)
//}
