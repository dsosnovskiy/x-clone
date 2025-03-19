package model

// ПЕРЕНЕСТИ В post

type Repost struct {
	UserID int `json:"user_id"`
	PostID int `json:"post_id"`
}
