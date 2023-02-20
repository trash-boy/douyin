package comment

import "github.com/go-playground/validator/v10"

type CommentActionRequest struct {
	Token string `form:"token"`
	VideoId string `form:"video_id" gorm:"video_id" validate:"required,ValidateCommentVideoId"`
	ActionType string `form:"action_type" validate:"required,ValidateCommentActionType"`
	CommentText string `form:"comment_text" gorm:"content" validate:"required,ValidateCommentText"`
	CommentId string `form:"comment_id"`
}

type CommentGetListRequest struct {
	Token string `form:"token"`
	VideoId string `form:"video_id"`
}

func ValidateCommentVideoId(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return len(value) > 0
}

func ValidateCommentActionType(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return value == "1" || value == "2"
}

func ValidateCommentText(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return len(value) > 0
}
