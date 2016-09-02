package market

import (
	"database/sql"
	//"fmt"
	//"errors"

	//stat "github.com/asiainfoLDP/datafoundry_appmarket/statistics"
	//"github.com/asiainfoLDP/datahub_commons/log"
)



type DatabaseUpgrader_2 struct {
	DatabaseUpgrader_Base
}

func newDatabaseUpgrader_2() *DatabaseUpgrader_2 {
	updater := &DatabaseUpgrader_2{}
	
	updater.currentTableCreationSqlFile = "initdb_v003.sql"
	
	updater.oldVersion = 2
	updater.newVersion = 3
	
	return updater
}

func (upgrader DatabaseUpgrader_2) Upgrade (db *sql.DB) error {
	return nil
}
