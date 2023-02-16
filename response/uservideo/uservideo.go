package uservideo


type UserFavoriteResponse struct {
	StatusCode int32 `form:"status_code" json:"status_code"`
	StatusMsg string `form:"status_msg" json:"status_msg"`
}

