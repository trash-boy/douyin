package Code

const (
	UserRegisterError = 10000
	UserNameNULL      = 10001
	UserNameLong      = 10002
	PasswordNULL      = 10003
	PasswordLong      = 10004
	PasswordError     = 10005
	UserLoginError    = 10006
	Success           = 0
	TokenProduceError = 10007
	UserLoginSuccess  = 10008
	UserNameNotExist  = 10009
	UserParamsError   = 100010
	UserIdNotExist    = 100011
	TokenInvalid      = 100012

	VideoActionParamsError  = 200001
	VideoPathError          = 200002
	VideoWriteError         = 200003
	VideoWriteDatabaseError = 200004

	VideoListParamsError  = 200005
	VideoFeedParamsError  = 200006
	VideoFeedGetDataError = 200007
	VideoTXYError         = 200008

	UserVideoFavoriteParamsError   = 300001
	UserVideoFavoriteDatabaseError = 300002

	UserGetFavoriteListParamsError    = 400001
	UserGetFavoriteListDatatbaseError = 400002

	CommentParamsError        = 500001
	CommentDatabaseError      = 500002
	CommentGetListParamsError = 500003

	RelationActionParamsError         = 600001
	RelationActionDatabaseError       = 600002
	RelationFollowListParamsError     = 600003
	RelationFollowListDatabaseError   = 600004
	RelationFollowerListParamsError   = 600005
	RelationFollowerListDatabaseError = 600006

	RelationDatabaseError = 600007
)

var Msg = map[int]string{
	UserRegisterError: "用户注册错误",
	UserLoginError:    "用户登录错误",
	UserNameNULL:      "用户名为空",
	UserNameLong:      "用户名不能超过32",
	PasswordNULL:      "密码为空",
	PasswordLong:      "密码不能超过32",
	PasswordError:     "密码错误",
	Success:           "成功",
	TokenProduceError: "token生成错误",
	UserLoginSuccess:  "用户登录成功",
	UserNameNotExist:  "用户不存在",
	UserParamsError:   "获取用户信息参数错误",
	UserIdNotExist:    "用户ID不存在",
	TokenInvalid:      "无效token",

	VideoActionParamsError:  "用户上传视频参数错误",
	VideoPathError:          "视频路径创建错误",
	VideoWriteError:         "视频写入本地错误",
	VideoWriteDatabaseError: "视频写入数据库错误",

	VideoListParamsError:           "用户查询上传视频参数错误",
	VideoFeedParamsError:           "视频流参数错误",
	VideoFeedGetDataError:          "视频流数据库获取数据错误",
	VideoTXYError:                  "腾讯云服务错误",
	UserVideoFavoriteParamsError:   "用户视频点赞参数错误",
	UserVideoFavoriteDatabaseError: "用户视频点赞数据库错误",

	UserGetFavoriteListParamsError:    "用户获取点赞视频列表参数错误",
	UserGetFavoriteListDatatbaseError: "用户获取点赞视频列表数据库错误",
	CommentParamsError:                "评论参数错误",
	CommentDatabaseError:              "评论数据库错误",
	CommentGetListParamsError:         "获取评论列表参数错误",

	RelationActionParamsError:         "用户关注参数错误",
	RelationActionDatabaseError:       "用户关注数据库错误",
	RelationFollowListParamsError:     "用户关注列表参数错误",
	RelationFollowListDatabaseError:   "用户关注列表数据库错误",
	RelationFollowerListParamsError:   "用户粉丝列表参数错误",
	RelationFollowerListDatabaseError: "用户粉丝列表数据库错误",
	RelationDatabaseError:             "用户关注数据库错误",
}

func GetMsg(code int) string {
	if val, ok := Msg[code]; ok {
		return val
	}
	return ""
}
