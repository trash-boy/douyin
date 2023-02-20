package video

import (
	"TinyTolk/config"
	"TinyTolk/model/Video"
	"TinyTolk/model/user"
	"TinyTolk/model/userfollow"
	video1 "TinyTolk/request/video"
	video2 "TinyTolk/response/video"
	"TinyTolk/utils"
	"TinyTolk/utils/Code"
	"TinyTolk/utils/JWT"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func VideoActionHandler(c *gin.Context){
	var videoActionRequest video1.VideoActionRequest
	if err := c.ShouldBind(&videoActionRequest);err != nil{
		log.Println(err.Error())
		data := *utils.FormVideoActionResponse(Code.VideoActionParamsError,Code.GetMsg(Code.VideoActionParamsError))
		c.JSON(http.StatusOK, data)
		return
	}

	//验证token
	//todo token的时间限制失效了
	if JWT.VerifyToken(c, videoActionRequest.Token) == false{
		data := *utils.FormVideoActionResponse(Code.TokenInvalid,Code.GetMsg(Code.TokenInvalid))
		c.JSON(http.StatusOK, data)
		return
	}

	//#确定视频存储的位置
	userId,exist := c.Get("id")
	if exist == false{
		data := *utils.FormVideoActionResponse(Code.UserIdNotExist,Code.GetMsg(Code.UserIdNotExist))
		c.JSON(http.StatusOK, data)
		return
	}
	//path.Join(config.VideoUrlPrefix, utils.UintToString(userId.(uint))
	VideoDir :=   config.VideoUrlPrefix + utils.UintToString(userId.(uint))
	_,err := utils.PathExists(VideoDir)
	if err != nil{
		log.Println("文件路径不存在")
		data := *utils.FormVideoActionResponse(Code.VideoPathError,Code.GetMsg(Code.VideoPathError))
		c.JSON(http.StatusOK, data)
		return
	}
	//创建MP4文件
	//temp := strings.Split(videoActionRequest.FileHeader.Filename, ".")
	//fileName,fileType := temp[0],temp[1]
	//fileName = fileName
	fileName := videoActionRequest.FileHeader.Filename
	_,err = os.Create(VideoDir + "\\" + fileName)
	if err != nil{
		log.Println("文件路径不存在")
		data := *utils.FormVideoActionResponse(Code.VideoPathError,Code.GetMsg(Code.VideoPathError))
		c.JSON(http.StatusOK, data)
		return
	}

	fileData,err := videoActionRequest.FileHeader.Open()
	if err != nil{
		log.Println("err",err)
	}
	//io.Copy(f, file)
	err = utils.WriteFile(VideoDir + "\\" + fileName,fileData)
	if err != nil{
		data := *utils.FormVideoActionResponse(Code.VideoWriteError,Code.GetMsg(Code.VideoWriteError))
		c.JSON(http.StatusOK, data)
		return
	}

	//todo 生成缩略图


	//C:\douyin\20\Pexels Videos 1474411.mp4
	//C:\douyin\20\Pexels Videos 1474411.jpeg
	fileUrl := VideoDir  +"\\" + fileName

	noTypeName := strings.Split(fileName,".")[0]
	coverName := utils.GetSnapshot(fileUrl, VideoDir + "\\" + noTypeName,1)

	err = Video.InsertVideo(userId.(uint), fileUrl, VideoDir + "\\" + coverName, videoActionRequest.Title)
	if err != nil{
		data := *utils.FormVideoActionResponse(Code.VideoWriteDatabaseError,Code.GetMsg(Code.VideoWriteDatabaseError))
		c.JSON(http.StatusOK, data)
		return
	}
	//工作数+1
	go func() {
		user.AddWorkCount(userId.(uint))
	}()

	data := *utils.FormVideoActionResponse(Code.Success,Code.GetMsg(Code.Success))
	c.JSON(http.StatusOK, data)
	return


}

func VideoListHandler(c *gin.Context){
	var videoListRequest video1.VideoListRequest
	var videoList []video2.Video
	if err := c.ShouldBind(&videoListRequest);err != nil{
		log.Println(err.Error())
		data := *utils.FormVideoListResponse(Code.VideoListParamsError,Code.GetMsg(Code.VideoListParamsError), []video2.Video{})
		c.JSON(http.StatusOK, data)
		return
	}

	// 校验token
	if JWT.VerifyToken(c, videoListRequest.Token) == false{
		data := *utils.FormVideoListResponse(Code.TokenInvalid,Code.GetMsg(Code.TokenInvalid), []video2.Video{})
		c.JSON(http.StatusOK, data)
		return
	}

	selfId, _ := c.Get("id")

	//根据user_id 查询出需要的数据
	err := Video.GetVideoListByUserId(utils.StringToUint(videoListRequest.UserId), &videoList)
	if err != nil{
		data := *utils.FormVideoListResponse(Code.TokenInvalid,Code.GetMsg(Code.TokenInvalid),[]video2.Video{})
		c.JSON(http.StatusOK, data)
		return
	}

	//构建user的信息
	for i := 0 ; i< len(videoList); i++{
		//videoList[i].Author =user2.User{}
		//utils.StringToUint(videoListRequest.UserId)
		exist := user.GetUserById(&videoList[i].Author, videoList[i].UserId)
		if exist != true{
			data := *utils.FormVideoListResponse(Code.UserNameNotExist,Code.GetMsg(Code.UserNameNotExist), []video2.Video{})
			c.JSON(http.StatusOK, data)
			return
		}
		exist, err  := userfollow.UserIsFollow(selfId.(uint),utils.StringToUint(videoListRequest.UserId) )
		if err != nil{
			log.Println(err)
			data := *utils.FormVideoListResponse(Code.UserParamsError,Code.GetMsg(Code.UserParamsError), []video2.Video{})
			c.JSON(http.StatusOK, data)
			return
		}
		videoList[i].Author.IsFollow = exist
	}
	data := *utils.FormVideoListResponse(Code.Success,Code.GetMsg(Code.Success), videoList)
	c.JSON(http.StatusOK, data)
	return

}

func VideoFeedHandler(c *gin.Context){
	var videoFeedRequest video1.VideoFeedRequest
	var videoList []video2.Video
	if err := c.ShouldBind(&videoFeedRequest);err != nil{
		data := *utils.FormVideoFeedResponse(Code.VideoFeedParamsError, Code.GetMsg(Code.VideoFeedParamsError), []video2.Video{},0)
		c.JSON(http.StatusOK, data)
		return
	}
	//检查参数
	//如果没有时间，就使用当前时间
	if videoFeedRequest.LastTime == ""{
		log.Println(time.Now().Unix())
		videoFeedRequest.LastTime = utils.Int64ToString(time.Now().Unix())
	}
	//设置一个随机用户的token
	if videoFeedRequest.Token == ""{
		videoFeedRequest.Token = config.Not_LOGIN_TOKEN
	}
	//验证token是否合法
	if !(videoFeedRequest.Token == config.INVALID_TOKEN || JWT.VerifyToken(c,videoFeedRequest.Token)){
		data := *utils.FormVideoFeedResponse(Code.TokenInvalid, Code.GetMsg(Code.TokenInvalid), []video2.Video{},0)
		c.JSON(http.StatusOK, data)
		return
	}
	//todo 利用用户token进行推荐算法
	//当前直接根据时间龊，拉最多拉30个
	err := Video.GetVideoListByTimeStamp(utils.StringToUint(videoFeedRequest.LastTime), &videoList)
	if err != nil{
		data := *utils.FormVideoFeedResponse(Code.VideoFeedGetDataError, Code.GetMsg(Code.VideoFeedGetDataError), []video2.Video{},0)
		c.JSON(http.StatusOK, data)
		return
	}

	for i := 0 ; i< len(videoList); i++{
		//videoList[i].Author =user2.User{}
		//utils.StringToUint(videoListRequest.UserId)
		exist := user.GetUserById(&videoList[i].Author, videoList[i].UserId)
		if exist != true{
			data := *utils.FormVideoFeedResponse(Code.UserNameNotExist,Code.GetMsg(Code.UserNameNotExist), []video2.Video{},0)
			c.JSON(http.StatusOK, data)
			return
		}
		if videoFeedRequest.Token != config.INVALID_TOKEN{
			selfId,_ := c.Get("id")
			exist, err  := userfollow.UserIsFollow(selfId.(uint),videoList[i].UserId)
			if err != nil{
				log.Println(err)
				data := *utils.FormVideoFeedResponse(Code.UserParamsError,Code.GetMsg(Code.UserParamsError), []video2.Video{},0)
				c.JSON(http.StatusOK, data)
				return
			}
			videoList[i].Author.IsFollow = exist
		}
	}
	data := *utils.FormVideoFeedResponse(Code.Success,Code.GetMsg(Code.Success), videoList,videoList[len(videoList) - 1].CreatedAt.Unix())
	c.JSON(http.StatusOK, data)
	return


}
