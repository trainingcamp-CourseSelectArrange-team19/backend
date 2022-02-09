package selectCourse

import (
	"backend/database"
	"backend/tools"
	"backend/types"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"strconv"
)
var (
	redisPool *redis.Pool
	userSetCount int
)
func InitRedisConfig() {
	redisPool = NewPool()
	userSetCount = 0
	_, err, users := database.GetAllValStudentInfo()
	if err != nil{
		tools.LogMsg(err)
		panic(err)
	}
	_, courses := database.GetAllCourse()
	if err != nil{
		tools.LogMsg(err)
		panic(err)
	}
	redisConn := redisPool.Get()
	defer redisConn.Close()
	val, err1 := redis.Strings(redisConn.Do("KEYS", "*"))
	if err1 != nil {
		tools.LogMsg(err1)
		panic(err1)
	}
	redisConn.Send("MULTI")
	for i, _ := range val{
		redisConn.Send("DEL", val[i])
	}
	redisConn.Do("EXEC")
	for ind := 0 ; ind < len(users) ; ind++{
		whichPart := ind / 5000
		if whichPart > userSetCount{
			userSetCount = whichPart
		}
		_, err := redis.Int64(redisConn.Do("SADD", "studentsID:" + strconv.Itoa(whichPart), strconv.Itoa(users[ind].Id)))
		if err != nil {
			tools.LogMsg(err)
			panic(err)
		}
	}
	for ind := 0 ; ind < len(courses) ; ind++{
		_, err := redisConn.Do("SET", "seckill:" + strconv.Itoa(courses[ind].Id) + ":stock", courses[ind].Capacity)
		if err != nil {
			tools.LogMsg(err)
			panic(err)
		}
		_, err = redisConn.Do("SET", "seckill:" + strconv.Itoa(courses[ind].Id) + ":end", 0)
		if err != nil {
			tools.LogMsg(err)
			panic(err)
		}
	}
}
//处理请求函数,根据请求将响应结果信息写入日志
func SelectCourse(c *gin.Context) {
	b := types.BookCourseResponse{
		Code: types.ParamInvalid,
	}
	var arg types.BookCourseRequest
	if err := c.ShouldBind(&arg); err != nil {
		c.JSON(500, b)
		panic(err.Error())
		return
	}
	StudentID, CourseID := arg.StudentID, arg.CourseID
	redisConn := redisPool.Get()
	defer redisConn.Close()
	courseExist, err := redisConn.Do("GET", "seckill:" + CourseID + ":end")
	if err != nil {
		tools.LogMsg(err)
		panic(err)
	}
	if courseExist == nil{
		b.Code = types.CourseNotExisted
		c.JSON(200, b)
		return
	}
	found := false
	for ind := 0 ; ind <= userSetCount ; ind++{
		studentExist, err := redis.Int64(redisConn.Do("sismember", "studentsID:" + strconv.Itoa(ind), StudentID))
		if err != nil {
			tools.LogMsg(err)
			panic(err)
		}
		if studentExist > 0{
			found = true
			break
		}
	}
	if !found {
		b.Code = types.StudentNotExisted
		c.JSON(200, b)
		return
	}
	success := RemoteDeductionStock(redisConn, CourseID, StudentID)
	if success {
		b.Code = types.OK
		c.JSON(200, b)
	} else {
		b.Code = types.CourseNotAvailable
		c.JSON(200, b)
	}
}
