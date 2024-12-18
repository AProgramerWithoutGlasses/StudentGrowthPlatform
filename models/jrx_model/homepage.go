package jrx_model

import "time"

// 个人主页信息
type HomepageMesStruct struct {
	Username     string  `json:"username"`
	Ban          bool    `json:"ban"`
	Name         string  `json:"name"`
	UserHeadShot string  `json:"user_headshot"`
	UserMotto    string  `json:"user_motto"`
	UserFans     int     `json:"user_fans"`
	UserConcern  int     `json:"user_concern"`
	UserLike     int     `json:"user_like"`
	Point        float64 `json:"user_point"`
	UserClass    string  `json:"user_class"`
}

// 个人资料信息
type HomepageDataStruct struct {
	Name         string `json:"name"`
	UserHeadShot string `json:"user_headshot"`
	UserMotto    string `json:"user_motto"`
	UserClass    string `json:"user_class"`
	UserGender   string `json:"user_gender"`
	PhoneNumber  string `json:"phone_number"`
	UserEmail    string `json:"user_email"`
	UserYear     string `json:"user_year"`
}

// 用于粉丝列表
type HomepageFanStruct struct {
	Username  string `json:"username"`
	Name      string `json:"name"`
	Motto     string `json:"user_motto"`
	HeadShot  string `json:"user_headshot"`
	IsConcern string `json:"is_concern"`
}

type HomepageArticleHistoryStruct struct {
	ID            int       `json:"article_id,omitempty,omitempty"`
	HeadShot      string    `json:"user_headshot,omitempty"`
	Name          string    `json:"name,omitempty"`
	Content       string    `json:"article_content,omitempty"`
	Pic           string    `json:"article_pic,omitempty"`
	CommentAmount int       `json:"comment_amount"`
	LikeAmount    int       `json:"like_amount"`
	CollectAmount int       `json:"collect_amount"`
	Topic         string    `json:"article_topic,omitempty"`
	Status        bool      `json:"article_status,omitempty"`
	ReadAmount    int       `json:"read_amount,omitempty"`
	ReportAmount  int       `json:"report_amount,omitempty"`
	Ban           bool      `json:"ban"`
	UserID        uint      `json:"user_id,omitempty"`
	IsLike        bool      `json:"is_like,omitempty"`
	IsCollect     bool      `json:"is_collect,omitempty"`
	CreatedAt     time.Time `json:"-"`
	PostTime      string    `json:"post_time,omitempty"`
	ArticleTags   []string  `json:"article_tags,omitempty"`
}

type HomepageClassmateStruct struct {
	Username string `json:"username"`
	HeadShot string `json:"user_headshot"`
	Name     string `json:"name"`
}

type HomepageTrack struct {
	ID              int    `json:"article_id"`
	Content         string `json:"article_content"`
	IType           string `json:"i_type"`
	Name            string `json:"name"`
	Created_at      string `json:"i_time"`
	LikeAmount      int    `json:"like_amount"`
	CommentAmount   int    `json:"comment_amount"`
	Comment_content string `json:"comment_content"`
}

type HomepageTopicPoint struct {
	StudyPoint     int     `json:"study_point"`
	HonorPoint     int     `json:"honor_point"`
	WorkPoint      int     `json:"work_point"`
	SocialPoint    int     `json:"social_point"`
	VolunteerPoint int     `json:"volunteer_point"`
	SportPoint     int     `json:"sport_point"`
	LifePoint      int     `json:"life_point"`
	TotalPoint     float64 `json:"total_point"`
}
