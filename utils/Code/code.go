package Code

const (
	UserRegisterError = 10000
	UserNameNULL =  10001
	UserNameLong = 10002
	PasswordNULL = 10003
	PasswordLong = 10004
	PasswordError = 10005
	UserLoginError = 10006
	Success = 0
	TokenProduceError = 10007
	UserLoginSuccess = 10008
	UserNameNotExist = 10009
	UserParamsError = 100010
	UserIdNotExist = 100011
	TokenInvalid=100012

	VideoActionParamsError = 200001
	VideoPathError = 200002
	VideoWriteError = 200003
	VideoWriteDatabaseError = 200004

	VideoListParamsError = 200005
	VideoFeedParamsError = 200006
	VideoFeedGetDataError = 200007

	UserVideoFavoriteParamsError = 300001
	UserVideoFavoriteDatabaseError = 300002
)

var Msg = map[int]string{
	UserRegisterError:"用户注册错误",
	UserLoginError:"用户登录错误",
	UserNameNULL : "用户名为空",
	UserNameLong: "用户名不能超过32",
	PasswordNULL : "密码为空",
	PasswordLong: "密码不能超过32",
	PasswordError: "密码错误",
	Success:"成功",
	TokenProduceError:"token生成错误",
	UserLoginSuccess:"用户登录成功",
	UserNameNotExist:"用户不存在",
	UserParamsError:"获取用户信息参数错误",
	UserIdNotExist:"用户ID不存在",
	TokenInvalid:"无效token",

	VideoActionParamsError:"用户上传视频参数错误",
	VideoPathError:"视频路径创建错误",
	VideoWriteError:"视频写入本地错误",
	VideoWriteDatabaseError:"视频写入数据库错误",

	VideoListParamsError:"用户查询上传视频参数错误",
	VideoFeedParamsError:"视频流参数错误",
	VideoFeedGetDataError:"视频流数据库获取数据错误",
	UserVideoFavoriteParamsError:"用户视频点赞参数错误",
	UserVideoFavoriteDatabaseError:"用户视频点赞数据库错误",


}

func GetMsg(code int)string{
	if val,ok := Msg[code] ;ok{
		return val
	}
	return ""
}


