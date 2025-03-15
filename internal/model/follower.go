package model

type Follower struct {
	FollowerID  int `json:"follower_id"`
	FollowingID int `json:"following_id"`
}
