package scheduleCourse

import (
	"backend/database"
	"backend/types"
	"github.com/gin-gonic/gin"
	"strconv"
)


/*
type GetTeacherCourseRequest struct {
	TeacherID string
}

type GetTeacherCourseResponse struct {
	Code ErrNo
	Data struct {
		CourseList []*TCourse
	}
}
type TCourse struct {
	CourseID string
	Name     string
	TeacherID string
}
 */
func GetTeacherCourse(c *gin.Context){
	b := types.GetTeacherCourseResponse{Code: types.ParamInvalid}
	var arg types.GetTeacherCourseRequest
	if err := c.ShouldBind(&arg); err != nil {
		c.JSON(200,b)
		return
	}
	TeacherID,_ := strconv.Atoi(arg.TeacherID)

	n,err,arr := database.GetTeacherCourse(TeacherID)
	if err != nil {
		b.Code = types.UnknownError
		c.JSON(200,b)
		return
	}
	tcourse := make([]*types.TCourse,n)
	for i := 0; i < int(n); i++ {
		tcourse[i].TeacherID = strconv.Itoa(TeacherID)
		tcourse[i].CourseID = strconv.Itoa(arr[i].CourseID)
		_,tempCourse := database.GetOneCourseName(arr[i].CourseID)
		tcourse[i].Name = tempCourse.Name
	}

	b.Data = struct{ CourseList []*types.TCourse }{CourseList: tcourse}
	b.Code = types.OK
	c.JSON(200,b)

}