package market

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"bytes"
	"strings"
	"io/ioutil"
	"path/filepath"

	stat "github.com/asiainfoLDP/datahub_commons/statistics"
	"github.com/asiainfoLDP/datahub_commons/log"
)

//=============================================================
//
//=============================================================

var dbUpgraders = []DatabaseUpgrader {
	newDatabaseUpgrader_0(),
}

const (
	DbPhase_Unkown    = -1
	DbPhase_Serving   = 0 // must be 0
	DbPhase_Upgrading = 1
)

var dbPhase = DbPhase_Unkown

func IsServing() bool {
	return dbPhase == DbPhase_Serving
}

// for ut, reallyNeedUpgrade is false
func TryToUpgradeDatabase(db *sql.DB, dbName string, reallyNeedUpgrade bool) error {
	
	if reallyNeedUpgrade {
	
		if len(dbUpgraders) == 0 {
			return errors.New("at least one db upgrader needed")
		}
		lastDbUpgrader := dbUpgraders[len(dbUpgraders) - 1]
		
		// create tables. (no table created in following _upgradeDatabase callings)
		
		err := lastDbUpgrader.TryToCreateTables(db)
		if err != nil {
			return err
		}
		
		// init version value stat as LatestDataVersion if it doesn't exist,
		// which means tables are just created. In this case, no upgrations are needed.
		
		// INSERT INTO DH_ITEM_STAT (STAT_KEY, STAT_VALUE) 
		// VALUES(dbName#version, LatestDataVersion) 
		// ON DUPLICATE KEY UPDATE STAT_VALUE=LatestDataVersion;
		dbVersionKey := stat.GetVersionKey(dbName)
		_, err = stat.SetStatIf(db, dbVersionKey, lastDbUpgrader.NewVersion(), 0)
		if err != nil && err != stat.ErrOldStatNotMatch {
			return err
		}
		
		current_version, err := stat.RetrieveStat(db, dbVersionKey)
		if err != nil {
			return err
		}
		
		// upgrade 
		
		if current_version != lastDbUpgrader.NewVersion() {
		
			log.DefaultLogger().Info("mysql start upgrading ...")
		
			dbPhase = DbPhase_Unkown
			
			for _, dbupgrader := range dbUpgraders {
				if err = _upgradeDatabase(db, dbName, dbupgrader); err != nil {
					return err
				}
			}
		}
	}
	
	dbPhase = DbPhase_Serving
	
	log.DefaultLogger().Info("mysql start serving ...")
	
	return nil
}

func _upgradeDatabase(db *sql.DB, dbName string, upgrader DatabaseUpgrader) error {
	dbVersionKey := stat.GetVersionKey(dbName)
	current_version, err := stat.RetrieveStat(db, dbVersionKey)
	if err != nil {
		return err
	}
	if current_version == 0 {
		current_version = 1
	}
	
	log.DefaultLogger().Info("TryToUpgradeDatabase current version: ", current_version) 
	
	if upgrader.NewVersion() <= current_version {
		return nil
	}
	if upgrader.OldVersion() != current_version {
		return fmt.Errorf("old version (%d) <= current version (%d)", upgrader.OldVersion(), current_version)
	}
	
	dbPhaseKey := stat.GetPhaseKey(dbName)
	phase, err := stat.SetStatIf(db, dbPhaseKey, DbPhase_Upgrading, DbPhase_Serving)

	log.DefaultLogger().Info("TryToUpgradeDatabase current phase: ", phase) 
	
	if err != nil {
		return err
	}
	
	// ...
	
	dbPhase = DbPhase_Upgrading
	
	err = upgrader.Upgrade(db)
	if err != nil {
		return err
	}
	
	// ...
	
	_, err = stat.SetStat(db, dbVersionKey, upgrader.NewVersion())
	if err != nil {
		return err
	}
	
	log.DefaultLogger().Info("TryToUpgradeDatabase new version: ", upgrader.NewVersion()) 
	
	time.Sleep(30 * time.Millisecond)
	
	_, err = stat.SetStatIf(db, dbPhaseKey, DbPhase_Serving, DbPhase_Upgrading)
	if err != nil {
		return err
	}
	
	return nil
}

type DatabaseUpgrader interface {
	OldVersion() int
	NewVersion() int
	Upgrade(db *sql.DB) error
	TryToCreateTables(db *sql.DB) error
}

type DatabaseUpgrader_Base struct {
	oldVersion int
	newVersion int
	
	currentTableCreationSqlFile string
}

func (upgrader DatabaseUpgrader_Base) OldVersion() int {
	return upgrader.oldVersion
}

func (upgrader DatabaseUpgrader_Base) NewVersion() int {
	return upgrader.newVersion
}

func (upgrader DatabaseUpgrader_Base) TryToCreateTables(db *sql.DB) error {
	
	if upgrader.currentTableCreationSqlFile == "" {
		return nil
	}
	
	data, err := ioutil.ReadFile(filepath.Join("_db", upgrader.currentTableCreationSqlFile))
	if err != nil {
		return err
	}
	
	sqls := bytes.SplitAfter(data, []byte("DEFAULT CHARSET=UTF8;"))
	sqls = sqls[:len(sqls)-1]
	for _, sql := range sqls {
		_, err = db.Exec(string(sql))
		if err != nil {
			return err
		}
	}
	
	return nil
}

//=============================================================
//
//=============================================================

type Star struct {
	//User_id          int     `json:"userid,omitempty"`
	User_name string `json:"username,omitempty"`
	//Dataitem_id     int     `json:"dataitemid,omitempty"`
	Repository_name string    `json:"repname,omitempty"`
	Dataitem_name   string    `json:"itemname,omitempty"`
	Optime          time.Time `json:"optime,omitempty"`
}

//=============================================================
//
//=============================================================

func CreateStar(db *sql.DB, userName string, repoName string, itemName string) (bool, error) {
	star, err := RetrieveStarByUserAndItem(db, userName, repoName, itemName)
	if star != nil {
		return false, errors.New("already subscribed")
	} else if err != nil {
		return false, err
	}

	nowstr := time.Now().Format("2006-01-02 15:04:05.999999")
	sqlstr := fmt.Sprintf(`insert into DH_STAR
							(USER_NAME, REPOSITORY_NAME, DATAITEM_NAME, OPTIME)
							values ('%s', '%s', '%s', '%s')
							`, userName, repoName, itemName, nowstr)
	_, err = db.Exec(sqlstr)
	if err != nil {
		return false, err
	}

	go func() {
		stat.UpdateStat(db, stat.GetStarsStatKey(repoName, itemName), 1)
		stat.UpdateStat(db, stat.GetStarsStatKey(repoName), 1)
	}()

	return true, nil
}

func CancelStar(db *sql.DB, userName string, repoName string, itemName string) (bool, error) {
	sqlstr := fmt.Sprintf(`delete from DH_STAR
							where USER_NAME='%s' and REPOSITORY_NAME='%s' and DATAITEM_NAME='%s'
							`, userName, repoName, itemName)
	result, err := db.Exec(sqlstr)
	if err != nil {
		return false, err
	}

	n, _ := result.RowsAffected()
	if n > 0 {
		go func() {
			stat.UpdateStat(db, stat.GetStarsStatKey(repoName, itemName), -int(n))
			stat.UpdateStat(db, stat.GetStarsStatKey(repoName), -int(n))
		}()
	}

	return true, nil
}

func RetrieveStarByUserAndItem(db *sql.DB, userName string, repoName string, itemName string) (*Star, error) {
	return getSingleStar(db,
		fmt.Sprintf("USER_NAME='%s' and REPOSITORY_NAME='%s' and DATAITEM_NAME='%s'", userName, repoName, itemName))
}

func RetrieveStarByID(db *sql.DB, StarId int) (*Star, error) {
	return getSingleStar(db, fmt.Sprintf("Star_ID=%d", StarId))
}

func getSingleStar(db *sql.DB, sqlWhere string) (*Star, error) {
	stars, err := queryStars(db, sqlWhere, 1, 0)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	if len(stars) == 0 {
		return nil, nil
	}

	return stars[0], nil
}

func GetUserStars(db *sql.DB, userName string, offset int64, limit int, sortOrder bool) (int64, []*Star, error) {
	count, stars, err := getStarList(db, offset, limit, fmt.Sprintf("USER_NAME='%s'", userName), sortOrder)
	for i := len(stars) - 1; i >= 0; i-- {
		star := stars[i]
		star.User_name = ""
	}
	return count, stars, err
}

func GetUserStarsInRepository(db *sql.DB, userName string, repoName string, offset int64, limit int, sortOrder bool) (int64, []*Star, error) {
	count, stars, err := getStarList(db, offset, limit,
		fmt.Sprintf("USER_NAME='%s' and REPOSITORY_NAME='%s'", userName, repoName), sortOrder)
	for i := len(stars) - 1; i >= 0; i-- {
		star := stars[i]
		star.User_name = ""
		star.Repository_name = ""
	}
	return count, stars, err
}

func GetStarsInRepository(db *sql.DB, repoName string, offset int64, limit int, sortOrder bool) (int64, []*Star, error) {
	count, stars, err := getStarList(db, offset, limit,
		fmt.Sprintf("REPOSITORY_NAME='%s'", repoName), sortOrder)
	for i := len(stars) - 1; i >= 0; i-- {
		star := stars[i]
		star.Repository_name = ""
	}
	return count, stars, err
}

func GetStarsOnDataItem(db *sql.DB, repoName string, itemName string, offset int64, limit int, sortOrder bool) (int64, []*Star, error) {
	count, stars, err := getStarList(db, offset, limit,
		fmt.Sprintf("REPOSITORY_NAME='%s' and DATAITEM_NAME='%s'", repoName, itemName), sortOrder)
	for i := len(stars) - 1; i >= 0; i-- {
		star := stars[i]
		star.Repository_name = ""
		star.Dataitem_name = ""
	}
	return count, stars, err
}

//================================================

func validateOffsetAndLimit(count int64, offset *int64, limit *int) {
	if *limit < 1 {
		*limit = 1
	}
	if *offset >= count {
		*offset = count - int64(*limit)
	}
	if *offset < 0 {
		*offset = 0
	}
	if *offset + int64(*limit) > count {
		*limit = int(count - *offset)
	}
}

const (
	SortOrder_Asc  = "asc"
	SortOrder_Desc = "desc"
)

// true: asc
// false: desc
var sortOrderText = map[bool]string{true: "asc", false: "desc"}

func ValidateSortOrder(sortOrder string, defaultOrder bool) bool {
	switch strings.ToLower(sortOrder) {
	case SortOrder_Asc:
		return true
	case SortOrder_Desc:
		return false
	}
	
	return defaultOrder
}

func getStarList(db *sql.DB, offset int64, limit int, sqlWhere string, sortOrder bool) (int64, []*Star, error) {
	if strings.TrimSpace(sqlWhere) == "" {
		return 0, nil, errors.New("sqlWhere can't be blank")
	}

	sql_where_all := sqlWhere
	
	count, err := queryStarsCount(db, sql_where_all)
	if err != nil {
		return 0, nil, err
	}
	if count == 0 {
		return 0, []*Star{}, nil
	}
	validateOffsetAndLimit(count, &offset, &limit)
	
	subs, err := queryStars(db,
		fmt.Sprintf(`%s order by OPTIME %s`, sql_where_all, sortOrderText[sortOrder]),
		limit, offset)
	
	return count, subs, err
}

const sqlSelectCountFromStar = `select COUNT(*) from DH_STAR`
const sqlSelectAllFromStar = `select
					USER_NAME, REPOSITORY_NAME, DATAITEM_NAME,
					OPTIME
					from DH_STAR`
func scanStarWithRows(rows *sql.Rows, s *Star) error {
	err := rows.Scan(&s.User_name, &s.Repository_name, &s.Dataitem_name, &s.Optime)
	return err
}


func queryStarsCount(db *sql.DB, sqlWhere string) (int64, error) {
	count := int64(0)
	
	sql_str := fmt.Sprintf(`%s where %s`, sqlSelectCountFromStar, sqlWhere)
	err := db.QueryRow(sql_str).Scan(&count)
	
	return count, err
}

func queryStars(db *sql.DB, sqlWhere string, limit int, offset int64) ([]*Star, error) {
	offset_str := ""
	if offset > 0 {
		offset_str = fmt.Sprintf("offset %d", offset)
	}
	sql_str := fmt.Sprintf(`
					%s
					where %s
					limit %d
					%s
					`,
		sqlSelectAllFromStar,
		sqlWhere,
		limit,
		offset_str)
	rows, err := db.Query(sql_str)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subs := make([]*Star, 32)
	num := 0
	for rows.Next() {
		s := &Star{}
		if err := scanStarWithRows(rows, s); err != nil {
			return nil, err
		}
		//validateStar(s) // already done in scanStarWithRows
		if num >= len(subs) {
			subs = append(subs, s)
		} else {
			subs[num] = s
		}
		num++
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subs[:num], nil
}

//=====================================================

func GetRepositoryNumStars(db *sql.DB, repoName string) (int, error) {
	return stat.RetrieveStat(db, stat.GetStarsStatKey(repoName))
}

func GetDataItemNumStars(db *sql.DB, repoName string, itemName string) (int, error) {
	return stat.RetrieveStat(db, stat.GetStarsStatKey(repoName, itemName))
}
