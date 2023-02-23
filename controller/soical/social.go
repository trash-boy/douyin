package soical

import (
	"TinyTolk/model/user"
	"TinyTolk/model/userfollow"
	"TinyTolk/request/social"
	user2 "TinyTolk/response/user"
	"TinyTolk/utils"
	"TinyTolk/utils/Code"
	"TinyTolk/utils/JWT"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sync"
)

func RelationActionHandler(c *gin.Context){
	//1.校验参数
	//2.验证token
	//2.在关注表中修改数据（如果没有就先创建)
	//3.对应用户关注数和被关注数对应增减
	var request social.RelationActionRequest
	if err := c.ShouldBind(&request);err != nil{
		data := * utils.FormRelationActionResponse(Code.RelationActionParamsError,Code.GetMsg(Code.RelationActionParamsError))
		c.JSON(http.StatusOK, data)
		return
	}

	if valid := JWT.VerifyToken(c, request.Token);valid != true {
		data := * utils.FormRelationActionResponse(Code.TokenInvalid,Code.GetMsg(Code.TokenInvalid))
		c.JSON(http.StatusOK, data)
		return
	}
	followerId ,_ := c.Get("id")
	err := userfollow.InsertNotExist(utils.StringToUint(request.ToUserId),followerId.(uint) )
	if err != nil{
		data := * utils.FormRelationActionResponse(Code.RelationActionDatabaseError,Code.GetMsg(Code.RelationActionDatabaseError))
		c.JSON(http.StatusOK, data)
		return
	}

	if  request.ActionType == "1"{
		err := userfollow.UpdateFollow(utils.StringToUint(request.ToUserId),followerId.(uint), true)
		if err != nil{
			data := * utils.FormRelationActionResponse(Code.RelationActionDatabaseError,Code.GetMsg(Code.RelationActionDatabaseError))
			c.JSON(http.StatusOK, data)
			return
		}
		//用户详情表关注数和被关注+1
		go func() {
			user.AddFollowCount(followerId.(uint))
			user.AddFollowerCount(utils.StringToUint(request.ToUserId))
		}()
	}else{
		err := userfollow.UpdateFollow(utils.StringToUint(request.ToUserId),followerId.(uint), false)
		if err != nil{
			data := * utils.FormRelationActionResponse(Code.RelationActionDatabaseError,Code.GetMsg(Code.RelationActionDatabaseError))
			c.JSON(http.StatusOK, data)
			return
		}
		//用户详情表关注数和被关注减1
		go func() {
			user.SubFollowCount(followerId.(uint))
			user.SubFollowerCount(utils.StringToUint(request.ToUserId))
		}()
	}
	data := * utils.FormRelationActionResponse(Code.Success,Code.GetMsg(Code.Success))
	c.JSON(http.StatusOK, data)
	return

}

func RelationFollowListHandler(c *gin.Context){
	var request social.RelationFollowListRequest
	if err := c.ShouldBind(&request); err != nil{
		data := *utils.FormRelationFollowListResponse(Code.RelationFollowListParamsError, Code.GetMsg(Code.RelationFollowListParamsError), &[]user2.User{})
		c.JSON(http.StatusOK, data)
		return
	}

	if valid := JWT.VerifyToken(c, request.Token); valid != true{
		data := *utils.FormRelationFollowListResponse(Code.TokenInvalid, Code.GetMsg(Code.TokenInvalid), &[]user2.User{})
		c.JSON(http.StatusOK, data)
		return
	}

	//检验token与当前用户是否一致
	id,_ := c.Get("id")
	if id.(uint) != utils.StringToUint(request.UserId){
		data := *utils.FormRelationFollowListResponse(Code.TokenInvalid, Code.GetMsg(Code.TokenInvalid), &[]user2.User{})
		c.JSON(http.StatusOK, data)
		return
	}

	followIdList, err := userfollow.GetFollowIdListByUserId(utils.StringToUint(request.UserId))
	if err != nil{
		data := *utils.FormRelationFollowListResponse(Code.RelationFollowListDatabaseError, Code.GetMsg(Code.RelationFollowListDatabaseError), &[]user2.User{})
		c.JSON(http.StatusOK, data)
		return
	}

	var userInfo = make([]user2.User, len(followIdList))
	var count sync.WaitGroup
	count.Add(len(followIdList))

	for idx,userId := range followIdList{
		go func(index int, userId uint) {
			_ = user.GetUserById(&userInfo[index],userId)

			userInfo[index].IsFollow,_=  userfollow.UserIsFollow(id.(uint),userId)

			//1.搜索详细的信息
			var userDetail user.UserInfo
			_ = user.GetUserInfoByUserId(&userDetail,  userId)

			//将userInfo 数据全部映射到tempUser中
			_ = utils.UserInfoToUser(&userInfo[index], &userDetail)
			defer func() {
				count.Done()
			}()
		}(idx, userId)
	}
	count.Wait()

	data := *utils.FormRelationFollowListResponse(Code.Success, Code.GetMsg(Code.Success), &userInfo)
	c.JSON(http.StatusOK, data)
	return
}


func RelationFollowerListHandler(c *gin.Context){
	var request social.RelationFollowerListRequest
	if err := c.ShouldBind(&request); err != nil{
		data := *utils.FormRelationFollowerListResponse(Code.RelationFollowerListParamsError, Code.GetMsg(Code.RelationFollowerListParamsError), &[]user2.User{})
		c.JSON(http.StatusOK, data)
		return
	}

	if valid := JWT.VerifyToken(c, request.Token); valid != true{
		data := *utils.FormRelationFollowListResponse(Code.TokenInvalid, Code.GetMsg(Code.TokenInvalid), &[]user2.User{})
		c.JSON(http.StatusOK, data)
		return
	}

	//检验token与当前用户是否一致
	id,_ := c.Get("id")
	if id.(uint) != utils.StringToUint(request.UserId){
		data := *utils.FormRelationFollowerListResponse(Code.TokenInvalid, Code.GetMsg(Code.TokenInvalid), &[]user2.User{})
		c.JSON(http.StatusOK, data)
		return
	}

	followerIdList, err := userfollow.GetFollowerIdListByUserId(utils.StringToUint(request.UserId))
	if err != nil{
		data := *utils.FormRelationFollowerListResponse(Code.RelationFollowerListDatabaseError, Code.GetMsg(Code.RelationFollowerListDatabaseError), &[]user2.User{})
		c.JSON(http.StatusOK, data)
		return
	}

	var userInfo = make([]user2.User, len(followerIdList))
	var count sync.WaitGroup
	count.Add(len(followerIdList))

	for idx,userId := range followerIdList{

		go func(index int, userId uint) {
			exist := user.GetUserById(&userInfo[index],userId)
			log.Println(exist)

			userInfo[index].IsFollow,err=  userfollow.UserIsFollow(userId,id.(uint))
			log.Println(err)

			//1.搜索详细的信息
			var userDetail user.UserInfo
			exist = user.GetUserInfoByUserId(&userDetail,  userId)
			log.Println(exist)

			//将userInfo 数据全部映射到tempUser中
			_ = utils.UserInfoToUser(&userInfo[index], &userDetail)
			defer func() {
				count.Done()
			}()
		}(idx, userId)
	}
	count.Wait()

	data := *utils.FormRelationFollowerListResponse(Code.Success, Code.GetMsg(Code.Success), &userInfo)
	c.JSON(http.StatusOK, data)
	return
}
