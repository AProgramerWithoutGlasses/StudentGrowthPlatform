package gorm_model

import "gorm.io/gorm"

type UserLikeRecord struct {
	gorm.Model
	UserID    uint
	User      User
	ArticleID uint
	Article   Article
	CommentID uint
	Comment   Comment
	IsRead    bool `json:"is_read"`
	Type      int  `json:"type"` // 区分文章点赞还是评论点赞:0文章，1评论
}
