package video

import "mime/multipart"

type VideoActionRequest struct {
	Token string `form:"token"`
	FileHeader  multipart.FileHeader  `form:"data"`
	Title string `form:"title"`
}


type VideoListRequest struct {
	UserId string `form:"user_id" json:"user_id"`
	Token string  `form:"token" json:"token"`
}

type VideoFeedRequest struct {
	LastTime string `form:"last_time" json:"last_time"`
	Token string `form:"token" json:"token"`
}