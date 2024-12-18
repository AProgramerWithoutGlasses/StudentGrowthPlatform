package nzx_model

// RedisLikeArticleData 文章点赞结构体
type RedisLikeArticleData struct {
	Operator string //操作：like/cancel_like
	Username string // 用户id
	Aid      int    // 文章id
}

// RedisLikeCommentData 评论点赞结构体
type RedisLikeCommentData struct {
	Operator string //操作:like/cancel_like
	Username string
	Cid      int
}
