package routers

import (
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/api"
	"school/helpers"
)

func HandlerSubjectList(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()
	filter, err := helpers.ParseFilter(ctx, r)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerSubjectList/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return subjectService.List(ctx, filter)

}

func HandlerSubjectDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()

	params := mux.Vars(r)

	subjectID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerSubjectDetail/parseID", helpers.BadRequestMessage, http.StatusBadRequest)
	}

	param := api.SubjectDetailParam{ID: subjectID}

	return subjectService.Detail(ctx, param)
}

func HandlerSubjectAdd(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.SubjectAddParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerSubjectAdd/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	return subjectService.Add(ctx, param)
}

func HandlerSubjectUpdate(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	subjectID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerSubjectUpdate/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.SubjectUpdateParam

	err = helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {

		return nil, helpers.ErrorWrap(err, "handler", "HandlerSubjectUpdate/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	param.ID = subjectID

	return subjectService.Update(ctx, param)
}

func HandlerSubjectDelete(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	subjectID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerSubjectDelete/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.SubjectDeleteParam

	param.ID = subjectID

	return subjectService.Delete(ctx, param)
}
