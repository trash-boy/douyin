package user

import (
	"TinyTolk/config"
	userReponse "TinyTolk/response/user"
	"github.com/jinzhu/gorm"
	"log"
)

//用户信息表
type User struct {
	gorm.Model
	Username string `gorm:"UNIQUE" ` //账号名
	Password string `gorm:"NOT NULL"`
	Name string //用户名
	FollowCount int64 `gorm:"DEFAULT:0"`
	FollowerCount int64 `gorm:"DEFAULT:0"`
}

//用户关注表

//生成用户数据表

func CreateUsersTable()error{
	db := config.DB.AutoMigrate(&User{})
	return db.Error

}
func InsertUser(user *User)*gorm.DB{
	return config.DB.Create(user)
}

func GetUserByUsername(user *User, username string)bool{
	config.DB.First(user,"username = ?",username)
	log.Println("查询到的user",user)
	return user.ID != 0
}

func GetUserById(user *userReponse.User,id uint)bool{
	config.DB.Model(&User{}).Find(user, "id = ?", id)
	return user.ID != 0
}


//关注表相关操作
