package routers

import (
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/api"
	"school/helpers"
)

func HandlerFacultyList(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	filter, err := helpers.ParseFilter(ctx, r)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerFacultyList/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return facultyService.List(ctx, filter)
}

func HandlerFacultyDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	facultyID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerFacultyDetail/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	param := api.FacultyDetailParam{ID: facultyID}

	return facultyService.Detail(ctx, param)
}

func HandlerFacultyAdd(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.FacultyAddParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerFacultyAdd/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	return facultyService.Add(ctx, param)
}

func HandlerFacultyUpdate(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	facultyID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerFacultyUpdate/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.FacultyUpdateParam

	err = helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {

		return nil, helpers.ErrorWrap(err, "handler", "HandlerUpdateStudents/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	param.ID = facultyID

	return facultyService.Update(ctx, param)
}

func HandlerFacultyDelete(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	facultyID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerFacultyUpdate/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.FacultyDeleteParam

	param.ID = facultyID

	return facultyService.Delete(ctx, param)
}
