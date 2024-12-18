package nzx_model

// RedisCollectData 文章收藏结构体
type RedisCollectData struct {
	Operator string //操作：collect/cancel_collect
	Username string
	Aid      int
}
