package model

type Repost struct {
	UserID         int `json:"user_id" gorm:"primaryKey;foreignKey:UserID;references:UserID;constraint:OnDelete:CASCADE"`
	RepostedPostID int `json:"reposted_post_id" gorm:"primaryKey;foreignKey:RepostedPostID;references:PostID;constraint:OnDelete:CASCADE"`
}
