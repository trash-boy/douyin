package comment

import (
	"TinyTolk/middleware"
	"TinyTolk/model/Video"
	comment3 "TinyTolk/model/comment"
	"TinyTolk/model/user"
	user3 "TinyTolk/model/user"
	"TinyTolk/model/userfollow"
	"TinyTolk/request/comment"
	comment2 "TinyTolk/response/comment"
	user2 "TinyTolk/response/user"
	"TinyTolk/utils"
	"TinyTolk/utils/Code"
	"TinyTolk/utils/JWT"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

func CommentActionHandler(c *gin.Context){
	var request comment.CommentActionRequest
	if err := c.ShouldBind(&request);err != nil{
		data := *utils.FormCommentActionResponse(Code.CommentParamsError, Code.GetMsg(Code.CommentParamsError), comment2.Comment{})
		c.JSON(http.StatusOK, data)
		return
	}
	if err := middleware.Validate.Struct(request); err !=  nil{
		data := *utils.FormCommentActionResponse(Code.CommentParamsError, Code.GetMsg(Code.CommentParamsError), comment2.Comment{})
		c.JSON(http.StatusOK, data)
		return
	}
	//1.检验token
	if valid := JWT.VerifyToken(c, request.Token);valid != true{
		data := *utils.FormCommentActionResponse(Code.TokenInvalid, Code.GetMsg(Code.TokenInvalid), comment2.Comment{})
		c.JSON(http.StatusOK, data)
		return
	}

	if request.ActionType == "1"{
		//插入评论
		id,_ := c.Get("id")
		result, err := comment3.InsertComment(utils.StringToUint(request.VideoId), id.(uint), request.CommentText)
		if err != nil{
			data := *utils.FormCommentActionResponse(Code.CommentDatabaseError, Code.GetMsg(Code.CommentDatabaseError), comment2.Comment{})
			c.JSON(http.StatusOK, data)
			return
		}

		var commentInfo comment2.Comment
		commentInfo.Content = result.Content
		commentInfo.Id = int64(result.ID)
		//commentInfo.CreateDate = time.Time{result.CreatedAt.Format("01-02")}
		commentInfo.CreatedAt =  result.CreatedAt.Format("01-02")
		exist := user.GetUserById(&commentInfo.User, id.(uint))
		if exist != true{
			data := *utils.FormCommentActionResponse(Code.CommentDatabaseError, Code.GetMsg(Code.CommentDatabaseError), comment2.Comment{})
			c.JSON(http.StatusOK, data)
			return
		}
		//视频评论+1
		go func() {
			Video.AddCommentCount(utils.StringToUint(request.VideoId))
		}()
		data := *utils.FormCommentActionResponse(Code.Success, Code.GetMsg(Code.Success),commentInfo)
		c.JSON(http.StatusOK, data)
		return
	}else{
		//1.删除评论 根据comment_id删除评论
		err := comment3.DeleteCommentById(utils.StringToUint(request.CommentId))
		if err != nil{
			data := *utils.FormCommentActionResponse(Code.CommentDatabaseError, Code.GetMsg(Code.CommentDatabaseError), comment2.Comment{})
			c.JSON(http.StatusOK, data)
			return
		}
		go func() {
			Video.SubCommentCount(utils.StringToUint(request.VideoId))
		}()
		data := *utils.FormCommentActionResponse(Code.Success, Code.GetMsg(Code.Success), comment2.Comment{})
		c.JSON(http.StatusOK, data)
		return

	}

}


func CommentGetListHandler(c *gin.Context){
	var request  comment.CommentGetListRequest
	if err := c.ShouldBind(&request); err != nil{
		data := *utils.FormCommentGetListResponse(Code.CommentGetListParamsError, Code.GetMsg(Code.CommentGetListParamsError),&[]comment2.Comment{})
		c.JSON(http.StatusOK, data)
		return
	}

	if valid := JWT.VerifyToken(c, request.Token);valid != true{
		data := *utils.FormCommentGetListResponse(Code.TokenInvalid, Code.GetMsg(Code.TokenInvalid),&[]comment2.Comment{})
		c.JSON(http.StatusOK, data)
		return
	}
	id,_ := c.Get("id")
	//根据视频id 获取评论
	commentInfo,err := comment3.GetCommentByVideoId(utils.StringToUint(request.VideoId))
	if err != nil{
		data := *utils.FormCommentGetListResponse(Code.CommentDatabaseError, Code.GetMsg(Code.CommentDatabaseError),&[]comment2.Comment{})
		c.JSON(http.StatusOK, data)
		return
	}

	var count sync.WaitGroup
	count.Add(len(*commentInfo))
	for i := 0; i < len(*commentInfo) ; i++ {

		go func(i int) {
			//修改时间格式

			//t,_ :=  time.Parse("2006-01-02 15:04:05",(*commentInfo)[i].CreatedAt)
			currTime := (*commentInfo)[i].CreatedAt[5:10]
			(*commentInfo)[i].CreatedAt = currTime

			//根据用户id获取用户信息
			var tempUser user2.User
			exist := user3.GetUserById(&tempUser, (*commentInfo)[i].UserId)
			if exist != true{
				data := *utils.FormCommentGetListResponse(Code.UserIdNotExist,Code.GetMsg(Code.UserIdNotExist), &[]comment2.Comment{} )
				c.JSON(http.StatusOK, data)
				return
			}

			tempUser.IsFollow,_=  userfollow.UserIsFollow(id.(uint),(*commentInfo)[i].UserId)

			//1.搜索详细的信息
			var userInfo user3.UserInfo
			exist = user3.GetUserInfoByUserId(&userInfo,  (*commentInfo)[i].UserId)
			if exist != true{
				data := *utils.FormCommentGetListResponse(Code.UserIdNotExist,Code.GetMsg(Code.UserIdNotExist), &[]comment2.Comment{}  )
				c.JSON(http.StatusOK, data)
				return
			}
			//将userInfo 数据全部映射到tempUser中
			_ = utils.UserInfoToUser(&tempUser, &userInfo)
			(*commentInfo)[i].User = tempUser
			//count.Done()
			defer func() {
				count.Done()
			}()
		}(i)
	}
	count.Wait()
	data := *utils.FormCommentGetListResponse(Code.Success, Code.GetMsg(Code.Success), commentInfo)
	c.JSON(http.StatusOK, data)
	return

}
