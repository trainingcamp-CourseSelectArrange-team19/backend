package database

// import (
// 	"strconv"
// 	"testing"
// )

// func TestAdd(t *testing.T) {
// 	RedisInitClient()
// 	r := map[string]string{
// 		"aid":                     "123",
// 		"platform":                "iOS",
// 		"download_count":          "0",
// 		"hit_count":               "0",
// 		"download_url":            "http://baidu.com",
// 		"update_version_code":     "1.1.1",
// 		"device_list":             "1,2,3,",
// 		"md5":                     "123",
// 		"max_update_version_code": "1.1.0",
// 		"min_update_version_code": "1.0.0",
// 		"max_os_api":              "0",
// 		"min_os_api":              "0",
// 		"cpu_arch":                "32",
// 		"channel":                 "App Store",
// 		"title":                   "Update",
// 		"update_tips":             "yes",
// 		"enabled":                 "1",
// 		"create_date":             "123",
// 	}
// 	tmp_id := cur_id
// 	RedisUpdateRuleWithList("1", r)
// 	//fmt.Println(err)

// 	val, _ := RedisCheckWhiteList(strconv.Itoa(tmp_id)+"s", "1")
// 	if val == false {
// 		t.Errorf("UnExpected!")
// 	}

// }

// func TestCheckAppidInWhiteList(t *testing.T) {
// 	RedisInitClient()
// 	one, _ := RedisCheckWhiteList(strconv.Itoa(cur_id-1)+"s", "1")

// 	two, _ := RedisCheckWhiteList(strconv.Itoa(cur_id-1)+"s", "2")

// 	three, _ := RedisCheckWhiteList(strconv.Itoa(cur_id-1)+"s", "3")
// 	four, _ := RedisCheckWhiteList(strconv.Itoa(cur_id-1)+"s", "4")
// 	if four == true {
// 		t.Error("Four UnExpected!")
// 	}
// 	if one == false || two == false || three == false {
// 		t.Errorf("WhiteList is not correct!%v %v %v", one, two, three)
// 	}
// }
