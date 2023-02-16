package main

import (
	"TinyTolk/config"
	user2 "TinyTolk/controller/user"
	"TinyTolk/controller/video"
	"TinyTolk/request/user"
	"TinyTolk/request/uservideo"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	uservideo2 "TinyTolk/controller/uservideo"
	//"github.com/go-playground/validator/v10"
	"gopkg.in/go-playground/validator.v8"
	"log"
)

func main()  {
	//配置MySQL连接参数
	if config.DB == nil{
		log.Panic("用户数据库连接错误")
	}


	r := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 这里的 key 和 fn 可以不一样最终在 struct 使用的是 key
		v.RegisterValidation("ValidateUsername", user.ValidateUsername)
		v.RegisterValidation("ValidatePassword", user.ValidatePassword)
		v.RegisterValidation("ValidateActionType", uservideo.ValidateActionType)
	}

	r.POST("/douyin/user/register", user2.UserRegisterHandler)
	r.POST("/douyin/user/login", user2.UserLoginHandler)
	r.GET("/douyin/user/",user2.GetUserHandler)
	r.POST("/douyin/publish/action",video.VideoActionHandler)
	r.GET("/douyin/publish/list/",video.VideoListHandler)
	r.GET("/douyin/feed/",video.VideoFeedHandler)
	r.POST("/douyin/favorite/action/", uservideo2.UserVideoFavoriteHandler)
	r.Run(":8080")
}


