package database

import (
	"fmt"
	"testing"
)

func TestAddRule(t *testing.T) {
	r := map[string]string{
		"aid":                     "133",
		"platform":                "iOS",
		"download_count":          "0",
		"hit_count":               "0",
		"download_url":            "http://baidu.com",
		"update_version_code":     "1.1.1",
		"device_list":             "1,2,3,4",
		"md5":                     "1233",
		"max_update_version_code": "1.1.0",
		"min_update_version_code": "1.0.0",
		"max_os_api":              "0",
		"min_os_api":              "0",
		"cpu_arch":                "32",
		"channel":                 "App Store",
		"title":                   "Update",
		"update_tips":             "yes",
		"enabled":                 "1",
		"create_date":             "123",
	}
	dl := []string{"1", "2", "3"}
	err := AddRule(&r, &dl)
	if err != nil {
		t.Error(err)
	}
}
func TestUpdateRule(t *testing.T) {
	r := map[string]string{
		"id":                      "1",
		"aid":                     "233",
		"platform":                "iOS",
		"download_count":          "0",
		"hit_count":               "0",
		"download_url":            "http://baidu.com",
		"update_version_code":     "1.1.1",
		"device_list":             "1,2,3,",
		"md5":                     "1231",
		"max_update_version_code": "1.1.0",
		"min_update_version_code": "1.0.0",
		"max_os_api":              "0",
		"min_os_api":              "0",
		"cpu_arch":                "32",
		"channel":                 "App Store",
		"title":                   "Update",
		"update_tips":             "no",
		"enabled":                 "1",
		"create_date":             "123",
	}
	dl := []string{"3", "2", "1"}
	err := UpdateRule(&r, &dl)
	if err != nil {
		t.Error(err)
	}
}
func TestDeleteRule(t *testing.T) {
	DeleteRule("1")
}
func TestQueryAllRules(t *testing.T) {
	fmt.Println(QueryAllRules())
}

func TestQueryRuleByID(t *testing.T) {
	fmt.Println(RedisGetAllKeys())
	fmt.Println(QueryRuleByID("1"))
}

func TestGetRuleAtt(t *testing.T) {
	fmt.Println(GetRuleAtt("1", "aid"))
}

func TestCheckDeviceIDInWhiteList(t *testing.T) {
	fmt.Println(CheckDeviceIDInWhiteList("1", "4"))
}

func TestGetIDList(t *testing.T) {
	fmt.Println(GetIDList())
}

func TestRedisDeleteAll(t *testing.T) {
	RedisDeleteAll()
	fmt.Println(RedisGetAllKeys())
}

func TestRedisGetAllKeys(t *testing.T) {
	fmt.Println(RedisGetAllKeys())
}
