package user

import (
	//"github.com/go-playground/validator/v10"
	"gopkg.in/go-playground/validator.v8"
	"reflect"
)

type UserRegisterRequest struct {
	Username string  `form:"username" json:"username" validate:"ValidateUsername"`
	Password string ` form:"password" json:"password" validate:"ValidatePassword"`
}

type UserLoginRequest struct {
	Username string `form:"username" json:"username" validate:"required, ValidateUsername"`
	Password string `form:"password" json:"password"  validate:"required, ValidateUsername"`
}

type UserReuqest struct {
	UserId string `form:"user_id"`
	Token string `form:"token"`
}



func ValidateUsername(v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool {

	if value, ok := field.Interface().(string); ok {
		// 字段不能为空，并且不等于  admin
		return value != "" && len(value) <= 32
	}

	return false
}

func ValidatePassword(v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool {

	if value, ok := field.Interface().(string); ok {
		// 字段不能为空，并且不等于  admin
		return value != "" && len(value) <= 32
	}

	return false
}
