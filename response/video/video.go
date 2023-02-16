package video

import (
	user2 "TinyTolk/response/user"
	"time"
)

type VideoActionResponse struct {
	StatusCode int32 `form:"status_code" json:"status_code"`
	StatusMsg string `form:"status_msg" json:"status_msg"`
}

type VideoListResponse struct {
	StatusCode int32 `json:"status_code" `
	StatusMsg string `json:"status_msg"`
	VideoList []Video `json:"video_list"`


}

type Video struct {
	Id uint `json:"id" gorm:"id"`

	CreatedAt time.Time`json:"-" gorm:"created_at"`
	Author user2.User `json:"author"`
	UserId uint `json:"-" gorm:"user_id"`
	PlayUrl string `json:"play_url" gorm:"play_url"`
	CoverUrl string `json:"cover_url" gorm:"cover_url"`
	FavoriteCount int64 `json:"favorite_count" gorm:"favorite_count"`
	CommentCount int64 `json:"comment_count" gorm:"comment_count"`
	IsFavorite bool `json:"is_favorite" gorm:"is_favorite"`
	Title string `json:"title" gorm:"title"`
}

type VideoFeedResponse struct {
	StatusCode int32 `json:"status_code" `
	StatusMsg string `json:"status_msg"`
	VideoList []Video `json:"video_list"`
	NextTime int64
}
