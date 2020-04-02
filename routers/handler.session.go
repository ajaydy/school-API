package routers

import (
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/api"
	"school/helpers"
)

func HandlerSessionDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()

	params := mux.Vars(r)

	sessionID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerSessionDetail/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	param := api.SessionDetailParam{ID: sessionID}

	return sessionService.Detail(ctx, param)
}

func HandlerSession(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()

	return sessionService.List(ctx)
}
