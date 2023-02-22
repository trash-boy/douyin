package social

type RelationActionRequest struct {
	Token string `form:"token"`
	ToUserId string `form:"to_user_id"`
	ActionType string `form:"action_type"  validate:"required,ValidateActionType"`
}

type RelationFollowListRequest struct {

	UserId string `form:"user_id"`
	Token string `form:"token"`

}

type RelationFollowerListRequest struct {

	UserId string `form:"user_id"`
	Token string `form:"token"`

}
