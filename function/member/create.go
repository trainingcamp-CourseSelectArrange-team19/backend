package member

import (
	"backend/database"
	"backend/function/selectCourse"
	"backend/tools"
	"backend/types"
	"strconv"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

/*
type CreateMemberRequest struct {
	Nickname string   // required，不小于 4 位 不超过 20 位
	Username string   // required，只支持大小写，长度不小于 8 位 不超过 20 位
	Password string   // required，同时包括大小写、数字，长度不少于 8 位 不超过 20 位
	UserType UserType // required, 枚举值
}
type CreateMemberResponse struct {
	Code ErrNo
	Data struct {
		UserID string // int64 范围
	}
}
*/

func CreateMember(c *gin.Context) {
	/* 参数不正确 提前设置返回值 */
	/* 发送来的不符合要求 */
	b := types.CreateMemberResponse{
		Code: types.ParamInvalid,
		Data: struct {
			UserID string
		}{
			UserID: "****",
		},
	}
	var arg types.CreateMemberRequest
	if err := c.ShouldBind(&arg); err != nil {
		c.JSON(200, b)
		return
	}
	Nickname, Username, Password, UserType := arg.Nickname, arg.Username, arg.Password, arg.UserType
	/* 整体长度校验 + 用户类型校验 */
	v1 := !(len(Nickname) >= 4 && len(Nickname) <= 20)
	v2 := !(len(Username) >= 8 && len(Username) <= 20)
	v3 := !(len(Password) >= 8 && len(Password) <= 20)
	v4 := !(UserType == types.Admin || UserType == types.Student || UserType == types.Teacher)

	/* 账户名校验 */
	for i := 0; i < len(Username); i++ {
		v2 = v2 || !((Username[i] >= 'a' && Username[i] <= 'z') || (Username[i] >= 'A' && Username[i] <= 'Z'))
	}
	/* 密码校验 */
	o1, o2, o3 := false, false, false
	for i := 0; i < len(Password); i++ {
		o1 = o1 || (Password[i] >= 'a' && Password[i] <= 'z')
		o2 = o2 || (Password[i] >= 'A' && Password[i] <= 'Z')
		o3 = o3 || (Password[i] >= '0' && Password[i] <= '9')
	}

	v3 = v3 || !(o1 && o2 && o3)

	if v1 || v2 || v3 || v4 {
		c.JSON(200, b)
		return
	}

	//鉴权
	cookie, err1 := c.Cookie("camp-session")
	if err1 != nil {
		panic(err1)
	}
	//fmt.Println(cookie)
	_, organizer := database.GetUserInfoByName(cookie)
	if organizer.Type != 1 {
		b.Code = types.PermDenied
		c.JSON(200, b)
		return
	}

	_, u := database.GetUserInfoByName(Username)
	if u.Id != 0 {
		b.Code = types.UserHasExisted
		c.JSON(200, b)
		return
	}

	database.CreateUser(Username, Nickname, Password, UserType)
	redisConn := selectCourse.RedisPool.Get()
	defer redisConn.Close()
	_, tempUser := database.GetUserInfoByName(Username)
	_, err := redis.Int64(redisConn.Do("BF.ADD", "studentsID", strconv.Itoa(tempUser.Id)))
	if err != nil {
		tools.LogMsg(err)
		panic(err)
	}

	b.Code = types.OK
	b.Data = struct{ UserID string }{UserID: strconv.Itoa(tempUser.Id)}
	c.JSON(200, b)

}
