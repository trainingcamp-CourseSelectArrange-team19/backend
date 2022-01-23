package database

import "testing"

func TestGetRuleAttr(t *testing.T) {

	testCases := []struct {
		input  string
		expect string
	}{
		{input: "aid", expect: "8"},
		{input: "download_url", expect: "https://baidu1.com"},
	}

	qObj := RuleObj{}

	for _, v := range testCases {
		s, _ := qObj.GetRuleAtt(v.input, v.input)
		if v.expect != s {
			t.Error("Unexp!")
		}
	}

}
