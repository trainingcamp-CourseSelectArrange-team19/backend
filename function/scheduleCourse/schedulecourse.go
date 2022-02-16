package scheduleCourse

import (
	"backend/types"
	"github.com/gin-gonic/gin"
)


/*
type ScheduleCourseRequest struct {
	TeacherCourseRelationShip map[string][]string // key 为 teacherID , val 为老师期望绑定的课程 courseID 数组
}

type ScheduleCourseResponse struct {
	Code ErrNo
	Data map[string]string   // key 为 teacherID , val 为老师最终绑定的课程 courseID
}
 */
func ScheduleCourse(c *gin.Context) {
	b := types.ScheduleCourseResponse{Code: types.ParamInvalid}
	var arg types.ScheduleCourseRequest
	if err := c.ShouldBindJSON(&arg); err != nil {
		c.JSON(200,b)
		return
	}
	ship := arg.TeacherCourseRelationShip
	match := make(map[string]string,len(ship))
	res := make(map[string]string,len(ship))
	count := 0
	var dfs func(string,map[string][]string,map[string]bool,map[string]string,map[string]string) bool

	dfs = func(u string, group map[string][]string, vis map[string]bool, match map[string]string, res map[string]string) bool {
		for _,v := range group[u]{
			if vis[v] == false {
				vis[v] = true
				value,ok := match[u]
				if !ok || dfs(value, group, vis, match, res) {
					match[v] = u
					res[u] = v
					return true
				}
			}
		}
		return false
	}
	vis := make(map[string]bool,len(ship))
	for k,_ := range ship{
		for _k,_ := range vis{
			vis[_k] = false
		}
		if dfs(k,ship,vis,match,res) {
			count++
		}
	}
	b.Code = types.OK
	b.Data = res
	c.JSON(200,b)
}