package utils

import user2 "TinyTolk/response/user"

const (
	UserPrefix = "douyin_"
)
func FormUserRegisterResponse(statusCode int32, statusMsg string, userId uint,token string)*user2.UserRegisterResponse{
	var response user2.UserRegisterResponse
	response.StatusCode = statusCode
	response.StatusMsg = statusMsg
	response.UserId = userId
	response.Token = token
	return &response
}


func FormUserLoginResponse(statusCode int32, statusMsg string, userId uint,token string)*user2.UserLoginResponse{
	var response user2.UserLoginResponse
	response.StatusCode = statusCode
	response.StatusMsg = statusMsg
	response.UserId = userId
	response.Token = token
	return &response
}

func FormUserResponse(statusCode int32, statusMsg string, user user2.User)*user2.UserResponse{
	var response user2.UserResponse
	response.StatusCode = statusCode
	response.StatusMsg = statusMsg
	response.User = user
	return &response
}



