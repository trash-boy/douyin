package userfollow

import (
	"TinyTolk/config"
	"time"
)

type UserFollow struct {

	UserId uint `gorm:"primaryKey;autoIncrement:false"`
	FollowId uint `gorm:"primaryKey;autoIncrement:false"`
	Status bool `gorm:"DEFAULT:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func CreateUserFollowTable()error{
	db := config.DB.AutoMigrate(&UserFollow{})
	return db.Error
}
func UserIsFollow(userId, followId uint)(bool,error){
	var userFollow UserFollow
	userFollow.UserId = userId
	userFollow.FollowId = followId

	result := config.DB.Find(&userFollow)
	if result.RowsAffected != 0  && userFollow.Status == true{
		return true,nil
	}
	return false, nil
}

