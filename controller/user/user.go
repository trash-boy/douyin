package user

import (
	"TinyTolk/config"
	"TinyTolk/middleware"
	user3 "TinyTolk/model/user"
	"TinyTolk/model/userfollow"
	"TinyTolk/request/user"
	user2 "TinyTolk/response/user"
	"TinyTolk/utils"
	"TinyTolk/utils/Code"
	"TinyTolk/utils/JWT"
	"TinyTolk/utils/encryption"
	"log"

	"github.com/gin-gonic/gin"
	"net/http"
)

func UserRegisterHandler(c *gin.Context) {

	var request user.UserRegisterRequest
	if err := c.Bind(&request); err != nil {
		data := *utils.FormUserRegisterResponse(utils.IntToInt32(Code.UserRegisterError), Code.GetMsg(Code.UserRegisterError), utils.IntToInt64(config.INVALID_USER_ID), config.INVALID_TOKEN)
		c.JSON(http.StatusOK, data)
		return
	}
	if valid := middleware.Validate.Struct(request); valid != nil {
		data := *utils.FormUserRegisterResponse(utils.IntToInt32(Code.UserRegisterError), Code.GetMsg(Code.UserRegisterError), utils.IntToInt64(config.INVALID_USER_ID), config.INVALID_TOKEN)
		c.JSON(http.StatusOK, data)
		return
	}

	userId, err := user3.InsertUser(request.Username, request.Password)
	if err != nil {
		data := *utils.FormUserRegisterResponse(utils.IntToInt32(Code.UserRegisterError), Code.GetMsg(Code.UserRegisterError), utils.IntToInt64(config.INVALID_USER_ID), config.INVALID_TOKEN)
		c.JSON(http.StatusOK, data)
		return
	}
	token, err := JWT.SetToken(userId)
	if err != nil {
		data := *utils.FormUserRegisterResponse(utils.IntToInt32(Code.TokenProduceError), Code.GetMsg(Code.TokenProduceError), utils.IntToInt64(config.INVALID_USER_ID), config.INVALID_TOKEN)
		c.JSON(http.StatusOK, data)
		return
	}
	data := *utils.FormUserRegisterResponse(utils.IntToInt32(Code.Success), Code.GetMsg(Code.Success), utils.UintToInt64(userId), token)
	c.JSON(http.StatusOK, data)
	return

}
func UserLoginHandler(c *gin.Context) {
	var request user.UserLoginRequest
	if err := c.Bind(&request); err != nil {
		data := *utils.FormUserLoginResponse(utils.IntToInt32(Code.UserLoginError), Code.GetMsg(Code.UserLoginError), utils.IntToInt64(config.INVALID_USER_ID), config.INVALID_TOKEN)
		c.JSON(http.StatusOK, data)
		return
	}

	//处理数据验证成功逻辑
	//从数据库里根据username获取密码进行对比
	var user user3.User
	exist := user3.GetUserByUsername(&user, request.Username)
	if exist != true {
		data := *utils.FormUserLoginResponse(utils.IntToInt32(Code.UserNameNotExist), Code.GetMsg(Code.UserNameNotExist), utils.UintToInt64(config.INVALID_USER_ID), config.INVALID_TOKEN)
		c.JSON(http.StatusOK, data)
		return
	}
	if encryption.Encoding(request.Password) != user.Password {
		data := *utils.FormUserLoginResponse(utils.IntToInt32(Code.PasswordError), Code.GetMsg(Code.PasswordError), utils.UintToInt64(config.INVALID_USER_ID), config.INVALID_TOKEN)
		c.JSON(http.StatusOK, data)
		return
	}

	token, err := JWT.SetToken(user.ID)
	if err != nil {
		data := *utils.FormUserRegisterResponse(Code.TokenProduceError, Code.GetMsg(Code.TokenProduceError), utils.UintToInt64(config.INVALID_USER_ID), config.INVALID_TOKEN)
		c.JSON(http.StatusOK, data)
		return
	}
	data := *utils.FormUserRegisterResponse(Code.Success, Code.GetMsg(Code.Success), utils.UintToInt64(user.ID), token)
	c.JSON(http.StatusOK, data)
	return
}
func GetUserHandler(c *gin.Context) {

	var request user.UserReuqest
	var tempUser user2.User

	//绑定数据
	if err := c.ShouldBindQuery(&request); err != nil {

		data := *utils.FormUserResponse(utils.IntToInt32(Code.UserParamsError), Code.GetMsg(Code.UserParamsError), &tempUser)
		c.JSON(http.StatusOK, data)
		return
	}

	//验证tokenid是否有效
	tokenIsValid := JWT.VerifyToken(c, request.Token)
	if tokenIsValid == false {
		data := *utils.FormUserResponse(utils.IntToInt32(Code.TokenInvalid), Code.GetMsg(Code.TokenInvalid), &tempUser)
		c.JSON(http.StatusOK, data)
		return
	}

	id, exist := c.Get("id")
	log.Println("id:", id)
	if exist != true {
		id = "0"
	}
	if id.(uint) != utils.Int64ToUInt(utils.StringToInt64(request.UserId)) {
		data := *utils.FormUserResponse(utils.IntToInt32(Code.TokenInvalid), Code.GetMsg(Code.TokenInvalid), &tempUser)
		c.JSON(http.StatusOK, data)
		return
	}

	//根据id查询响应的数据
	//首先根据id去用户信息表中查询对应数据，其次去关注表中查询是否有当前关注的字段

	err := user3.GetUserById(&tempUser, utils.Int64ToUInt(utils.StringToInt64(request.UserId)))
	if err != nil {
		data := *utils.FormUserResponse(utils.IntToInt32(Code.UserIdNotExist), Code.GetMsg(Code.UserIdNotExist), &tempUser)
		c.JSON(http.StatusOK, data)
		return
	}

	tempUser.IsFollow, _ = userfollow.UserIsFollow(id.(uint), utils.Int64ToUInt(utils.StringToInt64(request.UserId)))

	//1.搜索详细的信息
	//var userInfo user3.UserInfo
	//exist = user3.GetUserInfoByUserId(&userInfo,  utils.StringToUint(userRequest.UserId))
	//if exist != true{
	//	data := *utils.FormUserResponse(Code.UserIdNotExist,Code.GetMsg(Code.UserIdNotExist), tempUser )
	//	c.JSON(http.StatusOK, data)
	//	return
	//}
	////将userInfo 数据全部映射到tempUser中
	//_ = utils.UserInfoToUser(&tempUser, &userInfo)
	data := *utils.FormUserResponse(utils.IntToInt32(Code.Success), Code.GetMsg(Code.Success), &tempUser)
	c.JSON(http.StatusOK, data)
	return

}
