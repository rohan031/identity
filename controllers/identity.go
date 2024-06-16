package controllers

import (
	"net/http"

	"github.com/rohan031/identity/custom"
	"github.com/rohan031/identity/helper"
	"github.com/rohan031/identity/services"
)

func GetIdentity(w http.ResponseWriter, r *http.Request) {
	user, err := helper.DecodeJson[services.Identity](w, r)
	if err != nil {
		helper.HandleError(w, err)
		return
	}

	// validate req.body
	if valid := user.ValidateBody(); !valid {
		helper.HandleError(
			w,
			&custom.MalformedRequest{
				Status:  http.StatusBadRequest,
				Message: "Invalid req body values",
			},
		)
		return
	}

	res, err := user.GetIdentity()
	if err != nil {
		helper.HandleError(w, err)
		return
	}

	var payload services.JSONResponse
	payload.Data = res

	helper.EncodeJson(w, http.StatusOK, payload)
}
