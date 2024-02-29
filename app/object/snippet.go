package object

type Snippet struct {
	ID        int    `json:"id"`
	Content   string `json:"content"`
	Language  string `json:"language"`
	CreatedAt string `json:"created_at"`
	ExpiresAt string `json:"expires_at"`
}
