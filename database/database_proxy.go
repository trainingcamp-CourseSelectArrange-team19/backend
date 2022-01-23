package database

import "strings"

func qrnGetSObj() (*[]map[string]string, *[]string, error) {
	mp := []map[string]string{{
		"min_update_version_code": "8.4.0",
		"max_update_version_code": "8.8.8",
		"min_os_api":              "10",
		"max_os_api":              "20",
		"platform":                "Android",
		"cpu_arch":                "32",
		"channel":                 "dsd",
		"download_url":            "https://baidu1.com",
		"update_version_code":     "4.1",
		"md5":                     "aaa",
		"title":                   "asd",
		"update_tips":             "sad",
		"aid":                     "8",
		"enabled":                 "true",
	}}
	lst := []string{"aa", "bb"}
	return &mp, &lst, nil
}

type QueryObj interface {
	GetRuleAtt(ruleid string, field string) (string, error)
	CheckDeviceIDInWhiteList(ruleid string, userid string) (bool, error)
}

type RuleObj struct {
	initialed bool
	oldrule   string
	Rule      *map[string]string
	White     *[]string
}

func (r *RuleObj) GetRuleAtt(ruleid string, field string) (string, error) {
	if !r.initialed || strings.Compare(r.oldrule, ruleid) != 0 {
		e := r.InitRuleObj(ruleid)
		if e != nil {
			return "", e
		}
		r.initialed = true
		r.oldrule = ruleid
	}
	return (*r.Rule)[field], nil
}

func (r *RuleObj) InitRuleObj(ruleid string) error {
	a, w, e := QueryRuleByID(ruleid)
	// a, w, e := qrnGetSObj()
	if e != nil {
		return e
	}
	r.Rule = &(*a)[0]
	r.White = w
	return nil
}

func (r *RuleObj) CheckDeviceIDInWhiteList(ruleid string, userid string) (bool, error) {
	return CheckDeviceIDInWhiteList(ruleid, userid)
}
