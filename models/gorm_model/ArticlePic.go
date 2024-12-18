package gorm_model

import "gorm.io/gorm"

type ArticlePic struct {
	gorm.Model
	Pic       string `json:"article_pic"`
	ArticleID uint
	Article   Article
}
