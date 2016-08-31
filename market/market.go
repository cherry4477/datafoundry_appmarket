package market

import (
	//"database/sql"
	//"errors"
	//"fmt"
	"time"
	//"bytes"
	//"strings"
	//"io/ioutil"
	//"path/filepath"s

	//stat "github.com/asiainfoLDP/datafoundry_appmarket/statistics"
	//"github.com/asiainfoLDP/datahub_commons/log"
)

//=============================================================
//
//=============================================================

type SaasApp struct {
	App_id      int       `json:"appId,omitempty"`
	Provider    string    `json:"provider,omitempty"`
	Name        string    `json:"name,omitempty"`
	Version     string    `json:"version,omitempty"`
	Category    string    `json:"category,omitempty"`
	Description string    `json:"description,omitempty"`
	Icon_url    string    `json:"iconUrl,omitempty"`
	Create_time time.Time `json:"createTime,omitempty"`
	Hotness     int       `json:"-"`
}

//=============================================================
// 
//=============================================================
/*
func CreateApp(db *sql.DB, userName string, repoName string, itemName string) (bool, error) {
	star, err := RetrieveAppByUserAndItem(db, userName, repoName, itemName)
	if star != nil {
		return false, errors.New("already subscribed")
	} else if err != nil {
		return false, err
	}

	nowstr := time.Now().Format("2006-01-02 15:04:05.999999")
	sqlstr := fmt.Sprintf(`insert into DF_SAAS_APP
							(USER_NAME, REPOSITORY_NAME, DATAITEM_NAME, CREATE_TIME)
							values ('%s', '%s', '%s', '%s')
							`, userName, repoName, itemName, nowstr)
	_, err = db.Exec(sqlstr)
	if err != nil {
		return false, err
	}

	go func() {
		stat.UpdateStat(db, stat.GetAppsStatKey(repoName, itemName), 1)
		stat.UpdateStat(db, stat.GetAppsStatKey(repoName), 1)
	}()

	return true, nil
}

func DeleteApp(db *sql.DB, userName string, repoName string, itemName string) (bool, error) {
	sqlstr := fmt.Sprintf(`delete from DF_SAAS_APP
							where USER_NAME='%s' and REPOSITORY_NAME='%s' and DATAITEM_NAME='%s'
							`, userName, repoName, itemName)
	result, err := db.Exec(sqlstr)
	if err != nil {
		return false, err
	}

	n, _ := result.RowsAffected()
	if n > 0 {
		go func() {
			stat.UpdateStat(db, stat.GetAppsStatKey(repoName, itemName), -int(n))
			stat.UpdateStat(db, stat.GetAppsStatKey(repoName), -int(n))
		}()
	}

	return true, nil
}

func RetrieveAppByUserAndItem(db *sql.DB, userName string, repoName string, itemName string) (*SaasApp, error) {
	return getSingleApp(db,
		fmt.Sprintf("USER_NAME='%s' and REPOSITORY_NAME='%s' and DATAITEM_NAME='%s'", userName, repoName, itemName))
}

func RetrieveAppByID(db *sql.DB, AppId int) (*SaasApp, error) {
	return getSingleApp(db, fmt.Sprintf("App_ID=%d", AppId))
}

func getSingleApp(db *sql.DB, sqlWhere string) (*SaasApp, error) {
	stars, err := queryApps(db, sqlWhere, 1, 0)
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

func GetUserApps(db *sql.DB, userName string, offset int64, limit int, sortOrder bool) (int64, []*SaasApp, error) {
	count, stars, err := getAppList(db, offset, limit, fmt.Sprintf("USER_NAME='%s'", userName), sortOrder)
	for i := len(stars) - 1; i >= 0; i-- {
		star := stars[i]
		star.User_name = ""
	}
	return count, stars, err
}

func GetUserAppsInRepository(db *sql.DB, userName string, repoName string, offset int64, limit int, sortOrder bool) (int64, []*SaasApp, error) {
	count, stars, err := getAppList(db, offset, limit,
		fmt.Sprintf("USER_NAME='%s' and REPOSITORY_NAME='%s'", userName, repoName), sortOrder)
	for i := len(stars) - 1; i >= 0; i-- {
		star := stars[i]
		star.User_name = ""
		star.Repository_name = ""
	}
	return count, stars, err
}

func GetAppsInRepository(db *sql.DB, repoName string, offset int64, limit int, sortOrder bool) (int64, []*SaasApp, error) {
	count, stars, err := getAppList(db, offset, limit,
		fmt.Sprintf("REPOSITORY_NAME='%s'", repoName), sortOrder)
	for i := len(stars) - 1; i >= 0; i-- {
		star := stars[i]
		star.Repository_name = ""
	}
	return count, stars, err
}
*/

//================================================

/*
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

func getAppList(db *sql.DB, offset int64, limit int, sqlWhere string, sortOrder bool) (int64, []*SaasApp, error) {
	if strings.TrimSpace(sqlWhere) == "" {
		return 0, nil, errors.New("sqlWhere can't be blank")
	}

	sql_where_all := sqlWhere
	
	count, err := queryAppsCount(db, sql_where_all)
	if err != nil {
		return 0, nil, err
	}
	if count == 0 {
		return 0, []*SaasApp{}, nil
	}
	validateOffsetAndLimit(count, &offset, &limit)
	
	subs, err := queryApps(db,
		fmt.Sprintf(`%s order by CREATE_TIME %s`, sql_where_all, sortOrderText[sortOrder]),
		limit, offset)
	
	return count, subs, err
}

const sqlSelectCountFromApp = `select COUNT(*) from DF_SAAS_APP`
const sqlSelectAllFromApp = `select
					USER_NAME, REPOSITORY_NAME, DATAITEM_NAME,
					CREATE_TIME
					from DF_SAAS_APP`
func scanAppWithRows(rows *sql.Rows, s *SaasApp) error {
	err := rows.Scan(&s.User_name, &s.Repository_name, &s.Dataitem_name, &s.Optime)
	return err
}


func queryAppsCount(db *sql.DB, sqlWhere string) (int64, error) {
	count := int64(0)
	
	sql_str := fmt.Sprintf(`%s where %s`, sqlSelectCountFromApp, sqlWhere)
	err := db.QueryRow(sql_str).Scan(&count)
	
	return count, err
}

func queryApps(db *sql.DB, sqlWhere string, limit int, offset int64) ([]*SaasApp, error) {
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
		sqlSelectAllFromApp,
		sqlWhere,
		limit,
		offset_str)
	rows, err := db.Query(sql_str)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subs := make([]*SaasApp, 32)
	num := 0
	for rows.Next() {
		s := &SaasApp{}
		if err := scanAppWithRows(rows, s); err != nil {
			return nil, err
		}
		//validateApp(s) // already done in scanAppWithRows
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
*/