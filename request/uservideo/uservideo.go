package uservideo

import (

	"gopkg.in/go-playground/validator.v8"
	"reflect"
)


type UserFavoriteRequest struct {
	Token string `form:"token"  json:"token"`
	VideoId string `form:"video_id"  json:"video_id"`
	ActionType string `form:"action_type"  json:"action_type"  validate:"required,gte=1,lte=2"`
}

func ValidateActionType(v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool {

	if value, ok := field.Interface().(string); ok {
		return value == "1" || value == "2"
	}

	return false
}

