package social

import "TinyTolk/response/user"

type RelationActionResponse struct {
	
	StatusCode int32 `json:"status_code"`
	StatusMsg string `json:"status_msg"`
}

type RelationFollowListResponse struct {
	
	StatusCode int32 `json:"status_code"`
	StatusMsg string `json:"status_msg"`
	UserList []user.User `json:"user_list"`
}

type RelationFollowerListResponse struct {

	StatusCode int32 `json:"status_code"`
	StatusMsg string `json:"status_msg"`
	UserList []user.User `json:"user_list"`
}
