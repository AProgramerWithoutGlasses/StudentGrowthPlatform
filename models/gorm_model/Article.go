package gorm_model

import (
	"gorm.io/gorm"
	"studentGrow/models/constant"
)

type Article struct {
	gorm.Model
	Content       string `gorm:"size:350"json:"article_content"`
	WordCount     int    `gorm:"default:0"json:"word_count"`
	Pic           string
	Video         string
	Topic         string              `json:"article_topic"`
	Status        bool                `gorm:"not null"json:"article_status"`
	ReadAmount    int                 `gorm:"default:0"json:"read_amount"`
	LikeAmount    int                 `gorm:"default:0"json:"Like_amount"`
	CollectAmount int                 `gorm:"default:0"json:"collect_amount"`
	CommentAmount int                 `gorm:"default:0"json:"comment_amount"`
	ReportAmount  int                 `gorm:"default:0"json:"report_amount"`
	Point         int                 `gorm:"default:3"json:"point"`
	Quality       int                 `gorm:"default:0" json:"article_quality"`
	Ban           bool                `gorm:"default:false"json:"-"`
	UserID        uint                `gorm:"not null"` //文章属于用户
	User          User                `json:"user"`     //文章属于用户
	Comments      []Comment           //文章拥有评论
	ArticleLikes  []UserLikeRecord    //文章拥有点赞
	ArticleTags   []ArticleTag        `json:"article_tags"` //文章拥有标签
	ArticlePics   []ArticlePic        `json:"article_pics"` //文章拥有图片
	Collects      []UserCollectRecord //文章拥有收藏
	IsLike        bool                `gorm:"-" json:"is_like"`
	IsCollect     bool                `gorm:"-" json:"is_collect"`
	PostTime      string              `gorm:"-" json:"post_time"`
}

// Articles 自定义加权排序
type Articles []Article

func (a Articles) Len() int      { return len(a) }
func (a Articles) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a Articles) Less(i, j int) bool {
	return float64(a[i].LikeAmount)*constant.LikeWeightConstant+float64(a[i].CollectAmount)*constant.CollectWeightConstant+float64(a[i].CommentAmount)*constant.CommentWeightConstant > float64(a[j].LikeAmount)*constant.LikeWeightConstant+float64(a[j].CollectAmount)*constant.CollectWeightConstant+float64(a[j].CommentAmount)*constant.CommentWeightConstant
}
