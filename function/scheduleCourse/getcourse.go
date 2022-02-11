package scheduleCourse

import (
	"backend/database"
	"backend/types"
	"github.com/gin-gonic/gin"
	"strconv"
)


/*
// 获取课程
// Method: Get
type GetCourseRequest struct {
	CourseID string
}

type GetCourseResponse struct {
	Code ErrNo
	Data TCourse
}
CourseNotExisted   ErrNo = 12   // 课程不存在
OK                 ErrNo = 0
	ParamInvalid       ErrNo = 1  // 参数不合法
*/

func getcourse(c *gin.Context)  {
	b := types.GetCourseResponse{Code: types.ParamInvalid}

	var arg types.GetCourseRequest
	if err := c.ShouldBind(&arg); err != nil {
		c.JSON(200,b)
		return
	}
	courseID := arg.CourseID
	i,err := strconv.Atoi(courseID)
	if err != nil {
		c.JSON(200,b)
		return
	}
	// 不存在
	flag,u := database.GetOneCourseName(i)
	if flag != ""{
		b.Code = types.CourseNotExisted
		c.JSON(200,b)
		return
	}
	flag,r := database.GetTeacherFromCourse(i)
	b.Code= types.OK
	b.Data = types.TCourse{
		CourseID: courseID,
		Name:u.Name,
		TeacherID: strconv.Itoa(r.TeacherID),
	}
	c.JSON(200,b)
}
