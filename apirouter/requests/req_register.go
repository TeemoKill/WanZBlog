package requests

type RegisterRequest struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`

	Username string `form:"username" json:"username"`
}
