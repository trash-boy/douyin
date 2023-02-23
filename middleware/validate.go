package middleware

import (
	"TinyTolk/request/comment"
	"TinyTolk/request/user"
	"TinyTolk/request/uservideo"
	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init(){
	Validate = validator.New()
	Validate.RegisterValidation("ValidateUsername",	user.ValidateUsername)
	Validate.RegisterValidation("ValidatePassword",	user.ValidatePassword)
	Validate.RegisterValidation("ValidateActionType",uservideo.ValidateActionType)
	Validate.RegisterValidation("ValidateCommentVideoId",comment.ValidateCommentVideoId)
	Validate.RegisterValidation("ValidateCommentActionType",comment.ValidateCommentActionType)
	Validate.RegisterValidation("ValidateCommentText",comment.ValidateCommentText)
}