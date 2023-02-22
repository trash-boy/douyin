package utils

import (
	"TinyTolk/response/social"
	"TinyTolk/response/user"
)

func FormRelationActionResponse(statusCode int32, statusMsg string)*social.RelationActionResponse  {
	var response social.RelationActionResponse
	response.StatusCode = statusCode
	response.StatusMsg = statusMsg
	return &response
}

func FormRelationFollowListResponse(statusCode int32, statusMsg string, userList *[]user.User)*social.RelationFollowListResponse  {
	var response social.RelationFollowListResponse
	response.StatusCode = statusCode
	response.StatusMsg = statusMsg
	response.UserList = *userList
	return &response
}


func FormRelationFollowerListResponse(statusCode int32, statusMsg string, userList *[]user.User)*social.RelationFollowerListResponse  {
	var response social.RelationFollowerListResponse
	response.StatusCode = statusCode
	response.StatusMsg = statusMsg
	response.UserList = *userList
	return &response
}
