package gorm_model

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	Content      string `json:"comment_content"`
	LikeAmount   int    `gorm:"default:0"json:"comment_like_num"`
	IsRead       bool   `gorm:"default:false"json:"is_read"`
	UserID       uint   `gorm:"not null"`
	User         User
	Pid          uint             `json:"pid" gorm:"default:0"`     //回复评论的ID
	ReplyCount   int              `json:"comment_son_num" gorm:"-"` // 评论的子评论数量
	ArticleID    uint             `gorm:"not null"`                 //评论属于文章
	Article      Article          //评论属于文章
	CommentLikes []UserLikeRecord //评论拥有点赞
	IsLike       bool             `gorm:"-"json:"comment_if_like"`
	Time         string           `gorm:"-" json:"comment_time"`
	PostTime     string           `gorm:"-" json:"post_time"`
}

// Comments 自定义加权排序
type Comments []Comment

func (c Comments) Len() int      { return len(c) }
func (c Comments) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c Comments) Less(i, j int) bool {
	return c[i].LikeAmount > c[j].LikeAmount // 降序
}

// ByCreatedAt 根据创建时间排序
type ByCreatedAt []Comment

func (a ByCreatedAt) Len() int           { return len(a) }
func (a ByCreatedAt) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCreatedAt) Less(i, j int) bool { return a[i].CreatedAt.After(a[j].CreatedAt) }
