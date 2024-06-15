package controllers

import (
	"net/http"

	"github.com/rohan031/identity/helper"
	"github.com/rohan031/identity/services"
)

func GetIdentity(w http.ResponseWriter, r *http.Request) {
	var payload services.JSONResponse
	payload.Error = false
	payload.Message = "endpoint is working"

	helper.EncodeJson(w, http.StatusOK, payload)
}
