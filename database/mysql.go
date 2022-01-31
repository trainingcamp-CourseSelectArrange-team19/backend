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
		}
	}()

	var err error
	db := new(db_t)
	db.dbi, err = sql.Open("mysql", SQL_URL) //用户名:密码@/数据库名
	if (err != nil) {
		panic("wtf?!")
	}
	return db
}

/* 
	@检测表中是否有相关限制下的数据
*/
func (db *db_t) Find(table string, limit string) bool {
	str := PreExecuteSelect("user", "count(*)").where(limit).g()
	cnt := 0
	err := db.dbi.QueryRow(str).Scan(&cnt)
	if (err != nil) {
		v := fmt.Sprintf("[[FIND FAILED]] %v\n", err)
		TerminateSig <- v
	}
	return cnt > 0
}

/* 
	@ name 需要详细内容 如果参数不是全部信息 需要分别列出  eg. user(name, age)
	@ 防注入未测试
	@ 1/30 fix:返回数据库对象 可链式调用
 */
func (db *db_t) Insert(table string, value ...interface{}) *db_t {
	str := PreExecuteInsert(table, len(value)).g()
	ret, err := db.dbi.Exec(str, value...)
	if err != nil {
		v := fmt.Sprintf("[[INSERT FAILED]] %v\n", err)
		TerminateSig <- v
		return db
	}
	theID, err := ret.LastInsertId() // 新插入数据的id
	if err != nil {
		v := fmt.Sprintf("get last insert ID failed, err:%v\n", err)
		TerminateSig <- v
		return db
	}
	v := fmt.Sprintf("insert success, the id is %d.\n", theID)
	db.log = append(db.log, v)
	return db
}

/* 
	@注意limit在第二个位置
*/
func (db *db_t) Update(table string, limit string, tabs ...string) *db_t {
	str := PreExecuteUpdate(table, tabs...).where(limit).g()
	ret, err := db.dbi.Exec(str)
	if err != nil {
		v := fmt.Sprintf("[[UPDATE FAILED]] %v\n", err)
		TerminateSig <- v
		return db
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		v := fmt.Sprintf("get RowsAffected failed, err:%v\n", err)
		TerminateSig <- v
		return db
	}
	if (n != 1) {
		v := "[[WARNING]] update multiple rows while it should be one\n"
		TerminateSig <- v
	}
	return db
}

/* test */
type user struct {
	id   int
	age  int
	name string
}
func (db *db_t) Get(table string, limit string, tabs ...string) []string {
	str := PreExecuteSelect(table, tabs...).where(limit).g()
	rows, err := db.dbi.Query(str)
	var res []string
	if err != nil {
		v := fmt.Sprintf("[[QUERY FAILED]]: %v\n", err)
		TerminateSig <- v
		return res
	}
	defer rows.Close()

	for rows.Next() {
		var u user
		err := rows.Scan(&u.id, &u.name, &u.age)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return res
		}
		res = append(res, fmt.Sprintf("id:%d name:%s age:%d\n", u.id, u.name, u.age))
	}
	return res
}
/* 
	测试用 直接go run mysql.go
*/
func main() {
	db := UseDB()
	v := db.Get("user", "id > 5")
	for _, i := range(v) {
		fmt.Printf("%s", i)
	}
	time.Sleep(time.Duration(10) * time.Second)


}  