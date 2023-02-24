package userfollow

import (
	"TinyTolk/config"
	"github.com/jinzhu/gorm"
)

type UserFollow struct {
	gorm.Model
	UserId   uint `gorm:""`
	FollowId uint `gorm:""`
	Status   bool `gorm:"DEFAULT:true"`
}

func CreateUserFollowTable() error {
	db := config.DB.AutoMigrate(&UserFollow{})
	return db.Error
}

func GetFollowIdListByUserId(userId uint) ([]uint, error) {
	var ans []uint
	var userFollow []UserFollow
	result := config.DB.Model(UserFollow{}).Where("user_id = ? and status = ?", userId, true).Select("follow_id").Find(&userFollow)
	for i := 0; i < len(userFollow); i++ {
		ans = append(ans, userFollow[i].FollowId)
	}
	return ans, result.Error
}

func GetFollowerIdListByUserId(userId uint) ([]uint, error) {
	var ans []uint
	var userFollow []UserFollow
	result := config.DB.Model(UserFollow{}).Where("follow_id = ? and status = ?", userId, true).Select("user_id").Find(&userFollow)
	for i := 0; i < len(userFollow); i++ {
		ans = append(ans, userFollow[i].UserId)
	}
	return ans, result.Error
}
func InsertNotExist(userId uint, followerId uint) error {
	exist, _ := RecordIsExist(userId, followerId)

	if !exist {
		return InsertFollow(userId, followerId)
	}
	return nil
}

func InsertFollow(userId uint, followerId uint) error {
	var userFollow UserFollow
	userFollow.UserId = userId
	userFollow.FollowId = followerId
	userFollow.Status = true

	result := config.DB.Create(&userFollow)
	return result.Error
}

func UpdateFollow(userId uint, followerId uint, status bool) error {

	result := config.DB.Model(&UserFollow{}).Where("user_id = ? and follow_id = ?", userId, followerId).Update("status", status)
	return result.Error
}

func DeleteFollow(userId uint, followerId uint) error {
	var userFollow UserFollow
	userFollow.UserId = userId
	userFollow.FollowId = followerId
	result := config.DB.Model(&userFollow).Update("status", false)
	return result.Error
}

func RecordIsExist(userId uint, followerId uint) (bool, error) {
	var userFollow UserFollow

	result := config.DB.Model(&UserFollow{}).Where("user_id = ? and follow_id = ?", userId, followerId).First(&userFollow)
	return result.RowsAffected != 0, result.Error
}

func UserIsFollow(userId, followId uint) (bool, error) {
	var userFollow UserFollow
	userFollow.UserId = userId
	userFollow.FollowId = followId

	result := config.DB.Find(&userFollow)
	if result.RowsAffected != 0 && userFollow.Status == true {
		return true, nil
	}
	return false, nil
}
