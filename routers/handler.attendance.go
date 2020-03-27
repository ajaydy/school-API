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
