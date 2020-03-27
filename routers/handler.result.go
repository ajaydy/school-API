package routers

import (
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"school/api"
	"school/helpers"
)

func HandlerResultDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
	ctx := r.Context()

	params := mux.Vars(r)

	resultID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerResultDetail/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	param := api.ResultDetailParam{ID: resultID}

	return resultService.Detail(ctx, param)
}
