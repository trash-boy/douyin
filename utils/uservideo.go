package utils

import "TinyTolk/response/uservideo"

func FormUserVideoFavoriteResponse(statusCode int32, statusMsg string)*uservideo.UserFavoriteResponse{
	var uservideo uservideo.UserFavoriteResponse
	uservideo.StatusCode = statusCode
	uservideo.StatusMsg = statusMsg
	return &uservideo
}
