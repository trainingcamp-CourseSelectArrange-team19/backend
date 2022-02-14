package scheduleCourse

import (
	"backend/database"
	"backend/types"
	"github.com/gin-gonic/gin"
	"strconv"
)


/*
type BindCourseRequest struct {
	CourseID  string
	TeacherID string
}

type BindCourseResponse struct {
	Code ErrNo
}
 */
func BindCourse(c *gin.Context) {
	b := types.BindCourseResponse{Code: types.ParamInvalid}
	var arg types.BindCourseRequest
	if err := c.ShouldBind(&arg);err != nil{
		c.JSON(200,b)
		return
	}
	courseID,TeacherID := arg.CourseID,arg.TeacherID
	cid,_ := strconv.Atoi(courseID)
	tid,_ := strconv.Atoi(TeacherID)
	flag := database.BindTeacherCourse(tid,cid)
	if flag != "bind successes!" {
		b.Code = types.CourseNotExisted
		c.JSON(200,b)
		return
	}
	b.Code = types.CourseHasBound
	c.JSON(200,b)
}