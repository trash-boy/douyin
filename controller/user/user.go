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
func UserRegisterHandler(c *gin.Context){
	var userRegister user.UserRegisterRequest
	if err := c.Bind(&userRegister); err != nil {
		data := *utils.FormUserRegisterResponse(Code.UserRegisterError, Code.GetMsg(Code.UserRegisterError), 0, config.INVALID_TOKEN)
		c.JSON(http.StatusOK, data)
		return
	}
	if valid := middleware.Validate.Struct(userRegister);valid != nil{
		data := *utils.FormUserRegisterResponse(Code.UserRegisterError, Code.GetMsg(Code.UserRegisterError), 0, config.INVALID_TOKEN)
		c.JSON(http.StatusOK, data)
		return
	}
	//处理数据验证成功逻辑
	//生成数据存入数据库
	var user user3.User
	user.Username = userRegister.Username
	user.Password = encryption.Encoding( userRegister.Password)
	user.Name = utils.UserPrefix + userRegister.Username
	result := user3.InsertUser(&user)
	if result.Error != nil{
		data := *utils.FormUserRegisterResponse(Code.UserRegisterError, Code.GetMsg(Code.UserRegisterError),config.INVALID_USER_ID,"")
		c.JSON(http.StatusOK,data)
		return
	}
	token,err := JWT.SetToken(user.ID)
	if err != nil{
		data := *utils.FormUserRegisterResponse(Code.TokenProduceError, Code.GetMsg(Code.TokenProduceError),config.INVALID_USER_ID,"")
		c.JSON(http.StatusOK,data)
		return
	}
	//填写详细信息表
	go func() {
		var userInfo user3.UserInfo
		userInfo.Avatar = config.AVATAR_URL
		userInfo.BackgroudImage = config.BACKGROUD_IMAGE
		userInfo.UserId = user.ID
		err := user3.InsertUserInfo(&userInfo)
		if err != nil{
			log.Println(err)
		}

	}()


	data := *utils.FormUserRegisterResponse(Code.Success, Code.GetMsg(Code.Success),user.ID,token)
	c.JSON(http.StatusOK,data)
	return

}


func UserLoginHandler(c *gin.Context){
	var userLogin user.UserLoginRequest
	if err := c.Bind(&userLogin); err != nil{
		data := *utils.FormUserLoginResponse(Code.UserLoginError, Code.GetMsg(Code.UserLoginError),config.INVALID_USER_ID,"")
		c.JSON(http.StatusOK,data)
		return
	}

	//处理数据验证成功逻辑
	//从数据库里根据username获取密码进行对比
	var user user3.User
	exist := user3.GetUserByUsername(&user, userLogin.Username)
	if exist != true{
		data := *utils.FormUserLoginResponse(Code.UserNameNotExist, Code.GetMsg(Code.UserNameNotExist),user.ID,config.INVALID_TOKEN)
		c.JSON(http.StatusOK,data)
		return
	}
	if encryption.Encoding(userLogin.Password) != user.Password{
		data := *utils.FormUserLoginResponse(Code.PasswordError, Code.GetMsg(Code.PasswordError),user.ID,config.INVALID_TOKEN)
		c.JSON(http.StatusOK,data)
		return
	}

	token,err := JWT.SetToken(user.ID)
	if err != nil{
		data := *utils.FormUserRegisterResponse(Code.TokenProduceError, Code.GetMsg(Code.TokenProduceError),config.INVALID_USER_ID,config.INVALID_TOKEN)
		c.JSON(http.StatusOK,data)
		return
	}
	data := *utils.FormUserRegisterResponse(Code.Success, Code.GetMsg(Code.Success),user.ID,token)
	c.JSON(http.StatusOK,data)
	return
}

func GetUserHandler(c *gin.Context){

	var userRequest  user.UserReuqest
	var tempUser user2.User

	//绑定数据
	if err := c.ShouldBindQuery(&userRequest); err != nil {
		log.Println("err,",err.Error())
		data := *utils.FormUserResponse(Code.UserParamsError,Code.GetMsg(Code.UserParamsError), tempUser )
		c.JSON(http.StatusOK, data)
		return
	}

	//验证tokenid是否有效
	tokenIsValid := JWT.VerifyToken(c, userRequest.Token)
	if tokenIsValid == false{
		data := *utils.FormUserResponse(Code.TokenInvalid,Code.GetMsg(Code.TokenInvalid), tempUser )
		c.JSON(http.StatusOK, data)
		return
	}
	id,exist := c.Get("id")
	log.Println("id:",id)
	if exist != true{
		id = "0"
	}

	//根据id查询响应的数据
	//首先根据id去用户信息表中查询对应数据，其次去关注表中查询是否有当前关注的字段

	exist = user3.GetUserById(&tempUser, utils.StringToUint(userRequest.UserId))
	if exist != true{
		data := *utils.FormUserResponse(Code.UserIdNotExist,Code.GetMsg(Code.UserIdNotExist), tempUser )
		c.JSON(http.StatusOK, data)
		return
	}

	tempUser.IsFollow,_=  userfollow.UserIsFollow(id.(uint),utils.StringToUint(userRequest.UserId))

	//1.搜索详细的信息
	var userInfo user3.UserInfo
	exist = user3.GetUserInfoByUserId(&userInfo,  utils.StringToUint(userRequest.UserId))
	if exist != true{
		data := *utils.FormUserResponse(Code.UserIdNotExist,Code.GetMsg(Code.UserIdNotExist), tempUser )
		c.JSON(http.StatusOK, data)
		return
	}
	//将userInfo 数据全部映射到tempUser中
	_ = utils.UserInfoToUser(&tempUser, &userInfo)
	data := *utils.FormUserResponse(Code.Success,Code.GetMsg(Code.Success), tempUser )
	c.JSON(http.StatusOK, data)
	return


}
