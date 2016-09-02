package api

import (
	"net/http"
	"time"
	"crypto/rand"
	"fmt"
	mathrand "math/rand"

	"github.com/julienschmidt/httprouter"

	//"github.com/asiainfoLDP/datahub_commons/common"

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

type SaasApp struct {
	App_id      string    `json:"appId,omitempty"`
	Provider    string    `json:"provider,omitempty"`
	Url         string    `json:"url,omitempty"`
	Name        string    `json:"name,omitempty"`
	Version     string    `json:"version,omitempty"`
	Category    string    `json:"category,omitempty"`
	Description string    `json:"description,omitempty"`
	Icon_url    string    `json:"iconUrl,omitempty"`
	Create_time time.Time `json:"createTime,omitempty"`
	Hotness     int       `json:"-"`
}
func CreateApp(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// todo: auth
	
	// ...
	db := getDB()
	if db == nil {
		JsonResult(w, http.StatusInternalServerError, GetError(ErrorCodeDbNotInitlized), nil)
		return
	}

	provider := ""
	url := ""
	name := ""
	version := ""
	category := ""
	description := ""
	iconUrl := ""
	
	appId := genUUID()

	app := &market.SaasApp {
		App_id:      appId,
		Provider:    provider,
		Url:         url,
		Name:        name,
		Version:     version,
		Category:    category,
		Description: description,
		Icon_url:    iconUrl,
	}

	err := market.CreateApp(db, app)
	if err != nil {
		JsonResult(w, http.StatusBadRequest, GetError2(ErrorCodeCreateApp, err.Error()), nil)
		return
	}

	JsonResult(w, http.StatusOK, nil, appId)
}

func DeleteApp(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// todo: auth
	
	// ...
	db := getDB()
	if db == nil {
		JsonResult(w, http.StatusInternalServerError, GetError(ErrorCodeDbNotInitlized), nil)
		return
	}

	appId := ""

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

	appId := genUUID()

	provider := ""
	url := ""
	name := ""
	version := ""
	category := ""
	description := ""
	iconUrl := ""

	app := &market.SaasApp {
		App_id:      appId,
		Provider:    provider,
		Url:         url,
		Name:        name,
		Version:     version,
		Category:    category,
		Description: description,
		Icon_url:    iconUrl,
	}

	app, err := market.ModifyApp(db, app)
	if err != nil {
		JsonResult(w, http.StatusBadRequest, GetError2(ErrorCodeModifyApp, err.Error()), nil)
		return
	}


	JsonResult(w, http.StatusOK, nil, app)
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

	app, err := market.RetrieveApp(db, appId)
	if err != nil {
		JsonResult(w, http.StatusBadRequest, GetError2(ErrorCodeGetApp, err.Error()), nil)
		return
	}

	JsonResult(w, http.StatusOK, nil, app)
}

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

	provider := ""
	category := ""
	orderBy := ""

	apps, err := market.QueryApps(db, provider, category, orderBy)
	if err != nil {
		JsonResult(w, http.StatusBadRequest, GetError2(ErrorCodeQueryApps, err.Error()), nil)
		return
	}

	JsonResult(w, http.StatusOK, nil, newQueryListResult(int64(len(apps)), apps))
}


var appNewRelic = market.SaasApp{
	App_id:      "52FDFC07-2182-654F-163F-5F0F9A621D72",
	Provider:    "DaoCloud",
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