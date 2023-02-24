package Video

import (
	"TinyTolk/config"
	user3 "TinyTolk/model/user"
	"TinyTolk/response/video"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

type Video struct {
	gorm.Model
	UserId        uint
	User          user3.User `gorm:"ForeignKey:UserId"`
	PlayUrl       string
	CoverUrl      string
	FavoriteCount int64
	CommentCount  int64
	Title         string
}

func CreateVideoTable() error {
	db := config.DB.AutoMigrate(&Video{})
	return db.Error
}

func InsertVideo(userId uint, playUrl, coverUrl, title string) error {
	var video Video
	video.UserId = userId
	video.PlayUrl = playUrl
	video.CoverUrl = coverUrl
	video.Title = title

	result := config.DB.Create(&video)
	return result.Error
}

func GetVideoListByUserId(userId uint, videoList *[]video.Video) error {
	*videoList = make([]video.Video, 0)

	result := config.DB.Model(&Video{}).Where("user_id = ?", userId).Find(videoList)
	log.Print(userId)
	return result.Error
}

func GetVideoListByTimeStamp(timeStamp int64, videoList *[]video.Video) error {
	//将时间戳变为时间类型
	timeLayout := "2006-01-02 15:04:05"
	timeStr := time.Unix(timeStamp, 0).Format(timeLayout)
	*videoList = make([]video.Video, 0)
	result := config.DB.Model(&Video{}).Where("created_at <= ?", timeStr).Order("created_at desc").Limit(5).Find(videoList)
	return result.Error
}

func GetVideoByVideoId(VideoId uint, videoList *video.Video) error {
	result := config.DB.Model(&Video{}).Where("id = ?", VideoId).Scan(videoList)
	return result.Error
}

func GetUserIdByVideoId(videoID uint) uint {
	var v Video
	_ = config.DB.Model(&Video{}).Select("user_id").Where("id = ?", videoID).Find(&v)
	return v.UserId
}

func AddFavoriteCount(Id uint) error {
	result := config.DB.Model(&Video{}).Where(" id = ?", Id).UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1))
	return result.Error
}
func SubFavoriteCount(Id uint) error {
	result := config.DB.Model(&Video{}).Where(" id = ?", Id).UpdateColumn("favorite_count", gorm.Expr("favorite_count - ?", 1))
	return result.Error
}

func AddCommentCount(Id uint) error {
	result := config.DB.Model(&Video{}).Where(" id = ?", Id).UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1))
	return result.Error
}
func SubCommentCount(Id uint) error {
	result := config.DB.Model(&Video{}).Where(" id = ?", Id).UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1))
	return result.Error
}
