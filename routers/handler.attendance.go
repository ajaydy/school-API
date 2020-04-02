package routers

import (
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/api"
	"school/helpers"
)

func HandlerAttendanceDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()

	params := mux.Vars(r)

	attendanceID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerAttendanceDetail/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	param := api.AttendanceDetailParam{ID: attendanceID}

	return attendanceService.Detail(ctx, param)
}

func HandlerAttendanceByClass(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()
	filter, err := helpers.ParseFilter(ctx, r)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerAttendanceByClass/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return attendanceService.ListByClass(ctx, filter)
}

func HandlerAddAttendance(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.AttendanceAddParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerAddAttendance/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	return attendanceService.Add(ctx, param)
}

func HandlerUpdateAttendanceIsAttend(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	attendanceID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerUpdateAttendanceIsAttend/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.AttendanceUpdateParam

	err = helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {

		return nil, helpers.ErrorWrap(err, "handler", "HandlerUpdateAttendanceIsAttend/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	param.ID = attendanceID

	return attendanceService.UpdateIsAttend(ctx, param)
}
