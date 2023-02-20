package uservideo

import (
	"TinyTolk/middleware"
	"TinyTolk/model/Video"
	"TinyTolk/model/user"
	"TinyTolk/model/uservideoFavorite"
	"TinyTolk/request/uservideo"
	"TinyTolk/response/video"
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

	if valid := middleware.Validate.Struct(userVideoRequest);valid != nil{
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
	//操作数据表
	go func() {
		if userVideoRequest.ActionType == "1"{
			//点赞的数量加+1，获赞人数量+1
			user.AddFavoriteCount(id.(uint))
			Video.AddFavoriteCount(utils.StringToUint(userVideoRequest.VideoId))
			user.AddGetFavoriteCount(Video.GetUserIdByVideoId(utils.StringToUint(userVideoRequest.VideoId)))
		}else{
			user.SubFavoriteCount(id.(uint))
			Video.SubFavoriteCount(utils.StringToUint(userVideoRequest.VideoId))
			user.SubGetFavoriteCount(Video.GetUserIdByVideoId(utils.StringToUint(userVideoRequest.VideoId)))
		}
	}()
	data := *utils.FormUserVideoFavoriteResponse(Code.Success, Code.GetMsg(Code.Success))
	c.JSON(http.StatusOK, data)
	return




}

func UserGetFavoriteListHandler(c *gin.Context){
	var userGetFavorite uservideo.UserGetFavoriteListRequest
	if err := c.ShouldBind(&userGetFavorite);err != nil{
		data := *utils.FormUserGetFavoriteResponse(Code.UserGetFavoriteListParamsError,Code.GetMsg(Code.UserGetFavoriteListParamsError), []video.Video{})
		c.JSON(http.StatusOK, data)
		return
	}
	//1.校验token
	if JWT.VerifyToken(c, userGetFavorite.Token) != true{
		data := *utils.FormUserGetFavoriteResponse(Code.TokenInvalid,Code.GetMsg(Code.TokenInvalid), []video.Video{})
		c.JSON(http.StatusOK, data)
		return
	}
	//2.根据user_id去评论表中得到所有的视频id
	result,err := uservideoFavorite.GetVideoIdByUserId(utils.StringToUint(userGetFavorite.UserId))
	if err != nil{
		data := *utils.FormUserGetFavoriteResponse(Code.UserGetFavoriteListDatatbaseError,Code.GetMsg(Code.UserGetFavoriteListDatatbaseError), []video.Video{})
		c.JSON(http.StatusOK, data)
		return
	}

	//根据视频id获取视频数据
	videoList := []video.Video{}
	for i := 0 ; i < len(result); i++{
		var videoInfo video.Video
		err := Video.GetVideoByVideoId(result[i].VideoId,&videoInfo)
		if err != nil{
			data := *utils.FormUserGetFavoriteResponse(Code.UserGetFavoriteListDatatbaseError,Code.GetMsg(Code.UserGetFavoriteListDatatbaseError), []video.Video{})
			c.JSON(http.StatusOK, data)
			return
		}
		exist :=  user.GetUserById(&videoInfo.Author, videoInfo.UserId)
		if exist != true{
			data := *utils.FormUserGetFavoriteResponse(Code.UserGetFavoriteListDatatbaseError,Code.GetMsg(Code.UserGetFavoriteListDatatbaseError), []video.Video{})
			c.JSON(http.StatusOK, data)
			return
		}
		videoList = append(videoList, videoInfo)
	}

	data := *utils.FormUserGetFavoriteResponse(Code.Success,Code.GetMsg(Code.Success), videoList)
	c.JSON(http.StatusOK, data)
	return

}
