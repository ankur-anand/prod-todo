package auth

// UserModel represents individual User
type UserModel struct {
	ID       int64  `json:"id,omitempty"`
	Email    string `json:"email"`
	Password string `json:"-,"`
	Username string `json:"username"`
}
