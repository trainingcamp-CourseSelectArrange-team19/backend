package main

import (
	"database/sql"
	"fmt"
	"time"
	"os"
	_ "github.com/go-sql-driver/mysql"
)

const SQL_URL = "root:bytedancecamp@tcp(180.184.65.192:3306)/test"

/* 包装数据库 */
var TerminateSig chan string
type db_t struct {
	dbi *sql.DB
	log []string
}

/* 全局数据库对象 */
var db *db_t
func (db *db_t) Debug() {
	for _, s := range(db.log) {
		fmt.Printf("%s\n", s)
	}
}

/* 获取数据库实例 */
func UseDB() *db_t {
	if (db != nil) {
		return db
	} 
	TerminateSig = make(chan string)
	go func() {
		select {
		case v := <-TerminateSig:
			fmt.Printf("%s", v)
			os.Exit(1)
		case <- time.After(3 * time.Second):
			fmt.Printf("test")
		}
	}()
	var err error
	db := new(db_t)
	db.dbi, err = sql.Open("mysql", SQL_URL) //用户名:密码@/数据库名
	if (err != nil) {
		panic("wtf?! no db created")
	}
	return db
}


/* 
	@ name 需要详细内容 如果参数不是全部信息 需要分别列出  eg. user(name, age)
	@ 防注入未测试
	@ 1/30 fix:返回数据库对象 可链式调用
 */
func (db *db_t) Insert(table string, value ...interface{}) *db_t {
	sqlStr := PreExecuteInsert(table, len(value)).s
	//fmt.Printf("test sqlStr result: %s",sqlStr)
	ret, err := db.dbi.Exec(sqlStr, value...)
	if err != nil {
		v := fmt.Sprintf("insert failed, err:%v\n", err)
		TerminateSig<-v
		return db
	}
	theID, err := ret.LastInsertId() // 新插入数据的id
	if err != nil {
		v := fmt.Sprintf("get last insert ID failed, err:%v\n", err)
		TerminateSig<-v
		return db
	}
	v := fmt.Sprintf("insert success, the id is %d.\n", theID)
	db.log = append(db.log, v)
	return db
}

/* 
	测试用 直接go run mysql.go
*/
func main() {
	db := UseDB()
	k := PreExecuteSelect("user","name","age").where("name > 10 AND age < 10").s
	db.Debug()
	fmt.Printf("%s", k)


}  