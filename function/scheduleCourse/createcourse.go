package scheduleCourse

import (
	"backend/database"
	"backend/function/selectCourse"
	"backend/tools"
	"backend/types"
	"github.com/gin-gonic/gin"
	"strconv"
)


/*
type CreateCourseRequest struct {
	Name string
	Cap  int
}

type CreateCourseResponse struct {
	Code ErrNo
	Data struct {
		CourseID string
	}
}
 */
func Createcourse(c *gin.Context)  {
	b := types.CreateCourseResponse{Code: types.ParamInvalid}
	var arg types.CreateCourseRequest
	if err := c.ShouldBind(&arg);err != nil {
		c.JSON(200,b)
		panic(err.Error())
		return
	}
	Name,Cap := arg.Name,arg.Cap

	if database.CreateCourse(Name,Cap) != "create successes!"{
		b.Code = types.UnknownError
		c.JSON(200,b)
		return
	}
	redisConn := selectCourse.RedisPool.Get()
	defer redisConn.Close()
	b.Code = types.OK
	_,course := database.GetOneCourse(Name)
	b.Data = struct{ CourseID string }{CourseID:strconv.Itoa(course.Id)}
	_, err := redisConn.Do("SET", "seckill:" + strconv.Itoa(course.Id) + ":stock", course.Capacity)
	if err != nil {
		tools.LogMsg(err)
		panic(err)
	}
	_, err = redisConn.Do("SET", "seckill:" + strconv.Itoa(course.Id) + ":end", 0)
	if err != nil {
		tools.LogMsg(err)
		panic(err)
	}
	c.JSON(200,b)
}