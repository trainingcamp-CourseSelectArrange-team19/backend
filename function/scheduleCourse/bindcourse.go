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
	courseExisted, _ := database.GetOneCourseName(cid)
	if  courseExisted != ""{
		b.Code = types.CourseNotExisted
		c.JSON(200,b)
		return
	}
	t, schedule := database.GetCourseTeacher(cid)
	if t == ""{
		if schedule.TeacherID != 0{
			b.Code = types.CourseHasBound
			c.JSON(200,b)
			return
		}
	}
	flag := database.BindTeacherCourse(tid,cid)
	if flag != "bind successes!" {
		b.Code = types.UnknownError
		c.JSON(200,b)
		return
	}
	b.Code = types.OK
	c.JSON(200,b)
}