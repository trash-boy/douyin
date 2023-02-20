package comment

import (
	"TinyTolk/response/user"
)

type CommentActionResponse struct {
	StatusCode int32 `json:"status_code"`
	StatusMsg string `json:"status_msg"`
	Content Comment `json:"content"`
}

type Comment struct {
	Id int64 `json:"id" gorm:"id"`
	User user.User `json:"user"`
	UserId uint `json:"-" gorm:"user_id"`
	Content string  `json:"content" gorm:"content"`
	CreatedAt string`json:"create_date" gorm:"created_at"`
}

type CommentGetListResponse struct {
	StatusCode int32 `json:"status_code"`
	StatusMsg string `json:"status_msg"`
	ContentList []Comment `json:"content_list"`
}