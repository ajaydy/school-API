package routers

import (
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/api"
	"school/helpers"
)

func HandlerClassroomList(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	filter, err := helpers.ParseFilter(ctx, r)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerClassroomList/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return classroomService.List(ctx, filter)

}

func HandlerClassroomDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()

	params := mux.Vars(r)

	classroomID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerClassroomDetail/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	param := api.ClassroomDetailParam{ID: classroomID}

	return classroomService.Detail(ctx, param)
}

func HandlerClassroomAdd(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.ClassroomAddParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerClassroomAdd/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	return classroomService.Add(ctx, param)
}

func HandlerClassroomUpdate(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	classroomID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerClassroomUpdate/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.ClassroomUpdateParam

	err = helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {

		return nil, helpers.ErrorWrap(err, "handler", "HandlerClassroomUpdate/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	param.ID = classroomID

	return classroomService.Update(ctx, param)
}

func HandlerClassroomDelete(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	classroomID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerClassroomDelete/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.ClassroomDeleteParam

	param.ID = classroomID

	return classroomService.Delete(ctx, param)
}
