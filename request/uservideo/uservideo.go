package uservideo

import (
	"github.com/go-playground/validator/v10"
)

type UserFavoriteRequest struct {
	Token      string `form:"token" `
	VideoId    string `form:"video_id"  `
	ActionType string `form:"action_type"   validate:"required,ValidateActionType"`
}

type UserGetFavoriteListRequest struct {
	UserId string `form:"user_id" binding:"required"`
	Token  string `form:"token" binding:"required"`
	userId string
}

func ValidateActionType(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	return value == "1" || value == "2"

}
