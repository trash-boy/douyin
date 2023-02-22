package main

import (
	"TinyTolk/config"
	user2 "TinyTolk/controller/user"
	"TinyTolk/middleware"
	"TinyTolk/request/comment"
	"TinyTolk/request/user"
	"TinyTolk/request/uservideo"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"log"
)

func main()  {
	//配置MySQL连接参数
	if config.DB == nil{
		log.Panic("用户数据库连接错误")
	}


	r := gin.Default()
	middleware.Validate = validator.New()
	middleware.Validate .RegisterValidation("ValidateUsername",	user.ValidateUsername)
	middleware.Validate .RegisterValidation("ValidatePassword",	user.ValidatePassword)
	middleware.Validate .RegisterValidation("ValidateActionType",uservideo.ValidateActionType)
	middleware.Validate .RegisterValidation("ValidateCommentVideoId",comment.ValidateCommentVideoId)
	middleware.Validate .RegisterValidation("ValidateCommentActionType",comment.ValidateCommentActionType)
	middleware.Validate .RegisterValidation("ValidateCommentText",comment.ValidateCommentText)

	r.POST("/douyin/user/register", user2.UserRegisterHandler)
	r.POST("/douyin/user/login", user2.UserLoginHandler)
	r.GET("/douyin/user/",user2.GetUserHandler)
	//r.POST("/douyin/publish/action",video.VideoActionHandler)
	//r.GET("/douyin/publish/list/",video.VideoListHandler)
	//r.GET("/douyin/feed/",video.VideoFeedHandler)
	//r.POST("/douyin/favorite/action/", uservideo2.UserVideoFavoriteHandler)
	//r.GET("/douyin/favorite/list/", uservideo2.UserGetFavoriteListHandler)
	//r.POST("/douyin/comment/action/", comment2.CommentActionHandler)
	//r.GET("/douyin/comment/list/", comment2.CommentGetListHandler)
	//r.POST("/douyin/relation/action/", soical.RelationActionHandler)
	//r.GET("/douyin/relation/follow/list/", soical.RelationFollowListHandler)
	//r.GET("/douyin/relation/follower/list/", soical.RelationFollowerListHandler)

	r.Run(":8080")
}


