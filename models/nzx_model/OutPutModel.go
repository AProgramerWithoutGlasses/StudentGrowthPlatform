package nzx_model

import "time"

type Out struct {
	Username     string `json:"username"`
	Name         string `json:"name"`
	Content      string `json:"content"`
	UserHeadshot string `json:"user_headshot"`
	PostTime     string `json:"post_time"`
	IsRead       bool   `json:"is_read"`
	Type         int    `json:"type"`
	ArticleId    uint   `json:"article_id"`
	MsgId        uint   `json:"msg_id"`
	//CreatedAt    time.Time `json:"-"`
}

//type Outs []Out
//
//func (a Outs) Len() int           { return len(a) }
//func (a Outs) Less(i, j int) bool { return a[i].CreatedAt.After(a[j].CreatedAt) } // 逆序排序
//func (a Outs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type CommentMsg struct {
	Username     string    `json:"username"`
	Name         string    `json:"name"`
	Content      string    `json:"content"`
	UserHeadshot string    `json:"user_headshot"`
	PostTime     string    `json:"post_time"`
	IsRead       bool      `json:"is_read"`
	Type         int       `json:"type"`
	ArticleId    uint      `json:"article_id"`
	MsgId        uint      `json:"msg_id"`
	CreatedAt    time.Time `json:"-"`
}

type CommentMsgs []CommentMsg

func (a CommentMsgs) Len() int           { return len(a) }
func (a CommentMsgs) Less(i, j int) bool { return a[i].CreatedAt.After(a[j].CreatedAt) } // 逆序排序
func (a CommentMsgs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
