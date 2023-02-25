package main

import (
	"TinyTolk/config"
	comment2 "TinyTolk/controller/comment"
	"TinyTolk/controller/soical"
	user2 "TinyTolk/controller/user"
	uservideo2 "TinyTolk/controller/uservideo"
	"TinyTolk/controller/video"
	"TinyTolk/model/DataBaseInit"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	//配置MySQL连接参数
	if config.DB == nil {
		log.Panic("用户数据库连接错误")
	}

	r := gin.Default()
	DataBaseInit.Init()

	v1 := r.Group("douyin")
	{
		v1.POST("/user/register/", user2.UserRegisterHandler)
		v1.POST("/user/login/", user2.UserLoginHandler)
		v1.GET("/user/", user2.GetUserHandler)
		v1.POST("/publish/action/", video.VideoActionHandler)
		v1.GET("/publish/list/", video.VideoListHandler)
		v1.GET("/feed/", video.VideoFeedHandler)

		v1.POST("/favorite/action/", uservideo2.UserVideoFavoriteHandler)
		v1.GET("/favorite/list/", uservideo2.UserGetFavoriteListHandler)

		v1.POST("/comment/action/", comment2.CommentActionHandler)
		v1.GET("/comment/list/", comment2.CommentGetListHandler)

		v1.POST("/relation/action/", soical.RelationActionHandler)
		v1.GET("/relation/follow/list/", soical.RelationFollowListHandler)
		v1.GET("/relation/follower/list/", soical.RelationFollowerListHandler)
	}

	r.Run(":8080")
}
