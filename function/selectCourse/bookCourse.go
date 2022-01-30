package selectCourse

import (
	"backend/types"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)
var (
	redisPool *redis.Pool
)
func InitRedisConfig() {
	redisPool = NewPool()
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
	//加ID合法性校验
	redisConn := redisPool.Get()
	defer redisConn.Close()
	success := RemoteDeductionStock(redisConn, CourseID, StudentID)
	//可以考虑加读写锁
	if success {
		b.Code = types.OK
		c.JSON(200, b)
	} else {
		b.Code = types.CourseNotAvailable
		c.JSON(500, b)
	}
}
