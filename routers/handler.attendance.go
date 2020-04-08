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

func HandlerAttendanceList(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()
	filter, err := helpers.ParseFilter(ctx, r)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerAttendanceList/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return attendanceService.List(ctx, filter)
}

func HandlerAttendanceAdd(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.AttendanceAddParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerAddAttendance/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	return attendanceService.Add(ctx, param)
}

func HandlerAttendanceUpdate(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	attendanceID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerAttendanceUpdate/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.AttendanceUpdateParam

	err = helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {

		return nil, helpers.ErrorWrap(err, "handler", "HandlerAttendanceUpdate/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	param.ID = attendanceID

	return attendanceService.Update(ctx, param)
}

func HandlerAttendanceListByClass(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()

	params := mux.Vars(r)

	classID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerAttendanceListByClass/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.AttendanceListByClassParam

	err = helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {

		return nil, helpers.ErrorWrap(err, "handler", "HandlerAttendanceListByClass/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	param.ClassID = classID

	filter, err := helpers.ParseFilter(ctx, r)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerAttendanceListByClass/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return attendanceService.ListByClass(ctx, filter, param)
}
