package routers

import (
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/api"
	"school/helpers"
)

//func HandlerSubjects(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
//	ctx := r.Context()
//	filter, err := helpers.ParseFilter(ctx, r)
//	if err != nil {
//		return nil, helpers.ErrorWrap(err, "handler", "HandlerSubjects/parseFilter",
//			helpers.BadRequestMessage, http.StatusBadRequest)
//	}
//
//	return subjectService.List(ctx, filter)
//}

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

func HandlerAddSubject(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.SubjectParamAdd

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerAddSubject/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	return subjectService.Add(ctx, param)
}

//
//func HandlerUpdateSubjects(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
//
//	ctx := r.Context()
//
//	params := mux.Vars(r)
//
//	subjectID, err := uuid.FromString(params["id"])
//
//	if err != nil {
//		return nil, helpers.ErrorWrap(err, "handler", "HandlersUpdateSubjects/parseID",
//			helpers.BadRequestMessage, http.StatusBadRequest)
//	}
//
//	var param api.SubjectParamUpdate
//
//	err = helpers.ParseBodyRequestData(ctx, r, &param)
//	if err != nil {
//
//		return nil, helpers.ErrorWrap(err, "handler", "HandlerUpdateSubjects/ParseBodyRequestData",
//			helpers.BadRequestMessage, http.StatusBadRequest)
//
//	}
//
//	param.ID = subjectID
//
//	return subjectService.Update(ctx, param)
//}
