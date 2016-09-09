package market

import (
	"database/sql"
	//"fmt"
	//"errors"

	//stat "github.com/asiainfoLDP/datafoundry_appmarket/statistics"
	"github.com/asiainfoLDP/datahub_commons/log"
)



type DatabaseUpgrader_2 struct {
	DatabaseUpgrader_Base

	AlterSQL string
}

func newDatabaseUpgrader_2() *DatabaseUpgrader_2 {
	updater := &DatabaseUpgrader_2{}
	
	updater.currentTableCreationSqlFile = "initdb_v003.sql"
	
	updater.oldVersion = 2
	updater.newVersion = 3

	return updater
}

func (upgrader DatabaseUpgrader_2) Upgrade (db *sql.DB) error {

	log.DefaultLogger().Info("DatabaseUpgrader_2 started ... ") 
	
	// ...
	
	log.DefaultLogger().Info("DatabaseUpgrader_2 alter tables ... ") 
	
	_ = CreateApp(db, &appNewRelic)
	//if err != nil {
	//	return err
	//}

	_, _ = db.Exec("drop table DF_SAAS_APP")
	_, _ = db.Exec("drop table DF_SAAS_APP_INSTANCE")
	
	log.DefaultLogger().Info("DatabaseUpgrader_2, alter tables done. ")

	return nil
}


var appNewRelic = SaasApp{
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
	//Create_time: time.Now(),

}

/*
var appNewSMS = SaasApp{
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
	//Create_time: time.Now(),

}
*/
