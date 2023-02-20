package user

import (
	"TinyTolk/config"
	"TinyTolk/response/user"
	"testing"
)

func Test_GetUserInfoByUserId(t *testing.T){
	var tempUser user.User
	result :=config.DB.Model(&UserInfo{}).Where("id = ?", 1).Find(&tempUser)
	t.Log(result)
	if result.Error != nil{
		t.Error(result.Error)
	}

}
