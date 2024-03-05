package object

type User struct {
	ID            int    `json:"id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	PasswordHash  string `json:"password"`
	CreatedAt     string `json:"created_at"`
	EmailVerified bool   `json:"email_verified"`
}
