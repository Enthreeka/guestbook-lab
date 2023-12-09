package model

type Guestbook struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	ListID    int    `json:"list_id"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`

	UserName string `json:"login"`
	ListName string `json:"name"`
}
