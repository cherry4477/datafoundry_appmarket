package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	//"github.com/asiainfoLDP/datahub_commons/common"
	
	// "github.com/asiainfoLDP/datafoundry_appmarket/market"
)

//==================================================================
//
//==================================================================

func CreateApp(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	JsonResult(w, http.StatusOK, nil, nil)
}

func DeleteApp(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	JsonResult(w, http.StatusOK, nil, nil)
}

func ModifyApp(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	JsonResult(w, http.StatusOK, nil, nil)
}

func QueryApp(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	JsonResult(w, http.StatusOK, nil, nil)
}

func QueryAppList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	JsonResult(w, http.StatusOK, nil, nil)
}