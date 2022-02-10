package selectCourse

import (
	"backend/database"
	"backend/tools"
	"backend/types"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"strconv"
)

func FindCourse(c *gin.Context) {
	b := types.GetStudentCourseResponse{
		Code: types.ParamInvalid,

	}
	var arg types.GetStudentCourseRequest
	if err := c.ShouldBind(&arg); err != nil {
		c.JSON(500, b)
		panic(err.Error())
		return
	}
	StudentID, err1 := strconv.Atoi(arg.StudentID)
	if err1 != nil {
		tools.LogMsg(err1)
		panic(err1)
	}
	InsertSchedule(arg.StudentID)
	rows, err2, schedules := database.GetStudentCourse(StudentID)
	if err2 != nil {
		tools.LogMsg(err2)
		panic(err2)
	}
	if rows > 0 {
		b.Code = types.StudentHasCourse
		b.Data.CourseList = make([]types.TCourse, len(schedules))
		for ind := 0 ; ind < len(schedules) ; ind++ {
			b.Data.CourseList[ind].TeacherID = strconv.Itoa(schedules[ind].TeacherID)
			b.Data.CourseList[ind].Name = schedules[ind].CourseName
			b.Data.CourseList[ind].CourseID = strconv.Itoa(schedules[ind].CourseID)
		}
		c.JSON(200, b)
	} else {
		b.Code = types.StudentHasNoCourse
		c.JSON(200, b)
	}
}

func InsertSchedule(StudentID string){
	redisConn := redisPool.Get()
	defer redisConn.Close()
	val, err := redis.Ints(redisConn.Do("SMEMBERS", "courses:" + StudentID + ":uids"))
	if err != nil {
		tools.LogMsg(err)
		panic(err)
	}
	_, err = redisConn.Do("DEL", "courses:" + StudentID + ":uids")
	if err != nil {
		tools.LogMsg(err)
		panic(err)
	}
	SID, err1 := strconv.Atoi(StudentID)
	if err1 != nil {
		tools.LogMsg(err1)
		panic(err1)
	}
	for i, _ := range val{
		count, err, _ := database.GetStudentCourseAbsent(SID, val[i])
		if err != nil {
			tools.LogMsg(err)
			panic(err)
		}
		if count > 0 {
			continue
		}
		_, courseFound := database.GetOneCourseName(val[i])
		_, teacherFound := database.GetCourseTeacher(val[i])
		database.CreateStudentCourse(SID, val[i], courseFound.Name, teacherFound.TeacherID)
	}
}