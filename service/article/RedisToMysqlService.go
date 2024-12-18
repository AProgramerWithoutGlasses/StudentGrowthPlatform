package article

import (
	"fmt"
	nzxMod "studentGrow/models/nzx_model"
	myErr "studentGrow/pkg/error"
)

var (
	ArticleLikeChan    chan nzxMod.RedisLikeArticleData
	CommentLikeChan    chan nzxMod.RedisLikeCommentData
	ArticleCollectChan chan nzxMod.RedisCollectData
)

func InitMyMQ() {
	ArticleLikeChan = make(chan nzxMod.RedisLikeArticleData, 100)
	CommentLikeChan = make(chan nzxMod.RedisLikeCommentData, 100)
	ArticleCollectChan = make(chan nzxMod.RedisCollectData, 100)
	go writeToMysql()
}

// likeType:0代表文章，1代表评论
func writeToMysql() {
	for {
		select {
		case articleLike := <-ArticleLikeChan:
			switch articleLike.Operator {
			case "like":
				err := LikeToMysql(articleLike.Aid, 0, articleLike.Username)
				if err != nil {
					fmt.Println("writeToMysql() dao.redis.LikeToMysql err=", err)
				}
			case "cancel_like":
				err := CancelLikeToMysql(articleLike.Aid, 0, articleLike.Username)
				if err != nil {
					fmt.Println("writeToMysql() dao.redis.CancelLikeToMysql err=", err)
				}
			default:
				fmt.Println("writeToMysql() err=", myErr.DataFormatError())
			}
		case commentLike := <-CommentLikeChan:
			switch commentLike.Operator {
			case "like":
				err := LikeToMysql(commentLike.Cid, 1, commentLike.Username)
				if err != nil {
					fmt.Println("writeToMysql() dao.redis.LikeToMysql err=", err)
				}
			case "cancel_like":
				err := CancelLikeToMysql(commentLike.Cid, 1, commentLike.Username)
				if err != nil {
					fmt.Println("writeToMysql() dao.redis.CancelLikeToMysql err=", err)
				}
			}
		case articleCollect := <-ArticleCollectChan:
			switch articleCollect.Operator {
			case "collect":
				err := CollectToMysql(articleCollect.Aid, articleCollect.Username)
				if err != nil {
					fmt.Println("writeToMysql() dao.redis.CollectToMysql err=", err)
				}
			case "cancel_collect":
				err := CancelCollectToMysql(articleCollect.Aid, articleCollect.Username)
				if err != nil {
					fmt.Println("writeToMysql() dao.redis.CancelCollectToMysql err=", err)
				}
			}
		default:
			//fmt.Println("等待redis数据传入...")
		}
	}
}
