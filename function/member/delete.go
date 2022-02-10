package member

import (
	"backend/database"
	"backend/types"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

/*
type DeleteMemberRequest struct {
	UserID string
}

type DeleteMemberResponse struct {
	Code ErrNo
}
*/

func DeleteMember(c *gin.Context) {
	b := types.DeleteMemberResponse{Code: types.ParamInvalid}
	var arg types.DeleteMemberRequest
	//参数不对 返回ParamInvalid
	if err := c.ShouldBind(&arg); err != nil {
		c.JSON(200, b)
		return
	}

	//获取用户
	ID := arg.UserID
	u := database.GetUserInfoById(ID)

	//ID为0 用户不存在
	if u.Id == 0 {
		b.Code = types.UserNotExisted
		c.JSON(200, b)
		return
	}

	//已经删除
	if u.IsValid == 0 {
		b.Code = types.UserHasDeleted
		c.JSON(200, b)
		return
	}

	//执行删除
	database.DeleteUser(*u)
	b.Code = types.OK
	c.JSON(200, b)
}
