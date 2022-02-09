package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	Id       int
	Name     string
	Nickname string
	Password string
	Type     int
	IsValid  int
}

func (User) TableName() string {
	return "user"
}

type Course struct {
	Id       int
	Name     string
	Capacity int
}

func (Course) TableName() string {
	return "course"
}

type TeacherSchedule struct {
	ID        int
	TeacherID int
	CourseID  int
}

func (TeacherSchedule) TableName() string {
	return "teacher_schedule"
}

type StudentSchedule struct {
	ID        int
	StudentID int
	CourseID  int
	CourseName string
	TeacherID int
}

func (StudentSchedule) TableName() string {
	return "student_schedule"
}

var db *gorm.DB

//连接数据库
func Connect() {
	var err error
	//dsn := "root:bytedancecamp@tcp(180.184.65.192:3306)/test1?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := "root:ru19870528@tcp(127.0.0.1:3306)/test1?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(fmt.Sprintf("open mysql failed, err is %s", err))
	}
}

//创建成员
func CreateUser(Username string, Nickname string, Password string, Usertype int) string {
	newUser := &User{
		Name:     Username,
		Nickname: Nickname,
		Password: Password,
		Type:     Usertype,
	}
	if err := db.Create(newUser).Error; err != nil {
		return fmt.Sprintf("create failed, err is %s", err)
		//panic(fmt.Sprintf("create failed, err is %s", err))
	}
	return "create successes!"
}

//根据用户名获取单个用户信息
func GetUserInfoByName(Username string) (string, *User) {
	TempUser := new(User)
	if err := db.First(TempUser, "name = ?", Username).Error; err != nil {
		return fmt.Sprintf("query failed ,err is %s", err), TempUser
	}
	return "Success", TempUser
}
func GetUserInfoById(ID string) *User {
	TempUser := new(User)
	if err := db.First(TempUser, "id = ?", ID).Error; err != nil {
		return TempUser
	}
	return TempUser	
}

//获取所有用户信息，包括已删除的用户，
func GetAllUserInfo() (int64, []User) {
	var users []User
	rows := db.Find(&users).RowsAffected
	return rows, users
}

//获取所有有效的用户信息，不包括已删除的用户，
func GetAllValUserInfo() (int64, error, []User) {
	var users []User
	result := db.Where("is_valid = ?", "1").Find(&users)
	return result.RowsAffected, result.Error, users
}

//获取所有有效的学生，不包括已删除的用户，
func GetAllValStudentInfo() (int64, error, []User) {
	var users []User
	result := db.Where("is_valid = ? and type = ?", "1", "2").Find(&users)
	return result.RowsAffected, result.Error, users
}

//更新用户昵称
func UpdateUserNickname(ID string, Nickname string) error {
	user := GetUserInfoById(ID)
	result := db.Model(&user).Where("id = ?", user.Id).Update("nickname", Nickname)
	return result.Error
}

//软删除用户
func DeleteUser(user User) error {
	result := db.Model(&user).Update("is_valid", 0)
	return result.Error
}

//创建课程
func CreateCourse(name string, capacity int) string {
	newCourse := &Course{
		Name:     name,
		Capacity: capacity,
	}
	if err := db.Create(newCourse).Error; err != nil {
		return fmt.Sprintf("create failed, err is %s", err)
		//panic(fmt.Sprintf("create failed, err is %s", err))
	}
	return "create successes!"
}

//获取单个课程
func GetOneCourse(name string) (string, *Course) {
	TempCourse := new(Course)
	if err := db.First(TempCourse, "name = ?", name).Error; err != nil {
		return fmt.Sprintf("query failed ,err is %s", err), TempCourse
	}
	return "", TempCourse
}

//获取单个课程名称
func GetOneCourseName(ID int) (string, *Course) {
	TempCourse := new(Course)
	if err := db.First(TempCourse, "id = ?", ID).Error; err != nil {
		return fmt.Sprintf("query failed ,err is %s", err), TempCourse
	}
	return "", TempCourse
}

//获取所有用户信息，包括已删除的用户，
func GetAllCourse() (int64, []Course) {
	var courses []Course
	rows := db.Find(&courses).RowsAffected
	return rows, courses
}

//绑定老师和课程
func BindTeacherCourse(teacherID int, courseID int) string {
	newTeacherSchedule := &TeacherSchedule{
		TeacherID: teacherID,
		CourseID:  courseID,
	}
	if err := db.Create(newTeacherSchedule).Error; err != nil {
		return fmt.Sprintf("bind failed, err is %s", err)
	}
	return "bind successes!"
}

//解绑老师和课程
func UnbindTeacherCourse(teacherID int, courseID int) string {
	TempTeacherSchedule := new(TeacherSchedule)
	if err := db.First(TempTeacherSchedule, "teacher_id = ? AND course_id = ?", teacherID, courseID).Error; err != nil {
		return fmt.Sprintf("unbind failed ,err is %s", err)
	}
	if err := db.Delete(&TempTeacherSchedule).Error; err != nil {
		return fmt.Sprintf("unbind failed, err is %s", err)
	}
	return "unbind successes!"
}

//先使用GetOneUserInfo得到老师id，再使用该函数
func GetTeacherCourse(teacherID int) (int64, error, []TeacherSchedule) {
	var teacherSchedule []TeacherSchedule
	result := db.Where("teacher_id = ?", teacherID).Find(&teacherSchedule)
	return result.RowsAffected, result.Error, teacherSchedule
}

//先使用GetOneUserInfo得到老师id，再使用该函数
func GetCourseTeacher(courseID int) (string, *TeacherSchedule) {
	TempTeacherSchedule := new(TeacherSchedule)
	if err := db.First(TempTeacherSchedule, "course_id = ?", courseID).Error; err != nil {
		return fmt.Sprintf("query failed ,err is %s", err), TempTeacherSchedule
	}
	return "", TempTeacherSchedule
}

//创建学生课表
func CreateStudentCourse(studentID int, courseID int, courseName string, teacherID int) string {
	newStudentSchedule := &StudentSchedule{
		StudentID: studentID,
		CourseID:  courseID,
		CourseName: courseName,
		TeacherID:  teacherID,
	}
	if err := db.Create(newStudentSchedule).Error; err != nil {
		return fmt.Sprintf("create failed, err is %s", err)
	}
	return "create successes!"
}

//先使用GetOneUserInfo得到学生id，再使用该函数
func GetStudentCourse(studentID int) (int64, error, []StudentSchedule) {
	var studentSchedule []StudentSchedule
	result := db.Where("student_id = ?", studentID).Find(&studentSchedule)
	return result.RowsAffected, result.Error, studentSchedule
}

//查看是否已存在该课
func GetStudentCourseAbsent(studentID int, courseID int) (int64, error, []StudentSchedule) {
	var studentSchedule []StudentSchedule
	result := db.Where("student_id = ? and course_id = ?", studentID, courseID).Find(&studentSchedule)
	return result.RowsAffected, result.Error, studentSchedule
}


//func main() {
//	connect()
//	res := CreateUser("aaa", "张三", "asdfasdf", 2)
//	fmt.Println(res)
//	GetOneUserInfo("张三")
//}
