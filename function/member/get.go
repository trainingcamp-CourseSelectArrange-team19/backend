package member

import (
	"backend/database"
	"backend/types"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

/*
如果用户已删除请返回已删除状态码，不存在请返回不存在状态码
type GetMemberRequest struct {
	UserID string
}
type GetMemberResponse struct {
	Code ErrNo
	Data TMember
}
type TMember struct {
	UserID   string
	Nickname string
	Username string
	UserType UserType
}
*/

func GetUser(c *gin.Context) {
	b := types.GetMemberResponse{Code: types.ParamInvalid}

	var arg types.GetMemberRequest
	if err := c.ShouldBind(&arg); err != nil {
		c.JSON(200, b)
		return
	}
	ID := arg.UserID
	u := database.GetUserInfoById(ID)
	//不存在
	if u.Id == 0 {
		b.Code = types.UserNotExisted
		c.JSON(200, b)
		return
	}
	//已删除
	if u.IsValid == 0 {
		b.Code = types.UserHasDeleted
		c.JSON(200, b)
		return
	}

	//正常执行
	b.Data = types.TMember{
		UserID:   ID,
		Nickname: u.Nickname,
		Username: u.Name,
		UserType: u.Type,
	}
	b.Code = types.OK
	c.JSON(200, b)

}

/*
批量获取成员信息
type GetMemberListRequest struct {
	Offset int
	Limit  int
}
type GetMemberListResponse struct {
	Code ErrNo
	Data struct {
		MemberList []TMember
	}
}
*/

/* func GetUsers(c *gin.Context) {
	b := types.GetMemberListResponse{Code: types.ParamInvalid}

	var arg types.GetMemberListRequest
	if err := c.ShouldBind(&arg); err != nil {
		c.JSON(200, b)
		return
	}

	c.JSON(200, b)

}
*/
