package comment

import (
	"TinyTolk/config"
	"TinyTolk/model/Video"
	"TinyTolk/model/user"
	"TinyTolk/response/comment"
	"github.com/jinzhu/gorm"
)
//todo 回复评论
type Comment struct {
	gorm.Model
	Video Video.Video `gorm:"ForeignKey:VideoId"`
	VideoId uint
	User user.User `gorm:"ForeignKey:UserId"`
	UserId uint
	Content string
	Status bool  //表示是否删除， 默认为未删除，true表示未删除，false 表示已删除
}

func CreateCommentTable(){
	config.DB.AutoMigrate(Comment{})
}

func InsertComment(videoId uint, userId uint, content string)(*Comment,error){
	var c Comment
	c.VideoId = videoId
	c.UserId = userId
	c.Content = content
	c.Status = true
	result := config.DB.Create(&c)
	return &c, result.Error
}



func DeleteCommentById(commentId uint)error{
	result := config.DB.Model(&Comment{}).Where("id = ?", commentId).Update("status", false)
	return result.Error
}
func GetCommentByVideoId(videoId uint)(*[]comment.Comment, error){
	var response []comment.Comment
	result := config.DB.Model(&Comment{}).Where("video_id = ? and status = ?", videoId,1).Order("created_at desc").Scan(&response)
	return &response, result.Error
}
