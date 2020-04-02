package routers

import (
	"net/http"
	"school/api"
	"school/helpers"
)

func HandlerAdminLogin(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.AdminLoginParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerAdminLogin/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}
	return adminService.Login(ctx, param)
}
