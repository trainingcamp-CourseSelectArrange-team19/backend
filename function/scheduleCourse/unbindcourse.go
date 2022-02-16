package scheduleCourse

import (
	"backend/database"
	"backend/types"
	"github.com/gin-gonic/gin"
	"strconv"
)

/*
// 老师解绑课程
// Method： Post
type UnbindCourseRequest struct {
	CourseID  string
	TeacherID string
}

type UnbindCourseResponse struct {
	Code ErrNo
}
 */
func UnBindCourse(c *gin.Context) {
	b := types.UnbindCourseResponse{Code: types.ParamInvalid}
	var arg types.UnbindCourseRequest
	if err := c.ShouldBindJSON(&arg); err != nil {
		c.JSON(200,b)
		return
	}
	courseid,teacherid := arg.CourseID,arg.TeacherID
	cid,_ := strconv.Atoi(courseid)
	tid,_ := strconv.Atoi(teacherid)
	courseExisted, _ := database.GetOneCourseName(cid)
	if  courseExisted != ""{
		b.Code = types.CourseNotExisted
		c.JSON(200,b)
		return
	}
	t, schedule := database.GetCourseTeacher(cid)
	if t == ""{
		if schedule.TeacherID == 0{
			b.Code = types.CourseNotBind
			c.JSON(200,b)
			return
		}
	}
	flag := database.UnbindTeacherCourse(tid, cid)
	if flag != "unbind successes!" {
		b.Code = types.UnknownError
		c.JSON(200,b)
		return
	}
	b.Code = types.OK
	c.JSON(200,b)
}