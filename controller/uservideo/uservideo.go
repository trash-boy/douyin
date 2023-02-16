package uservideo

import (
	"TinyTolk/model/uservideoFavorite"
	"TinyTolk/request/uservideo"
	"TinyTolk/utils"
	"TinyTolk/utils/Code"
	"TinyTolk/utils/JWT"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserVideoFavoriteHandler(c *gin.Context){
	var userVideoRequest uservideo.UserFavoriteRequest
	if err:= c.ShouldBind(&userVideoRequest);err != nil{
		data := *utils.FormUserVideoFavoriteResponse(Code.UserVideoFavoriteParamsError, Code.GetMsg(Code.UserVideoFavoriteParamsError))
		c.JSON(http.StatusOK, data)
		return
	}

	//1.校验token
	//2.查询数据库是否有（user_id,video_id）没有就创建，有就直接改状态
	if JWT.VerifyToken(c, userVideoRequest.Token) != true{
		data := *utils.FormUserVideoFavoriteResponse(Code.TokenInvalid, Code.GetMsg(Code.TokenInvalid))
		c.JSON(http.StatusOK, data)
		return
	}

	id,_ := c.Get("id")
	if uservideoFavorite.CreateAndUpdate(uint(id.(uint)), utils.StringToUint(userVideoRequest.VideoId),utils.StringToInt(userVideoRequest.ActionType)) != nil{
		data := *utils.FormUserVideoFavoriteResponse(Code.UserVideoFavoriteDatabaseError, Code.GetMsg(Code.UserVideoFavoriteDatabaseError))
		c.JSON(http.StatusOK, data)
		return
	}

	data := *utils.FormUserVideoFavoriteResponse(Code.Success, Code.GetMsg(Code.Success))
	c.JSON(http.StatusOK, data)
	return




}
