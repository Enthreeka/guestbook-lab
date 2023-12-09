package model

type Session struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Token  string `json:"token"`

	User User `json:"user"`
}
