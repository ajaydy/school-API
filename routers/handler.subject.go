package routers

import (
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/api"
	"school/helpers"
)

func HandlerSubjects(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()
	filter, err := helpers.ParseFilter(ctx, r)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerSubjects/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	return subjectService.List(ctx, filter)
}

func HandlerSubjectsDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()

	params := mux.Vars(r)

	subjectID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerSubjectsDetail/parseID", helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return subjectService.Detail(ctx, subjectID)
}

func HandlerAddSubjects(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.SubjectParamAdd

	err := helpers.ParsePOSTRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerAddSubjects/ParsePOSTRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	return subjectService.Add(ctx, param)
}

func HandlerUpdateSubjects(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	subjectID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlersUpdateSubjects/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.SubjectParamUpdate

	err = helpers.ParsePOSTRequestData(ctx, r, &param)
	if err != nil {

		return nil, helpers.ErrorWrap(err, "handler", "HandlerUpdateSubjects/ParsePOSTRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	param.ID = subjectID

	return subjectService.Update(ctx, param)
}
