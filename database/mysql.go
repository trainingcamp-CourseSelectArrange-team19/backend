package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"
	"techtrainingcamp-courseSelectArrange/tools"

	_ "github.com/go-sql-driver/mysql"
)

type UpdCnt struct {
	Id      string
	Success bool
}

var db *sql.DB

func CloseMysql() {
	db.Close()
}

func OpenMysql() error {
	var err error
	if os.Getenv("IS_DOCKER") == "1" {
		db, err = sql.Open("mysql", "group1:group1_1@rdsmysqlh138cf58d0185285d.rds.ivolces.com:3306/group1_mysql") //用户名:密码@/数据库名
	} else {
		db, err = sql.Open("mysql", "test:123456@/app") //用户名:密码@/数据库名
	}
	if err != nil {
		tools.LogMsg("数据库链接错误", err)
		// fmt.Println("数据库链接错误", err)
	}
	//延迟到函数结束关闭链接
	//defer db.Close()
	return err
}

func MysqlCreateTable() {
	//OpenMysql()
	//defer db.Close()
	_, err := db.Exec("SELECT COUNT(*) FROM rules")
	if err != nil {
		db.Exec("CREATE TABLE rules(id int UNSIGNED AUTO_INCREMENT,aid INT UNSIGNED,hit_count INT UNSIGNED DEFAULT 0,download_count INT UNSIGNED DEFAULT 0,  platform CHAR(16),download_url VARCHAR(128),update_version_code	VARCHAR(128),device_list TEXT,md5	VARCHAR(128),max_update_version_code	VARCHAR(128),min_update_version_code	VARCHAR(128),max_os_api	TINYINT UNSIGNED,min_os_api	TINYINT UNSIGNED,cpu_arch	TINYINT UNSIGNED,channel	VARCHAR(128),title	VARCHAR(256),update_tips	VARCHAR(1024),enabled	BOOLEAN DEFAULT true,create_date DATETIME DEFAULT CURRENT_TIMESTAMP,PRIMARY KEY ( id ));")
	}
	// db.Exec("CREATE TABLE rules(id int UNSIGNED AUTO_INCREMENT,aid INT UNSIGNED,hit_count INT UNSIGNED DEFAULT 0,download_count INT UNSIGNED DEFAULT 0,  platform CHAR(16),download_url VARCHAR(128),update_version_code	VARCHAR(128),device_list TEXT,md5	VARCHAR(128),max_update_version_code	VARCHAR(128),min_update_version_code	VARCHAR(128),max_os_api	TINYINT UNSIGNED,min_os_api	TINYINT UNSIGNED,cpu_arch	TINYINT UNSIGNED,channel	VARCHAR(128),title	VARCHAR(256),update_tips	VARCHAR(1024),enabled	BOOLEAN DEFAULT true,create_date DATETIME DEFAULT CURRENT_TIMESTAMP,PRIMARY KEY ( id ));")
}

func MysqlAddRule(rulemap *map[string]string, devicelst *[]string) (int64, error) {
	//OpenMysql()
	//defer db.Close()
	var devices string
	var res sql.Result
	var err error
	if devicelst != nil {
		devices = ltos(devicelst)
	} else {
		devices = ""
	}

	res, err = db.Exec("insert into rules(aid,platform,download_url,update_version_code,device_list,md5,max_update_version_code,min_update_version_code,max_os_api,min_os_api,cpu_arch,channel,title,update_tips) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)", (*rulemap)["aid"], (*rulemap)["platform"], (*rulemap)["download_url"], (*rulemap)["update_version_code"], devices, (*rulemap)["md5"], (*rulemap)["max_update_version_code"], (*rulemap)["min_update_version_code"], (*rulemap)["max_os_api"], (*rulemap)["min_os_api"], (*rulemap)["cpu_arch"], (*rulemap)["channel"], (*rulemap)["title"], (*rulemap)["update_tips"])
	if err != nil {
		tools.LogMsg(err)
		panic(err)
	}
	val, _ := res.LastInsertId()
	//fmt.Printf("res: %v\n", val)
	return val, err
}

func MysqlUpdateRule(rulemap *map[string]string, devicelst *[]string) error {
	//OpenMysql()
	//defer db.Close()
	var devices string
	if devicelst != nil {
		devices = ltos(devicelst)
	}
	id := (*rulemap)["id"]
	delete(*rulemap, "id")
	if id == "" {
		return errors.New("id can't be none")
	}
	var err error
	if devicelst != nil {
		_, err = db.Exec("update rules set device_list=? where id=?", devices, id)
	}
	checkErr(err)
	for key, val := range *rulemap {
		s := fmt.Sprintf("update rules set %s=? where id=?", key)
		_, err := db.Exec(s, val, id)
		checkErr(err)
	}
	return err
}

func MysqlDeleteRule(ruleid string) error {
	//OpenMysql()
	//defer db.Close()
	_, err := db.Exec("delete from rules where id=?", ruleid)
	checkErr(err)
	return err
}
func GetDownloadRatio(ruleid string) (float64, error) {
	//OpenMysql()
	//defer db.Close()
	qres, _, err := MysqlQueryRules(ruleid)
	checkErr(err)
	res := (*qres)[0]
	return float64(ToInt(res["download_count"])) / float64(ToInt(res["hit_count"])), err
}

func RowsToMap(rows *sql.Rows) *[]map[string]string {
	result := make([]map[string]string, 0)

	for rows.Next() {
		cols, err := rows.Columns()
		checkErr(err)

		colsTypes, err := rows.ColumnTypes()
		checkErr(err)

		dest := make([]interface{}, len(cols))
		destPointer := make([]interface{}, len(cols))
		for i := range dest {
			destPointer[i] = &dest[i]
		}

		err = rows.Scan(destPointer...)
		checkErr(err)

		rowResult := make(map[string]interface{})
		for i, colVal := range dest {
			colName := cols[i]
			itemType := colsTypes[i].ScanType()
			//fmt.Printf("type %v \n", itemType)

			switch itemType.Kind() {
			case BytesKind:
				rowResult[colName] = ToStr(colVal)

			case reflect.Int, reflect.Int8,
				reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:

				rowResult[colName] = ToInt(colVal)

			case TimeKind:
				// fmt.Println("time", colVal, reflect.TypeOf(colVal))
				rowResult[colName] = ToStr(colVal)
			default:
				rowResult[colName] = ToStr(colVal)
			}
		}
		result = append(result, itos(&rowResult))
	}
	return &result
}

func itos(m *map[string]interface{}) map[string]string {
	ret := make(map[string]string, len(*m))
	for k, v := range *m {
		ret[k] = fmt.Sprint(v)
	}
	return ret
}
func ltos(sl *[]string) string {
	s := ""
	for _, val := range *sl {
		s += val + ","
	}
	return s
}

//根据id查询规则，"0"代表全部
func MysqlQueryRules(ruleid string) (*[]map[string]string, *[]string, error) {
	//OpenMysql()
	//defer db.Close()
	if ruleid == "0" {
		dbrows, err := db.Query("select * from rules")
		if err != nil {
			tools.LogMsg(err)
			panic(err)
			// return nil, err
		}
		s := make([]string, 0)
		return RowsToMap(dbrows), &s, err
	} else {
		dbrows, err := db.Query("select * from rules where id=?", ruleid)
		if err != nil {
			tools.LogMsg(err)
			panic(err)
			// return nil, err
		}
		rowsmap := RowsToMap(dbrows)
		s := strings.Split((*rowsmap)[0]["device_list"], ",")
		return rowsmap, &s, err
	}
}

var AddLst []int32 = make([]int32, 1000)
var HitLst []int32 = make([]int32, 1000)

func AddHitCnt(c UpdCnt) {
	i, _ := strconv.Atoi(c.Id)
	atomic.AddInt32(&HitLst[i], 1)
	if c.Success {
		atomic.AddInt32(&AddLst[i], 1)
	}
}

func CommitAll() {
	tx, err := db.Begin()
	if err != nil {
		tools.LogMsg(err)
		return
	}
	for i, v := range AddLst {
		if v > 0 {
			qStr := "UPDATE rules SET download_count = download_count + ? WHERE id = ?"
			_, err := tx.Exec(qStr, v, i)
			if err != nil {
				_ = tx.Rollback()
				tools.LogMsg(err)
			}
		}
		atomic.StoreInt32(&AddLst[i], 0)
	}
	for i, v := range HitLst {
		if v > 0 {
			qStr := "UPDATE rules SET hit_count = hit_count + ? WHERE id = ?"
			_, err := tx.Exec(qStr, v, i)
			if err != nil {
				_ = tx.Rollback()
				tools.LogMsg(err)
			}
		}
		atomic.StoreInt32(&HitLst[i], 0)
	}
	e := tx.Commit()
	if e != nil {
		tools.LogMsg(err)
	}
}
