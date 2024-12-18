package article

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"studentGrow/models/constant"
	myErr "studentGrow/pkg/error"
	res "studentGrow/pkg/response"
	"studentGrow/service/article"
	readUtil "studentGrow/utils/readMessage"
	"studentGrow/utils/token"
)

// GetArticleIdController article_id	获取文章详情
func GetArticleIdController(c *gin.Context) {
	in := struct {
		ArticleId int    `json:"article_id"`
		Username  string `json:"username"`
	}{}

	err := c.ShouldBindJSON(&in)
	if err != nil {
		zap.L().Error("GetArticleIdController() controller.article.ShouldBindJSON err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	// 获取文章详情
	art, err := article.GetArticleService(in.Username, in.ArticleId)
	if err != nil {
		zap.L().Error("GetArticleIdController() controller.article.GetArticleService err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	var tags []string
	for _, tag := range art.ArticleTags {
		tags = append(tags, tag.Tag.TagName)
	}
	var pics []string
	for _, pic := range art.ArticlePics {
		pics = append(pics, pic.Pic)
	}

	articleContent := map[string]any{
		"article_image": pics,
		"article_text":  art.Content,
		"article_video": art.Video,
	}
	// 若是不可被查看的文章除了文章状态，其余字段返回nil
	var data map[string]any
	if (art.Ban == true && art.User.Username != in.Username) || art.Status == false {
		data = map[string]any{
			"ban":             art.Ban,
			"status":          art.Status,
			"user_headshot":   nil,
			"name":            nil,
			"username":        nil,
			"user_class":      nil,
			"article_tags":    nil,
			"post_time":       nil,
			"article_content": nil,
			"like_amount":     nil,
			"collect_amount":  nil,
			"comment_amount":  nil,
			"is_like":         nil,
			"is_collect":      nil,
		}
	} else {
		data = map[string]any{
			"ban":             art.Ban,
			"status":          art.Status,
			"user_headshot":   art.User.HeadShot,
			"name":            art.User.Name,
			"username":        art.User.Username,
			"user_class":      art.User.Class,
			"article_tags":    tags,
			"post_time":       art.PostTime,
			"article_content": articleContent,
			"like_amount":     art.LikeAmount,
			"collect_amount":  art.CollectAmount,
			"comment_amount":  art.CommentAmount,
			"is_like":         art.IsLike,
			"is_collect":      art.IsCollect,
			"article_quality": art.Quality,
		}
	}

	res.ResponseSuccess(c, data)
}

// GetArticleListController 获取文章列表
func GetArticleListController(c *gin.Context) {
	in := struct {
		Page     int    `json:"page"`
		Limit    int    `json:"limit"`
		SortType string `json:"sort"`
		Order    string `json:"order"`
		StartAt  string `json:"start_at"`
		EndAt    string `json:"end_at"`
		IsBan    bool   `json:"article_ban"`
		Name     string `json:"name"`
		Topic    string `json:"topic"`
		KeyWords string `json:"key_words"`
	}{}

	err := c.ShouldBindJSON(&in)
	if err != nil {
		zap.L().Error("GetArticleListController() controller.article.ShouldBindJSON err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	aToken := token.NewToken(c)
	user, exist := aToken.GetUser()
	if !exist {
		res.ResponseError(c, res.TokenError)
		zap.L().Error("token错误")
		return
	}

	// 获取身份
	role, err := aToken.GetRole()
	if err != nil {
		zap.L().Error("GetArticleListController() controller.article.GetRole err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	//查询文章列表
	result, articleAmount, err := article.GetArticleListService(in.Page, in.Limit, in.SortType, in.Order, in.StartAt, in.EndAt, in.Topic, in.KeyWords, in.Name, in.IsBan, role, user.Username)
	if err != nil {
		zap.L().Error("GetArticleListController() controller.article.GetArticleListService err=", zap.Error(myErr.DataFormatError()))
		myErr.CheckErrors(err, c)
		return
	}

	list := make([]map[string]any, 0)
	for _, val := range result {
		list = append(list, map[string]any{
			"article_id":      val.ID,
			"article_content": val.Content,
			"user_headshot":   val.User.HeadShot,
			"article_ban":     val.Ban,
			"upvote_amount":   val.LikeAmount,
			"comment_amount":  val.CommentAmount,
			"username":        val.User.Username,
			"created_at":      val.CreatedAt,
			"name":            val.User.Name,
			"collect_amount":  val.CollectAmount,
			"article_quality": val.Quality,
		})
	}

	res.ResponseSuccess(c, map[string]any{
		"list":           list,
		"article_amount": articleAmount,
		"role":           role,
	})
}

// BannedArticleController 封禁文章
func BannedArticleController(c *gin.Context) {
	//获取前端发送的数据
	json, err := readUtil.GetJsonvalue(c)

	if err != nil {
		zap.L().Error("BannedArticleController() controller.article.GetJsonvalue err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	aToken := token.NewToken(c)
	user, exist := aToken.GetUser()
	if !exist {
		res.ResponseError(c, res.TokenError)
		zap.L().Error("token错误")
		return
	}
	username := user.Username
	role, err := aToken.GetRole()
	if err != nil {
		zap.L().Error("BannedArticleController() controller.article.GetRole err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	// 对应帖子进行封禁或解封操作
	err = article.BannedArticleService(json, role, username)
	// 检查错误
	if err != nil {
		zap.L().Error("BannedArticleController() controller.article.BannedArticleService err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	res.ResponseSuccess(c, struct{}{})

}

// DeleteArticleController 删除文章
func DeleteArticleController(c *gin.Context) {

	aToken := token.NewToken(c)
	user, exist := aToken.GetUser()
	if !exist {
		res.ResponseError(c, res.TokenError)
		zap.L().Error("token错误")
		return
	}

	username := user.Username
	role, err := aToken.GetRole()
	if err != nil {
		zap.L().Error("BannedArticleController() controller.article.GetRole err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	in := struct {
		ArticleId int `json:"article_id"`
	}{}

	err = c.ShouldBindJSON(&in)
	if err != nil {
		zap.L().Error("DeleteArticleController() controller.article.ShouldBindJSON err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	// 对文章进行删除操作
	err = article.DeleteArticleService(in.ArticleId, role, username)
	if err != nil {
		zap.L().Error("DeleteArticleController() controller.article.DeleteArticleService err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	res.ResponseSuccess(c, struct{}{})

}

// ReportArticle 举报文章
func ReportArticle(c *gin.Context) {
	//获取前端发送的数据
	json, err := readUtil.GetJsonvalue(c)
	if err != nil {
		zap.L().Error("ReportArticle() controller.article.GetJsonvalue err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	aToken := token.NewToken(c)
	user, exist := aToken.GetUser()
	if !exist {
		res.ResponseError(c, res.TokenError)
		zap.L().Error("token错误")
		return
	}
	username := user.Username

	// 对文章进行举报并记录
	err = article.ReportArticleService(json, username)

	if err != nil {
		zap.L().Error("ReportArticle() controller.article.ReportArticleService err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	res.ResponseSuccess(c, struct{}{})
}

// GetHotArticlesOfDayController 获取今日十条热帖
func GetHotArticlesOfDayController(c *gin.Context) {
	//获取前端发送的数据
	json, err := readUtil.GetJsonvalue(c)
	if err != nil {
		zap.L().Error("GetHotArticlesOfDayController() controller.article.GetJsonvalue err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	articles, err := article.SearchHotArticlesOfDayService(json)
	if err != nil {
		zap.L().Error("GetHotArticlesOfDayController() controller.article.SearchHotArticlesOfDayService err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	list := make([]map[string]any, 0)
	for _, a := range articles {
		list = append(list, map[string]any{
			"article_id":    a.ID,
			"article_title": a.Content,
		})
	}

	res.ResponseSuccess(c, map[string]any{
		"article_list": list,
	})
}

// SelectArticleAndUserListByPageFirstPageController 前台首页模糊搜索文章列表
func SelectArticleAndUserListByPageFirstPageController(c *gin.Context) {
	in := struct {
		Username string `json:"username"`
		KeyWords string `json:"key_word"`
		Topic    string `json:"topic_name"`
		SortWay  string `json:"article_sort"`
		Limit    int    `json:"article_count"`
		Page     int    `json:"article_page"`
	}{}

	err := c.ShouldBindJSON(&in)
	if err != nil {
		zap.L().Error("SelectArticleAndUserListByPageFirstPageController() controller.article.ShouldBindJSON err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	articles, err := article.SelectArticleAndUserListByPageFirstPageService(in.Username, in.KeyWords, in.Topic, in.SortWay, in.Limit, in.Page)
	if err != nil {
		zap.L().Error("SelectArticleAndUserListByPageFirstPageController() controller.article.SelectArticleAndUserListByPageFirstPageService err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	list := make([]map[string]any, 0)
	for _, a := range articles {
		var pics []string
		var tags []string
		for _, pic := range a.ArticlePics {
			pics = append(pics, pic.Pic)
		}
		for _, tag := range a.ArticleTags {
			tags = append(tags, tag.Tag.TagName)
		}

		list = append(list, map[string]any{
			"user_headshot":   a.User.HeadShot,
			"user_class":      a.User.Class,
			"name":            a.User.Name,
			"article_id":      a.ID,
			"like_amount":     a.LikeAmount,
			"collect_amount":  a.CollectAmount,
			"comment_amount":  a.CommentAmount,
			"article_content": a.Content,
			"article_pics":    pics,
			"article_video":   a.Video,
			"article_tags":    tags,
			"article_topic":   a.Topic,
			"is_like":         a.IsLike,
			"is_collect":      a.IsCollect,
			"post_time":       a.PostTime,
			"username":        a.User.Username,
			"article_quality": a.Quality,
			"user_identity":   a.User.Identity,
		})
	}

	res.ResponseSuccess(c, map[string]any{
		"content": list,
	})
}

// GetArticleByClassController 班级分类获取文章列表
func GetArticleByClassController(c *gin.Context) {
	input := struct {
		Username string `json:"username"`
		KeyWords string `json:"key_word"`
		SortWay  string `json:"article_sort"`
		Class    string `json:"class_name"`
		Limit    int    `json:"article_count"`
		Page     int    ` json:"article_page"`
	}{}

	err := c.ShouldBindJSON(&input)
	if err != nil {
		zap.L().Error("GetArticleByClassController() controller.article.ShouldBindJSON err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	// 获取列表
	articles, err := article.GetArticlesByClassService(input.KeyWords, input.Username, input.SortWay, input.Limit, input.Page, input.Class)
	if err != nil {
		zap.L().Error("GetArticleByClassController() controller.article.GetArticlesByClassService err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	list := make([]map[string]any, 0)
	for _, a := range articles {
		var pics []string
		var tags []string
		for _, pic := range a.ArticlePics {
			pics = append(pics, pic.Pic)
		}
		for _, tag := range a.ArticleTags {
			tags = append(tags, tag.Tag.TagName)
		}

		list = append(list, map[string]any{
			"user_headshot":   a.User.HeadShot,
			"user_class":      a.User.Class,
			"name":            a.User.Name,
			"article_id":      a.ID,
			"like_amount":     a.LikeAmount,
			"collect_amount":  a.CollectAmount,
			"comment_amount":  a.CommentAmount,
			"article_content": a.Content,
			"article_pics":    pics,
			"article_video":   a.Video,
			"article_tags":    tags,
			"article_topic":   a.Topic,
			"is_like":         a.IsLike,
			"is_collect":      a.IsCollect,
			"post_time":       a.PostTime,
			"username":        a.User.Username,
			"article_quality": a.Quality,
		})
	}
	res.ResponseSuccess(c, map[string]any{
		"content": list,
	})
}

// PublishArticleController 发布文章
func PublishArticleController(c *gin.Context) {
	aToken := token.NewToken(c)
	user, exist := aToken.GetUser()
	if !exist {
		res.ResponseError(c, res.TokenError)
		zap.L().Error("token错误")
		return
	}
	username := user.Username

	err := c.Request.ParseMultipartForm(constant.MemoryLimit) // 最大 80MB

	if err != nil {
		zap.L().Error("PublishArticleController() controller.article.getArticle.ParseMultipartForm err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		zap.L().Error("PublishArticleController() controller.article.getArticle.MultipartForm err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	// 获取基本数据
	content := form.Value["article_content"][0]
	countString := form.Value["word_count"][0]
	var wordCount int
	if countString != "" {
		wordCount, err = strconv.Atoi(countString)
		if err != nil {
			zap.L().Error("PublishArticleController() controller.article.getArticle.Atoi err=", zap.Error(err))
			myErr.CheckErrors(err, c)
			return
		}
	}

	tags := form.Value["article_tags"]
	topic := form.Value["article_topic"][0]

	statusString := form.Value["article_status"][0]
	var status bool

	if statusString != "" {
		status, err = strconv.ParseBool(statusString)
		if err != nil {
			zap.L().Error("PublishArticleController() controller.article.getArticle.ParseBool err=", zap.Error(err))
			myErr.CheckErrors(err, c)
			return
		}
	}

	// 获取图片和视频文件
	pics := form.File["pic"]
	video := form.File["video"]

	err = article.PublishArticleService(username, content, topic, wordCount, tags, pics, video, status)
	if err != nil {
		zap.L().Error("PublishArticleController() controller.article.getArticle.PublishArticleService err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	res.ResponseSuccess(c, struct{}{})

}

// ReviseArticleStatusController 修改文章私密状态
func ReviseArticleStatusController(c *gin.Context) {
	in := struct {
		ArticleId int  `json:"article_id"`
		Status    bool `json:"article_status"`
	}{}

	err := c.ShouldBindJSON(&in)
	if err != nil {
		zap.L().Error("ReviseArticleStatusController() controller.article.ShouldBindJSON err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	err = article.ReviseArticleStatusService(in.ArticleId, in.Status)
	if err != nil {
		zap.L().Error("ReviseArticleStatusController() controller.article.ReviseArticleStatusService err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	res.ResponseSuccess(c, struct{}{})
}

// SelectGoodArticleController 评选优秀帖子
func SelectGoodArticleController(c *gin.Context) {
	in := struct {
		ArticleId      int `json:"article_id"`
		ArticleQuality int `json:"article_quality"`
	}{}

	err := c.ShouldBindJSON(&in)
	if err != nil {
		zap.L().Error("SelectGoodArticleController() controller.article.ShouldBindJSON err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	aToken := token.NewToken(c)
	user, exist := aToken.GetUser()
	if !exist {
		res.ResponseError(c, res.TokenError)
		zap.L().Error("token错误")
		return
	}
	username := user.Username

	// 获取身份-token
	role, err := aToken.GetRole()
	if err != nil {
		zap.L().Error("SelectGoodArticleController() controller.article.GetRole err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	// 修改文章优秀等级
	err = article.SelectGoodArticleService(in.ArticleId, in.ArticleQuality, role, username)
	if err != nil {
		zap.L().Error("SelectGoodArticleController() controller.article.SelectGoodArticleService err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	res.ResponseSuccess(c, struct{}{})
}

// AdvancedArticleFilteringController 高级筛选文章
func AdvancedArticleFilteringController(c *gin.Context) {
	in := struct {
		KeyWords     string   `json:"key_words"`
		TopicName    string   `json:"topic_name"`
		ArticleCount int      `json:"article_count"`
		ArticlePage  int      `json:"article_page"`
		Username     string   `json:"username"`
		Sort         string   `json:"sort"`
		Order        string   `json:"order"`
		StartAt      string   `json:"start_at"`
		EndAt        string   `json:"end_at"`
		Class        []string `json:"class"`
		Name         string   `json:"name"`
		Grade        int      `json:"grade"`
		Role         string   `json:"role"`
	}{}

	err := c.ShouldBindJSON(&in)
	if err != nil {
		zap.L().Error("AdvancedArticleFilteringController() controller.article.ShouldBindJSON err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	articles, err := article.AdvancedArticleFilteringService(in.ArticlePage, in.ArticleCount, in.Sort, in.Order, in.StartAt, in.EndAt, in.TopicName, in.KeyWords, in.Name, in.Username, in.Class, in.Grade, in.Role)
	if err != nil {
		zap.L().Error("AdvancedArticleFilteringController() controller.article.AdvancedArticleFilteringService err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	content := make([]map[string]any, 0)
	for _, at := range articles {
		var pics []string
		var tags []string
		for _, pic := range at.ArticlePics {
			pics = append(pics, pic.Pic)
		}
		for _, tag := range at.ArticleTags {
			tags = append(tags, tag.Tag.TagName)
		}
		fmt.Println(pics)
		content = append(content, map[string]any{
			"user_headshot":   at.User.HeadShot,
			"user_class":      at.User.Class,
			"name":            at.User.Name,
			"article_id":      at.ID,
			"like_amount":     at.LikeAmount,
			"collect_amount":  at.CollectAmount,
			"comment_amount":  at.CommentAmount,
			"article_content": at.Content,
			"article_pic":     pics,
			"article_video":   at.Video,
			"article_tags":    tags,
			"article_topic":   at.Topic,
			"is_like":         at.IsLike,
			"is_collect":      at.IsCollect,
			"post_time":       at.PostTime,
			"username":        at.User.Username,
			"article_quality": at.Quality,
			"user_identity":   at.User.Identity,
		})
	}
	res.ResponseSuccess(c, map[string]any{
		"content": content,
	})
}

// GetGoodArticlesController 获取优秀帖子
func GetGoodArticlesController(c *gin.Context) {
	in := struct {
		Page     int    `json:"page"`
		Limit    int    `json:"limit"`
		Sort     string `json:"sort"`
		Order    string `json:"order"`
		StartAt  string `json:"start_at"`
		EndAt    string `json:"end_at"`
		Topic    string `json:"topic"`
		KeyWords string `json:"key_words"`
		Name     string `json:"name"`
	}{}

	err := c.ShouldBindJSON(&in)
	if err != nil {
		zap.L().Error("GetGoodArticlesController() controller.article.QueryGoodArticlesByRoleService err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	aToken := token.NewToken(c)
	user, exist := aToken.GetUser()
	if !exist {
		res.ResponseError(c, res.TokenError)
		zap.L().Error("token错误")
		return
	}
	username := user.Username

	// 获取身份-token
	role, err := aToken.GetRole()
	if err != nil {
		zap.L().Error("SelectGoodArticleController() controller.article.GetRole err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}
	// 查询优秀帖子
	articles, count, err := article.QueryGoodArticlesByRoleService(role, in.StartAt, in.EndAt, in.Topic, in.KeyWords, in.Sort, in.Order, in.Name, username, in.Page, in.Limit)
	if err != nil {
		zap.L().Error("GetGoodArticlesController() controller.article.QueryGoodArticlesByRoleService err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	list := make([]map[string]any, 0)
	for _, at := range articles {
		list = append(list, map[string]any{
			"article_id":      at.ID,
			"article_content": at.Content,
			"user_headshot":   at.User.HeadShot,
			"upvote_amount":   at.LikeAmount,
			"comment_amount":  at.CommentAmount,
			"username":        at.User.Username,
			"created_at":      at.CreatedAt,
			"collect_amount":  at.CollectAmount,
			"article_quality": at.Quality,
			"name":            at.User.Name,
		})
	}

	res.ResponseSuccess(c, map[string]any{
		"list":           list,
		"article_amount": count,
		"role":           role,
	})
}
