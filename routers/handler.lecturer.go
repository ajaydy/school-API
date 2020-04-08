package routers

import (
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/api"
	"school/helpers"
)

func HandlerLecturerList(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()
	filter, err := helpers.ParseFilter(ctx, r)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerLecturerList/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return lecturerService.List(ctx, filter)
}

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

	var param api.LecturerLoginParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerLecturerLogin/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}
	return lecturerService.Login(ctx, param)
}

func HandlerLecturerAdd(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.LecturerAddParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerLecturerAdd/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	return lecturerService.Add(ctx, param)
}

func HandlerLecturerUpdate(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	lecturerID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerLecturerUpdate/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.LecturerUpdateParam

	err = helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {

		return nil, helpers.ErrorWrap(err, "handler", "HandlerLecturerUpdate/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	param.ID = lecturerID

	return lecturerService.Update(ctx, param)
}

func HandlerLecturerDelete(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	lecturerID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerLecturerDelete/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.LecturerDeleteParam

	param.ID = lecturerID

	return lecturerService.Delete(ctx, param)
}

func HandlerLecturerPasswordUpdate(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	lecturerID := uuid.FromStringOrNil(ctx.Value("user_id").(string))

	var param api.LecturerPasswordUpdateParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {

		return nil, helpers.ErrorWrap(err, "handler", "HandlerLecturerPasswordUpdate/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	param.ID = lecturerID

	return lecturerService.PasswordUpdate(ctx, param)
}
