package requests

type UserProfileRequest struct {
	UserUUID string `form:"user_uuid" json:"user_uuid"`
}
