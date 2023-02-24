package DataBaseInit

import "TinyTolk/model/user"
import "TinyTolk/model/userfollow"
import "TinyTolk/model/comment"
import "TinyTolk/model/uservideoFavorite"
import "TinyTolk/model/Video"

func Init() {
	user.CreateUsersTable()
	userfollow.CreateUserFollowTable()
	comment.CreateCommentTable()
	uservideoFavorite.CreateUserVideoFavoriteTable()
	Video.CreateVideoTable()
}
