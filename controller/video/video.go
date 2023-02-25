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
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tencentyun/cos-go-sdk-v5"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

func VideoActionHandler(c *gin.Context) {
	var request video1.VideoActionRequest
	if err := c.ShouldBind(&request); err != nil {
		data := *utils.FormVideoActionResponse(utils.IntToInt32(Code.VideoActionParamsError), Code.GetMsg(Code.VideoActionParamsError))
		c.JSON(http.StatusOK, data)
		return
	}

	//验证token
	//todo token的时间限制失效了
	if JWT.VerifyToken(c, request.Token) == false {
		data := *utils.FormVideoActionResponse(utils.IntToInt32(Code.TokenInvalid), Code.GetMsg(Code.TokenInvalid))
		c.JSON(http.StatusOK, data)
		return
	}

	//#确定视频存储的位置
	userId, exist := c.Get("id")
	if exist == false {
		data := *utils.FormVideoActionResponse(utils.IntToInt32(Code.TokenInvalid), Code.GetMsg(Code.TokenInvalid))
		c.JSON(http.StatusOK, data)
		return
	}

	//path.Join(config.VideoUrlPrefix, utils.UintToString(userId.(uint))
	VideoDir := config.VideoUrlPrefix + utils.UintToString(userId.(uint))

	//不存在文件夹就创建
	_, err := utils.PathExists(VideoDir)
	if err != nil {

		data := *utils.FormVideoActionResponse(utils.IntToInt32(Code.VideoPathError), Code.GetMsg(Code.VideoPathError))
		c.JSON(http.StatusOK, data)
		return
	}
	//创建MP4文件
	//temp := strings.Split(videoActionRequest.FileHeader.Filename, ".")
	//fileName,fileType := temp[0],temp[1]
	//fileName = fileName
	fileName := request.FileHeader.Filename
	_, err = os.Create(VideoDir + "\\" + fileName)
	if err != nil {
		log.Println("文件路径不存在")
		data := *utils.FormVideoActionResponse(utils.IntToInt32(Code.VideoPathError), Code.GetMsg(Code.VideoPathError))
		c.JSON(http.StatusOK, data)
		return
	}

	fileData, err := request.FileHeader.Open()
	if err != nil {
		log.Println("err", err)
	}
	//io.Copy(f, file)
	err = utils.WriteFile(VideoDir+"\\"+fileName, fileData)
	if err != nil {
		data := *utils.FormVideoActionResponse(utils.IntToInt32(Code.VideoWriteError), Code.GetMsg(Code.VideoWriteError))
		c.JSON(http.StatusOK, data)
		return
	}

	//C:\douyin\20\Pexels Videos 1474411.mp4
	//C:\douyin\20\Pexels Videos 1474411.jpeg
	fileUrl := VideoDir + "\\" + fileName

	//将视频文件储存在腾讯云服务
	key, err := TengXunYunFuWu_UpLoad(fileUrl, string(userId.(uint)), "mp4")
	VideoUrl := TenXunYunFuWu_DownLoad(key)
	noTypeName := strings.Split(fileName, ".")[0]
	coverName := utils.GetSnapshot(fileUrl, VideoDir+"\\"+noTypeName, 1)
	keyP, _ := TengXunYunFuWu_UpLoad(VideoDir+"\\"+coverName, string(userId.(uint)), "jpeg")
	photoGraphUrl := TenXunYunFuWu_DownLoad(keyP)
	if err != nil {
		data := *utils.FormVideoActionResponse(Code.VideoTXYError, Code.GetMsg(Code.VideoTXYError))
		c.JSON(http.StatusOK, data)
		return
	}

	err = Video.InsertVideo(userId.(uint), VideoUrl, photoGraphUrl, request.Title)
	if err != nil {
		data := *utils.FormVideoActionResponse(Code.VideoWriteDatabaseError, Code.GetMsg(Code.VideoWriteDatabaseError))
		c.JSON(http.StatusOK, data)
		return
	}

	//工作数+1
	go func() {
		user.AddWorkCount(userId.(uint))
	}()

	data := *utils.FormVideoActionResponse(utils.IntToInt32(Code.Success), Code.GetMsg(Code.Success))
	c.JSON(http.StatusOK, data)
	return

}

func VideoListHandler(c *gin.Context) {
	var request video1.VideoListRequest
	var videoList []video2.Video
	if err := c.ShouldBind(&request); err != nil {
		data := *utils.FormVideoListResponse(utils.IntToInt32(Code.VideoListParamsError), Code.GetMsg(Code.VideoListParamsError), &[]video2.Video{})
		c.JSON(http.StatusOK, data)
		return
	}

	// 校验token
	if JWT.VerifyToken(c, request.Token) == false {
		data := *utils.FormVideoListResponse(utils.IntToInt32(Code.TokenInvalid), Code.GetMsg(Code.TokenInvalid), &[]video2.Video{})
		c.JSON(http.StatusOK, data)
		return
	}

	selfId, _ := c.Get("id")
	//根据user_id 查询出需要的数据
	err := Video.GetVideoListByUserId(selfId.(uint), &videoList)
	if err != nil {
		data := *utils.FormVideoListResponse(utils.IntToInt32(Code.TokenInvalid), Code.GetMsg(Code.TokenInvalid), &[]video2.Video{})
		c.JSON(http.StatusOK, data)
		return
	}

	//构建user的信息
	var count sync.WaitGroup

	count.Add(len(videoList))

	for i := 0; i < len(videoList); i++ {
		//videoList[i].Author =user2.User{}
		//utils.StringToUint(videoListRequest.UserId)
		go func(index int) {
			defer func() {
				count.Done()
			}()
			err := user.GetUserById(&videoList[index].Author, videoList[index].UserId)
			if err != nil {
				//data := *utils.FormVideoListResponse(utils.IntToInt32(Code.UserNameNotExist),Code.GetMsg(Code.UserNameNotExist), &[]video2.Video{})
				//c.JSON(http.StatusOK, data)
				return
			}
			exist, err := userfollow.UserIsFollow(selfId.(uint), utils.StringToUint(request.UserId))
			if err != nil {
				//log.Println(err)
				//data := *utils.FormVideoListResponse(utils.IntToInt32(Code.RelationDatabaseError),Code.GetMsg(Code.RelationDatabaseError), &[]video2.Video{})
				//c.JSON(http.StatusOK, data)
				return
			}
			videoList[index].Author.IsFollow = exist
		}(i)

	}
	count.Wait()
	data := *utils.FormVideoListResponse(Code.Success, Code.GetMsg(Code.Success), &videoList)
	c.JSON(http.StatusOK, data)
	return

}

func VideoFeedHandler(c *gin.Context) {
	var request video1.VideoFeedRequest
	var videoList []video2.Video
	if err := c.ShouldBind(&request); err != nil {
		log.Print("异常0")
		data := *utils.FormVideoFeedResponse(utils.IntToInt32(Code.VideoFeedParamsError), Code.GetMsg(Code.VideoFeedParamsError), &[]video2.Video{}, 0)
		c.JSON(http.StatusOK, data)
		return
	}
	//检查参数
	//如果没有时间，就使用当前时间
	if request.LastTime == "" {
		request.LastTime = utils.Int64ToString(time.Now().Unix())
	}
	//设置一个随机用户的token
	if request.Token == "" {
		request.Token = config.Not_LOGIN_TOKEN
	}
	//验证token是否合法
	if !(request.Token == config.INVALID_TOKEN || JWT.VerifyToken(c, request.Token)) {
		log.Print("异常1")
		data := *utils.FormVideoFeedResponse(utils.IntToInt32(Code.TokenInvalid), Code.GetMsg(Code.TokenInvalid), &[]video2.Video{}, 0)
		c.JSON(http.StatusOK, data)
		return
	}
	//todo 利用用户token进行推荐算法
	//当前直接根据时间龊，拉最多拉30个
	err := Video.GetVideoListByTimeStamp(utils.UintToInt64(utils.StringToUint(request.LastTime)), &videoList)
	if err != nil {
		log.Print("异常2")
		data := *utils.FormVideoFeedResponse(utils.IntToInt32(Code.VideoFeedGetDataError), Code.GetMsg(Code.VideoFeedGetDataError), &[]video2.Video{}, 0)
		c.JSON(http.StatusOK, data)
		return
	}

	var count sync.WaitGroup
	count.Add(len(videoList))

	for i := 0; i < len(videoList); i++ {
		//videoList[i].Author =user2.User{}
		//utils.StringToUint(videoListRequest.UserId)
		go func(index int) {
			defer func() {
				count.Done()
			}()
			err := user.GetUserById(&videoList[index].Author, videoList[index].UserId)
			if err != nil {
				return
			}
			if request.Token != config.INVALID_TOKEN {
				selfId, _ := c.Get("id")
				exist, err := userfollow.UserIsFollow(selfId.(uint), videoList[index].UserId)
				if err != nil {
					return
				}
				videoList[index].Author.IsFollow = exist
			}
		}(i)

	}
	count.Wait()
	data := *utils.FormVideoFeedResponse(Code.Success, Code.GetMsg(Code.Success), &videoList, videoList[len(videoList)-1].CreatedAt.Unix())
	log.Print(videoList)
	c.JSON(http.StatusOK, data)
	log.Print("正常")
	return

}

func TengXunYunFuWu_UpLoad(filePath string, userId string, typename string) (string, error) {
	// 将 examplebucket-1250000000 和 COS_REGION 修改为真实的信息
	// 存储桶名称，由 bucketname-appid 组成，appid 必须填入，可以在 COS 控制台查看存储桶名称。https://console.cloud.tencent.com/cos5/bucket
	// COS_REGION 可以在控制台查看，https://console.cloud.tencent.com/cos5/bucket, 关于地域的详情见 https://cloud.tencent.com/document/product/436/6224
	u, _ := url.Parse("https://1375958258-1314408528.cos.ap-nanjing.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv(config.SecretId),  // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考 https://cloud.tencent.com/document/product/598/37140
			SecretKey: os.Getenv(config.SecretKey), // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考 https://cloud.tencent.com/document/product/598/37140
		},
	})
	// 对象键（Key）是对象在存储桶中的唯一标识。
	// 例如，在对象的访问域名 `examplebucket-1250000000.cos.COS_REGION.myqcloud.com/test/objectPut.go` 中，对象键为 test/objectPut.go
	name := "douyin/" + userId + time.Now().String() + "." + typename
	// 2.通过本地文件上传对象
	_, err := c.Object.PutFromFile(context.Background(), name, filePath, nil)
	if err != nil {
		panic(err)
	}
	return name, nil
}

func TenXunYunFuWu_DownLoad(key string) string {
	// 存储桶名称，由 bucketname-appid 组成，appid 必须填入，可以在 COS 控制台查看存储桶名称。 https://console.cloud.tencent.com/cos5/bucket
	// 替换为用户的 region，存储桶 region 可以在 COS 控制台“存储桶概览”查看 https://console.cloud.tencent.com/ ，关于地域的详情见 https://cloud.tencent.com/document/product/436/6224 。
	u, _ := url.Parse("https://1375958258-1314408528.cos.ap-nanjing.myqcloud.com")
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 通过环境变量获取密钥
			// 环境变量 SECRETID 表示用户的 SecretId，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretID: os.Getenv(config.SecretId), // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参见 https://cloud.tencent.com/document/product/598/37140
			// 环境变量 SECRETKEY 表示用户的 SecretKey，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
			SecretKey: os.Getenv(config.SecretKey), // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参见 https://cloud.tencent.com/document/product/598/37140
		},
	})
	ourl := client.Object.GetObjectURL(key)
	fmt.Println(ourl)
	return ourl.String()
}
