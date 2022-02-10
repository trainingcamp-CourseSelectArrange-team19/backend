package member

import (
	"backend/database"
	"backend/types"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

/*
type UpdateMemberRequest struct {
	UserID   string
	Nickname string
}

type UpdateMemberResponse struct {
	Code ErrNo
}
*/

func UpdateName(c *gin.Context) {
	b := types.UpdateMemberResponse{Code: types.ParamInvalid}

	var arg types.UpdateMemberRequest
	if err := c.ShouldBind(&arg); err != nil {
		c.JSON(200, b)
		return
	}
	Nickname, ID := arg.Nickname, arg.UserID
	u := database.GetUserInfoById(ID)
	if u.Id == 0 {
		b.Code = types.UserNotExisted
		c.JSON(200, b)
		return
	}
	if !(len(Nickname) >= 4 && len(Nickname) <= 20) {
		c.JSON(200, b)
		return
	}
	database.UpdateUserNickname(*u, Nickname)
	b.Code = types.OK
	c.JSON(200, b)
}
