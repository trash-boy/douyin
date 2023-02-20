package utils

import (
	"TinyTolk/response/uservideo"
	"TinyTolk/response/video"
)

func FormUserVideoFavoriteResponse(statusCode int32, statusMsg string)*uservideo.UserFavoriteResponse{
	var uservideo uservideo.UserFavoriteResponse
	uservideo.StatusCode = statusCode
	uservideo.StatusMsg = statusMsg
	return &uservideo
}

func FormUserGetFavoriteResponse(statusCode int32, statusMsg string,  videoList []video.Video)*uservideo.UserGetFavoriteListResponse{
	var response uservideo.UserGetFavoriteListResponse
	response.StatusCode = statusCode
	response.StatusMsg = statusMsg
	response.VideoList = videoList
	return &response
}

