package main


type sqlStr struct {
	s string 
}

func PreExecuteInsert(table string, length int) *sqlStr {
	res := "INSERT INTO "
	res += table
	res += " "
	res += "VALUES"
	if (length == 1) {
		res += "?"
	} else {
		res += "("
		for i := 0; i < length; i++ {
			if (i > 0) {
				res += ", "
			}
			res += "?"
		}
		res += ")"
	}
	t := sqlStr{s:res}
	return &t;
}

func PreExecuteSelect(table string, tabs ...string) *sqlStr {
	res := "SELECT"
	for idx, i := range tabs {
		res += " "
		res += i
		if (idx != len(tabs) - 1) {
			res += ","
		}
	}
	res += " "
	res += "FROM "
	res += table;
	t := sqlStr{s:res}
	return &t
}
func (s *sqlStr) where(info string) *sqlStr {
	s.s += " WHERE"
	s.s += " "
	s.s += info
	return s
}