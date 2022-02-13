package main

import (
	"backend/function/auth"
	"backend/function/member"
	"backend/function/scheduleCourse"
	"backend/function/selectCourse"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.Engine) {
	g := r.Group("/api/v1")

	// 成员管理
	g.GET("/member", member.GetUser)
	g.GET("/member/list", member.GetUsers)
	g.POST("/member/create", member.CreateMember)
	g.POST("/member/update", member.UpdateName)
	g.POST("/member/delete", member.DeleteMember)

	// 登录
	g.POST("/auth/login", auth.Login)
	g.POST("/auth/logout", auth.AuthMiddleWare(), auth.Logout)
	g.GET("/auth/whoami", auth.AuthMiddleWare(), auth.GetInfo)

	// 排课
	g.POST("/course/create", scheduleCourse.Createcourse)
	g.GET("/course/get", scheduleCourse.Getcourse)
	//g.POST("/teacher/bind_course")
	//g.POST("/teacher/unbind_course")
	g.GET("/teacher/get_course", scheduleCourse.GetTeacherCourse)
	//g.POST("/course/schedule")

	// 抢课
	g.POST("/student/book_course", selectCourse.SelectCourse)
	g.GET("/student/course", selectCourse.FindCourse)
	g.POST("/student/ping", selectCourse.Pong)
}
