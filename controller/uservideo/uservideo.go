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
	"log"
	"net/http"
	"sync"
)

func UserVideoFavoriteHandler(c *gin.Context) {
	var request uservideo.UserFavoriteRequest
	if err := c.ShouldBind(&request); err != nil {
		data := *utils.FormUserVideoFavoriteResponse(utils.IntToInt32(Code.UserVideoFavoriteParamsError), Code.GetMsg(Code.UserVideoFavoriteParamsError))
		c.JSON(http.StatusOK, data)
		return
	}

	if valid := middleware.Validate.Struct(request); valid != nil {
		data := *utils.FormUserVideoFavoriteResponse(utils.IntToInt32(Code.UserVideoFavoriteParamsError), Code.GetMsg(Code.UserVideoFavoriteParamsError))
		c.JSON(http.StatusOK, data)
		return
	}
	//1.校验token
	//2.查询数据库是否有（user_id,video_id）没有就创建，有就直接改状态
	if JWT.VerifyToken(c, request.Token) != true {
		data := *utils.FormUserVideoFavoriteResponse(utils.IntToInt32(Code.TokenInvalid), Code.GetMsg(Code.TokenInvalid))
		c.JSON(http.StatusOK, data)
		return
	}

	id, _ := c.Get("id")

	if uservideoFavorite.CreateAndUpdate(id.(uint), utils.Int64ToUInt(utils.StringToInt64(request.VideoId)), utils.StringToInt(request.ActionType)) != nil {
		data := *utils.FormUserVideoFavoriteResponse(utils.IntToInt32(Code.UserVideoFavoriteDatabaseError), Code.GetMsg(Code.UserVideoFavoriteDatabaseError))
		c.JSON(http.StatusOK, data)
		return
	}
	//操作数据表
	go func() {
		if request.ActionType == "1" {
			//点赞的数量加+1，获赞人数量+1
			//利用redis优化
			_ = user.AddFavoriteCount(id.(uint))
			_ = user.AddGetFavoriteCount(Video.GetUserIdByVideoId(utils.Int64ToUInt(utils.StringToInt64(request.VideoId))))
			_ = Video.AddFavoriteCount(utils.Int64ToUInt(utils.StringToInt64(request.VideoId)))

		} else {
			//利用redis优化
			_ = user.SubFavoriteCount(id.(uint))
			_ = user.SubGetFavoriteCount(Video.GetUserIdByVideoId(utils.Int64ToUInt(utils.StringToInt64(request.VideoId))))
			_ = Video.SubFavoriteCount(utils.Int64ToUInt(utils.StringToInt64(request.VideoId)))

		}
	}()
	data := *utils.FormUserVideoFavoriteResponse(utils.IntToInt32(Code.Success), Code.GetMsg(Code.Success))
	c.JSON(http.StatusOK, data)
	return

}

func UserGetFavoriteListHandler(c *gin.Context) {
	var request uservideo.UserGetFavoriteListRequest
	if err := c.ShouldBind(&request); err != nil {
		data := *utils.FormUserGetFavoriteResponse(utils.IntToInt32(Code.UserGetFavoriteListParamsError), Code.GetMsg(Code.UserGetFavoriteListParamsError), &[]video.Video{})
		c.JSON(http.StatusOK, data)
		return
	}
	//1.校验token
	if JWT.VerifyToken(c, request.Token) != true {
		data := *utils.FormUserGetFavoriteResponse(utils.IntToInt32(Code.TokenInvalid), Code.GetMsg(Code.TokenInvalid), &[]video.Video{})
		c.JSON(http.StatusOK, data)
		return
	}

	userId, _ := c.Get("id")
	//if userId.(uint) != utils.Int64ToUInt(utils.StringToInt64(request.UserId)) {
	//	data := *utils.FormUserGetFavoriteResponse(utils.IntToInt32(Code.TokenInvalid), Code.GetMsg(Code.TokenInvalid), &[]video.Video{})
	//	c.JSON(http.StatusOK, data)
	//	return
	//}
	//2.根据user_id去评论表中得到所有的视频id
	log.Print(userId)
	result, err := uservideoFavorite.GetVideoIdByUserId(userId.(uint))
	if err != nil {
		data := *utils.FormUserGetFavoriteResponse(utils.IntToInt32(Code.UserGetFavoriteListDatatbaseError), Code.GetMsg(Code.UserGetFavoriteListDatatbaseError), &[]video.Video{})
		c.JSON(http.StatusOK, data)
		return
	}

	//根据视频id获取视频数据
	//videoList := []video.Video{}
	videoList := make([]video.Video, len(result))

	var count sync.WaitGroup
	count.Add(len(result))

	for i := 0; i < len(result); i++ {

		go func(index int) {
			defer func() {
				count.Done()
			}()
			err := Video.GetVideoByVideoId(result[index].VideoId, &videoList[index])
			if err != nil {
				return
			}
			err = user.GetUserById(&videoList[index].Author, videoList[index].UserId)
			if err != nil {
				return
			}
		}(i)

	}

	count.Wait()
	data := *utils.FormUserGetFavoriteResponse(Code.Success, Code.GetMsg(Code.Success), &videoList)
	c.JSON(http.StatusOK, data)
	return

}
