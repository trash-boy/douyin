package comment

import (
	"TinyTolk/middleware"
	"TinyTolk/model/Video"
	comment3 "TinyTolk/model/comment"
	"TinyTolk/model/user"
	"TinyTolk/model/userfollow"
	"TinyTolk/request/comment"
	comment2 "TinyTolk/response/comment"
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
		data := *utils.FormCommentActionResponse(utils.IntToInt32(Code.CommentParamsError), Code.GetMsg(Code.CommentParamsError), &comment2.Comment{})
		c.JSON(http.StatusOK, data)
		return
	}

	if err := middleware.Validate.Struct(request); err !=  nil{
		data := *utils.FormCommentActionResponse(utils.IntToInt32(Code.CommentParamsError), Code.GetMsg(Code.CommentParamsError), &comment2.Comment{})
		c.JSON(http.StatusOK, data)
		return
	}
	//1.检验token
	if valid := JWT.VerifyToken(c, request.Token);valid != true{
		data := *utils.FormCommentActionResponse(utils.IntToInt32(Code.TokenInvalid), Code.GetMsg(Code.TokenInvalid), &comment2.Comment{})
		c.JSON(http.StatusOK, data)
		return
	}

	if request.ActionType == "1"{
		//插入评论
		id,_ := c.Get("id")

		result, err := comment3.InsertComment(utils.Int64ToUInt(utils.StringToInt64(request.VideoId)), id.(uint), request.CommentText)
		if err != nil{
			data := *utils.FormCommentActionResponse(utils.IntToInt32(Code.CommentDatabaseError), Code.GetMsg(Code.CommentDatabaseError), &comment2.Comment{})
			c.JSON(http.StatusOK, data)
			return
		}

		var commentInfo comment2.Comment

		commentInfo.Content = result.Content
		commentInfo.Id = int64(result.ID)
		//commentInfo.CreateDate = time.Time{result.CreatedAt.Format("01-02")}
		commentInfo.CreatedAt =  result.CreatedAt.Format("01-02")

		err = user.GetUserById(&commentInfo.User, id.(uint))
		if err != nil{
			data := *utils.FormCommentActionResponse(utils.IntToInt32(Code.CommentDatabaseError), Code.GetMsg(Code.CommentDatabaseError), &comment2.Comment{})
			c.JSON(http.StatusOK, data)
			return
		}
		//视频评论+1
		go func() {
			_ = Video.AddCommentCount(utils.Int64ToUInt(utils.StringToInt64(request.VideoId)))
		}()
		data := *utils.FormCommentActionResponse(utils.IntToInt32(Code.Success), Code.GetMsg(Code.Success),&commentInfo)
		c.JSON(http.StatusOK, data)
		return

	}else{
		//1.删除评论 根据comment_id删除评论
		//todo 评论表设计其实不用status，硬删除即可
		err := comment3.DeleteCommentById(utils.StringToUint(request.CommentId))
		if err != nil{
			data := *utils.FormCommentActionResponse(utils.IntToInt32(Code.CommentDatabaseError), Code.GetMsg(Code.CommentDatabaseError), &comment2.Comment{})
			c.JSON(http.StatusOK, data)
			return
		}
		go func() {
			_ = Video.SubCommentCount(utils.Int64ToUInt(utils.StringToInt64(request.VideoId)))
		}()
		data := *utils.FormCommentActionResponse(utils.IntToInt32(Code.Success), Code.GetMsg(Code.Success), &comment2.Comment{})
		c.JSON(http.StatusOK, data)
		return

	}

}


func CommentGetListHandler(c *gin.Context){
	var request  comment.CommentGetListRequest
	if err := c.ShouldBind(&request); err != nil{
		data := *utils.FormCommentGetListResponse(utils.IntToInt32(Code.CommentGetListParamsError), Code.GetMsg(Code.CommentGetListParamsError),&[]comment2.Comment{})
		c.JSON(http.StatusOK, data)
		return
	}

	if valid := JWT.VerifyToken(c, request.Token);valid != true{
		data := *utils.FormCommentGetListResponse(utils.IntToInt32(Code.TokenInvalid), Code.GetMsg(Code.TokenInvalid),&[]comment2.Comment{})
		c.JSON(http.StatusOK, data)
		return
	}

	id,_ := c.Get("id")
	//根据视频id 获取评论
	commentInfo,err := comment3.GetCommentByVideoId(utils.Int64ToUInt(utils.StringToInt64(request.VideoId)))
	if err != nil{
		data := *utils.FormCommentGetListResponse(utils.IntToInt32(Code.CommentDatabaseError), Code.GetMsg(Code.CommentDatabaseError),&[]comment2.Comment{})
		c.JSON(http.StatusOK, data)
		return
	}

	var count sync.WaitGroup
	count.Add(len(*commentInfo))
	for i := 0; i < len(*commentInfo) ; i++ {

		go func(index int) {
			//修改时间格式
			defer func() {
				count.Done()
			}()
			//t,_ :=  time.Parse("2006-01-02 15:04:05",(*commentInfo)[i].CreatedAt)
			currTime := (*commentInfo)[index].CreatedAt[5:10]
			(*commentInfo)[index].CreatedAt = currTime

			//根据用户id获取用户信息

			err := user.GetUserById(&(*commentInfo)[index].User, (*commentInfo)[index].UserId)
			if err != nil{
				data := *utils.FormCommentGetListResponse(utils.IntToInt32(Code.UserIdNotExist),Code.GetMsg(Code.UserIdNotExist), &[]comment2.Comment{} )
				c.JSON(http.StatusOK, data)
				return
			}
			(*commentInfo)[index].User.IsFollow ,_=  userfollow.UserIsFollow(id.(uint),(*commentInfo)[index].UserId)

		}(i)
	}
	count.Wait()
	data := *utils.FormCommentGetListResponse(utils.IntToInt32(Code.Success), Code.GetMsg(Code.Success), commentInfo)
	c.JSON(http.StatusOK, data)
	return

}
