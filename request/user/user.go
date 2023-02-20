package user

import (
	"github.com/go-playground/validator/v10"
)

type UserRegisterRequest struct {
	Username string  `form:"username"  validate:"required,ValidateUsername"`
	Password string ` form:"password" validate:"required,ValidatePassword"`
}

type UserLoginRequest struct {
	Username string `form:"username"  validate:"required,ValidateUsername"`
	Password string `form:"password"  validate:"required,ValidatePassword"`
}

type UserReuqest struct {
	UserId string `form:"user_id"`
	Token string `form:"token"`
}


//todo validate 失效

func ValidateUsername(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return value != "" && len(value) <= 32

}


func ValidatePassword(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return value != "" && len(value) <= 32

}

