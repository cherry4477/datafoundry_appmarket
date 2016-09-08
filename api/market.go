package api

import (
	"net/http"
	"time"
	"crypto/rand"
	"fmt"
	mathrand "math/rand"

	"github.com/julienschmidt/httprouter"

	"github.com/asiainfoLDP/datahub_commons/common"

	"github.com/asiainfoLDP/datafoundry_appmarket/market"
)

//==================================================================
//
//==================================================================

func init() {
	mathrand.Seed(time.Now().UnixNano())
}

func genUUID() string {
	bs := make([]byte, 16)
	_, err := rand.Read(bs)
	if err != nil {
		Logger.Warning("genUUID error: ", err.Error())

		mathrand.Read(bs)
	}

	return fmt.Sprintf("%X-%X-%X-%X-%X", bs[0:4], bs[4:6], bs[6:8], bs[8:10], bs[10:])
}

//==================================================================
//
//==================================================================

func validateAppID(appId string) *Error {
	// GetError2(ErrorCodeInvalidParameters, err.Error())
	_, e := _mustStringParam("appid", appId, 50, StringParamType_UnicodeUrlWord)
	return e
}

func validateAppInfo(app *market.SaasApp) *Error {

	/*
	app.Provider
	app.Url
	app.Name
	app.Version
	app.Category
	app.Description
	app.Icon_url
	*/

	return nil
}

func validateProvider(provider string) (string, *Error) {
	if provider != "" {
		provider_param, e := _mustStringParam("provider", provider, 50, StringParamType_UnicodeUrlWord)
		if e != nil {
			return "", e
		}
		provider = provider_param
	}

	return provider, nil
}

func validateCategory(category string) (string, *Error) {
	if category != "" {
		category_param, e := _mustStringParam("category", category, 25, StringParamType_UnicodeUrlWord)
		if e != nil {
			return "", e
		}
		category = category_param
	}

	return category, nil
}

//==================================================================
//
//==================================================================

func CreateApp(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// todo: auth
	
	// ...
	db := getDB()
	if db == nil {
		JsonResult(w, http.StatusInternalServerError, GetError(ErrorCodeDbNotInitlized), nil)
		return
	}

	app := &market.SaasApp{}
	err := common.ParseRequestJsonInto(r, app)
	if err != nil {
		JsonResult(w, http.StatusBadRequest, GetError2(ErrorCodeParseJsonFailed, err.Error()), nil)
		return
	}

	e := validateAppInfo(app)
	if e != nil {
		JsonResult(w, http.StatusBadRequest, e, nil)
		return
	}
	
	app.App_id = genUUID()
	// followings will be ignored
	//app.Create_time = time.Now()
	//app.Hotness = 0

	err = market.CreateApp(db, app)
	if err != nil {
		JsonResult(w, http.StatusBadRequest, GetError2(ErrorCodeCreateApp, err.Error()), nil)
		return
	}

	JsonResult(w, http.StatusOK, nil, app.App_id)
}

func DeleteApp(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// todo: auth
	
	// ...
	db := getDB()
	if db == nil {
		JsonResult(w, http.StatusInternalServerError, GetError(ErrorCodeDbNotInitlized), nil)
		return
	}

	appId := params.ByName("appid")

	e := validateAppID(appId)
	if e != nil {
		JsonResult(w, http.StatusBadRequest, e, nil)
		return
	}

	err := market.DeleteApp(db, appId)
	if err != nil {
		JsonResult(w, http.StatusBadRequest, GetError2(ErrorCodeDeleteApp, err.Error()), nil)
		return
	}

	JsonResult(w, http.StatusOK, nil, nil)
}

func ModifyApp(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// todo: auth
	
	// ...
	db := getDB()
	if db == nil {
		JsonResult(w, http.StatusInternalServerError, GetError(ErrorCodeDbNotInitlized), nil)
		return
	}

	appId := params.ByName("appid")

	e := validateAppID(appId)
	if e != nil {
		JsonResult(w, http.StatusBadRequest, e, nil)
		return
	}

	app := &market.SaasApp{}
	err := common.ParseRequestJsonInto(r, app)
	if err != nil {
		JsonResult(w, http.StatusBadRequest, GetError2(ErrorCodeParseJsonFailed, err.Error()), nil)
		return
	}

	e = validateAppInfo(app)
	if e != nil {
		JsonResult(w, http.StatusBadRequest, e, nil)
		return
	}

	err = market.ModifyApp(db, app)
	if err != nil {
		JsonResult(w, http.StatusBadRequest, GetError2(ErrorCodeModifyApp, err.Error()), nil)
		return
	}


	JsonResult(w, http.StatusOK, nil, nil)
}

func RetrieveApp(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	JsonResult(w, http.StatusOK, nil, appNewRelic)
	return

	// todo: auth
	
	// ...
	db := getDB()
	if db == nil {
		JsonResult(w, http.StatusInternalServerError, GetError(ErrorCodeDbNotInitlized), nil)
		return
	}

	appId := params.ByName("appid")

	e := validateAppID(appId)
	if e != nil {
		JsonResult(w, http.StatusBadRequest, e, nil)
		return
	}

	app, err := market.RetrieveAppByID(db, appId)
	if err != nil {
		JsonResult(w, http.StatusBadRequest, GetError2(ErrorCodeGetApp, err.Error()), nil)
		return
	}

	JsonResult(w, http.StatusOK, nil, app)
}

/*
category: app的类别。可选。如果忽略，表示所有类别。
provider: 提供方。可选。如果忽略，表示所有提供方。
orderby: 排序依据。可选。合法值包括hotness|createtime，默认为hotness。
sortOrder: 排序方向。可选。合法值包括asc|desc，默认为desc。
page: 第几页。可选。最小值为1。默认为1。
size: 每页最多返回多少条数据。可选。最小为1，最大为100。默认为30。
*/

func QueryAppList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	apps := []*market.SaasApp{
		&appNewRelic,
	}

	JsonResult(w, http.StatusOK, nil, newQueryListResult(int64(len(apps)), apps))
	return



	// todo: auth
	
	// ...
	db := getDB()
	if db == nil {
		JsonResult(w, http.StatusInternalServerError, GetError(ErrorCodeDbNotInitlized), nil)
		return
	}

	r.ParseForm()

	provider, e := validateProvider(r.Form.Get("provider"))
	if e != nil {
		JsonResult(w, http.StatusBadRequest, e, nil)
		return
	}

	category, e := validateCategory(r.Form.Get("category"))
	if e != nil {
		JsonResult(w, http.StatusBadRequest, e, nil)
		return
	}
	
	offset, size := optionalOffsetAndSize(r, 30, 1, 100)
	orderBy := market.ValidateOrderBy(r.Form.Get("orderby"))
	sortOrder := market.ValidateSortOrder(r.Form.Get("sortorder"), false)

	count, apps, err := market.QueryApps(db, provider, category, orderBy, sortOrder, offset, size)
	if err != nil {
		JsonResult(w, http.StatusBadRequest, GetError2(ErrorCodeQueryApps, err.Error()), nil)
		return
	}

	JsonResult(w, http.StatusOK, nil, newQueryListResult(count, apps))
}


var appNewRelic = market.SaasApp{
	App_id:      "98DED98A-F7A1-EDF2-3DF7-B799333D2FD2",
	Provider:    "New Relic",
	Url:         "https://dashboard.daocloud.io/orgs/asiainfo_dev/services/fec195f5-3440-4f13-94da-48d5008b6eb6",
	Name:        "New Relic",
	Version:     "",
	Category:    "monitor",
	Description: 	`
New Relic是一款基于 SaaS 的云端应用监测与管理平台，可以监测和管理云端、网络端及移动端的应用，能让开发者以终端用户、服务器端或应用代码端的视角来监控自己的应用。 目前New Relic 提供的服务包括终端用户行为监控、应用监控、数据库监控、基础底层监控以及单个平台的监控，能为应用的健康提供实时的可预见性。例如，当出现大量用户无法登录帐号时，New Relic 提供的实时服务能让用户在投诉蜂拥而至之前找到问题的症结所在，进而让开发运营团队实时管理其应用的表现。
`,
	Icon_url:    "https://dn-dao-pr.qbox.me/website/icon/yEDRfH2o.jpeg",
	Create_time: time.Now(),

}

/*
var appNewSMS = market.SaasApp{
	App_id:      "DC3E7112-4202-8593-771D-824197CE79D0",
	Provider:    "AsiaInfo",
	Url:         "http://124.207.3.112:18351/smsservice/send",
	Name:        "SMS Gateway",
	Version:     "",
	Category:    "sms",
	Description: 	`
亚信短信网关
`,
	Icon_url:    "",
	Create_time: time.Now(),

}
*/