package scheduleCourse

import (
	"backend/database"
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
		return
	}
	Name,Cap := arg.Name,arg.Cap

	if database.CreateCourse(Name,Cap) != "create successes!"{
		b.Code = types.UnknownError
		c.JSON(200,b)
		return
	}
	b.Code = types.OK
	_,course := database.GetOneCourse(Name)
	b.Data = struct{ CourseID string }{CourseID:strconv.Itoa(course.Id)}
	c.JSON(200,b)
}