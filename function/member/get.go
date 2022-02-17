package member

import (
	"backend/database"
	"backend/types"
	"strconv"

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
	_, u := database.GetUserInfoById(ID)
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

func GetUsers(c *gin.Context) {
	b := types.GetMemberListResponse{Code: types.ParamInvalid, Data: struct{ MemberList []types.TMember }{MemberList: make([]types.TMember, 0)}}
	var arg types.GetMemberListRequest
	if err := c.ShouldBind(&arg); err != nil {
		c.JSON(200, b)
		return
	}
	offset, limit := arg.Offset, arg.Limit
	sz, _, users := database.GetAllValUserInfo()
	if sz <= int64(offset) {
		c.JSON(200, b)
		return
	}
	var res []types.TMember
	b.Code = types.OK
	for i := int64(offset); i < Min(sz, int64(offset+limit)); i++ {
		res = append(res, types.TMember{
			UserID:   strconv.Itoa(users[i].Id),
			Nickname: users[i].Nickname,
			Username: users[i].Name,
			UserType: users[i].Type,
		})
	}

	b.Data = struct{ MemberList []types.TMember }{MemberList: res}
	c.JSON(200, b)

}

func Min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}
