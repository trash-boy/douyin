package uservideo

import "TinyTolk/response/video"
type UserFavoriteResponse struct {
	StatusCode int32 `form:"status_code" json:"status_code"`
	StatusMsg string `form:"status_msg" json:"status_msg"`
}

type UserGetFavoriteListResponse struct{
	StatusCode int32 `json:"status_code" `
	StatusMsg string `json:"status_msg"`
	VideoList []video.Video `json:"video_list"`
}
