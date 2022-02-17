package scheduleCourse

import (
	"backend/database"
	"backend/types"
	"github.com/gin-gonic/gin"
	"net/http"
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
func GetTeacherCourse(c *gin.Context) {
	b := types.GetTeacherCourseResponse{Code: types.ParamInvalid}
	var arg types.GetTeacherCourseRequest
	if err := c.ShouldBind(&arg); err != nil {
		c.JSON(200, b)
		return
	}
	TeacherID, _ := strconv.Atoi(arg.TeacherID)

	flag, TmpTeacher := database.GetUserInfoById(arg.TeacherID)

	if flag != "Success" {
		b.Code = types.UserNotExisted
		c.JSON(http.StatusOK, b)
		return
	}

	if TmpTeacher.Type != 3 {
		b.Code = types.UserNotExisted
		c.JSON(http.StatusOK, b)
		return
	}

	if TmpTeacher.IsValid != 1 {
		b.Code = types.UserHasDeleted
		c.JSON(http.StatusOK, b)
		return
	}

	n, err, arr := database.GetTeacherCourse(TeacherID)
	if err != nil {
		b.Code = types.UnknownError
		c.JSON(200, b)
		return
	}

	tcourse := make([]*types.TCourse, n)
	for i := 0; i < int(n); i++ {
		_, tempCourse := database.GetOneCourseName(arr[i].CourseID)
		tCourse := types.TCourse{
			TeacherID: arg.TeacherID,
			CourseID:  strconv.Itoa(arr[i].CourseID),
			Name:      tempCourse.Name,
		}
		tcourse[i] = &tCourse
	}

	b.Code = types.OK
	b.Data = struct{ CourseList []*types.TCourse }{CourseList: tcourse}

	c.JSON(200, b)

}
