package database

import (
	"database/sql"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"backend/tools"

	"github.com/go-sql-driver/mysql"
)

var BytesKind = reflect.TypeOf(sql.RawBytes{}).Kind()
var TimeKind = reflect.TypeOf(mysql.NullTime{}).Kind()
var timecnt int64 = 0

const UPDATETIME = 5

func checkErr(err error) {
	if err != nil {
		tools.LogfMsg("checkErr:%v\n", err)
		// fmt.Printf("checkErr:%v\n", err)
	}
}

func ToStr(strObj interface{}) string {
	switch v := strObj.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", strObj)
	}
}

func ToInt(intObj interface{}) int {
	// 假定int == int64，运行在64位机
	switch v := intObj.(type) {
	case []byte:
		return ToInt(string(v))
	case int:
		return v
	case int8:
		return int(v)
	case int16:
		return int(v)
	case int32:
		return int(v)
	case int64:
		return int(v)
	case uint:
		return int(v)
	case uint8:
		return int(v)
	case uint16:
		return int(v)
	case uint32:
		return int(v)
	case uint64:
		if v > math.MaxInt64 {
			info := fmt.Sprintf("ToInt, error, overflowd %v", v)
			tools.LogfMsg("ToInt, error, overflowd %v", v)
			panic(info)
		}
		return int(v)
	case float32:
		return int(v)
	case float64:
		return int(v)
	case string:
		strv := v
		if strings.Contains(v, ".") {
			strv = strings.Split(v, ".")[0]
		}
		if strv == "" {
			return 0
		}
		if intv, err := strconv.Atoi(strv); err == nil {
			return intv
		}
	}
	// fmt.Printf(fmt.Sprintf("ToInt err, %v, %v not supportted\n", intObj,
	// reflect.TypeOf(intObj).Kind()))
	tools.LogfMsg("ToInt err, %v, %v not supportted\n", intObj,
		reflect.TypeOf(intObj).Kind())
	return 0
}

func CheckDeviceIDInWhiteList(ruleid string, userid string) (bool, error) {
	res, err := RedisCheckWhiteList(ruleid, userid)
	if err != nil {
		qres, wls, err2 := MysqlQueryRules(ruleid)
		if err2 != nil {
			return false, err2
		}
		RedisUpdateRule(ruleid, &(*qres)[0], wls)
		res, err = RedisCheckWhiteList(ruleid, userid)
	} else {
		return res, err
	}
	return res, err
}

func GetRuleAtt(ruleid string, field string) (string, error) {
	val, err := RedisGetRuleAttr(ruleid, field)
	if err != nil || val == "" {
		qres, wls, err2 := MysqlQueryRules(ruleid)
		if err2 != nil || len(*qres) == 0 {
			return "Not Match!", err2
		}
		RedisUpdateRule(ruleid, &(*qres)[0], wls)
		val, err = RedisGetRuleAttr(ruleid, field)
	} else {
		return val, err
	}
	return val, err
}

func UpdateUserDownloadStatus(ruleid string, status bool) error {
	err := RedisUpdateDownloadStatus(ruleid, status)
	checkErr(err)
	AddHitCnt(UpdCnt{ruleid, status})
	// if time.Now().Unix()-timecnt > UPDATETIME {
	// 	timecnt = time.Now().Unix()
	// 	val, wls, _ := RedisQueryRuleByID(ruleid)
	// 	(*val)[0]["id"] = ruleid
	// 	MysqlUpdateRule(&(*val)[0], wls)
	// }
	return err
}

//查询所有规则，为了保证完整性，对 mysql 查询
func QueryAllRules() (*[]map[string]string, error) {
	val, _, err := MysqlQueryRules("0")
	return val, err
}

//优先对 redis 查询，若没查询到，对 mysql 查询并更新 redis
func QueryRuleByID(ruleid string) (*[]map[string]string, *[]string, error) {
	res, devices, err := RedisQueryRuleByID(ruleid)
	if err != nil || len(*res) == 0 {
		// fmt.Println(res)
		// tools.LogMsg(err)
	} else {
		return res, devices, err
	}
	// fmt.Println("Redis not found, query mysql next...")
	res, devices, err = MysqlQueryRules(ruleid)
	if err != nil || len(*res) == 0 {
		// fmt.Println("Wrong ID!")
		tools.LogMsg("Wrong ID!")
		return res, devices, err
	}
	RedisUpdateRuleWithList(ruleid, &(*res)[0])
	return res, devices, err
}

//提供一个 string-string 的哈希表和白名单，向 mysql 添加规则。
func AddRule(rulemap *map[string]string, devicelst *[]string) error {
	// fmt.Println(rulemap, devicelst)
	id, err := MysqlAddRule(rulemap, devicelst)
	checkErr(err)
	// fmt.Printf("!!")
	// fmt.Println(id)
	err = RedisUpdateRule(ToStr(id), rulemap, devicelst)
	checkErr(err)
	return err
}

func UpdateRule(rulemap *map[string]string, devicelst *[]string) error {
	// r, _, _ := QueryRuleByID((*rulemap)["id"])
	// for ky, v := range *rulemap {
	// 	(*r)[0][ky] = v
	// }
	// if tools.JudgeLegalRule(rulemap) == false {
	// 	return errors.New("Rule is not legal!")
	// }
	// fmt.Println(rulemap, devicelst)
	err := RedisUpdateRule((*rulemap)["id"], rulemap, devicelst)
	checkErr(err)
	err = MysqlUpdateRule(rulemap, devicelst)
	checkErr(err)
	return err
}

func DeleteRule(ruleid string) error {
	err := MysqlDeleteRule(ruleid)
	checkErr(err)
	err = RedisDeleteRule(ruleid)
	checkErr(err)
	return err
}

// 这个接口直接放在了 mysql.go 中
// func GetDownloadRatio(ruleid string) (float64, error)

// 这个接口直接放在了 redis.go 中
// func GetIDList()(*[]string,error)
