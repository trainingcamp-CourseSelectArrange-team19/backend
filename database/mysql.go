package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
func DBStart() {
	var err error
	db, err = sql.Open("mysql", "root:bytedancecamp@tcp(180.184.65.192:3306)/test") //用户名:密码@/数据库名
	if err != nil {
		fmt.Printf("bad")
	}
}

func DBEnd() {
	db.Close()
}


func preExecuteInsert(name string, length int) string {
	//生成 SQL 语句
	res := "INSERT INTO "
	res += name
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
	return res;
}



/* 
	@ name 需要详细内容 如果参数不是全部信息 需要分别列出 
	@ eg. user(name, age)
	@ 防注入未测试
 */
func Insert(name string, value ...interface{}) {
	sqlStr := preExecuteInsert(name, len(value))
	fmt.Printf("%s",sqlStr)
	ret, err := db.Exec(sqlStr, value...)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	theID, err := ret.LastInsertId() // 新插入数据的id
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err:%v\n", err)
		return
	}
	fmt.Printf("insert success, the id is %d.\n", theID)

}


func main() {
	DBStart()


	Insert("user", 5, "fuck", 24)


	DBEnd()
}