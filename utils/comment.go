package utils

import "TinyTolk/response/comment"

func FormCommentActionResponse(statusCode int32, statusMsg string, c comment.Comment)*comment.CommentActionResponse{
	var response comment.CommentActionResponse
	response.StatusCode = statusCode
	response.StatusMsg = statusMsg
	response.Content = c
	return &response

}

func FormCommentGetListResponse(statusCode int32, statusMsg string, c *[]comment.Comment)*comment.CommentGetListResponse{
	var response comment.CommentGetListResponse
	response.StatusCode = statusCode
	response.StatusMsg = statusMsg
	response.ContentList = *c
	return &response
}
