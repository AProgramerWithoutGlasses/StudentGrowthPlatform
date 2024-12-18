package article

import (
	"errors"
	"fmt"
	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mime/multipart"
	"sort"
	"strconv"
	"studentGrow/dao/mysql"
	"studentGrow/dao/redis"
	"studentGrow/models/constant"
	model "studentGrow/models/gorm_model"
	myErr "studentGrow/pkg/error"
	"studentGrow/utils/fileProcess"
	"studentGrow/utils/timeConverter"
	"time"
)

// GetArticleService 获取文章详情
func GetArticleService(username string, aid int) (article *model.Article, err error) {
	// 若为游客
	if username == "passenger" {
		err, article = mysql.QueryArticleByIdOfPassenger(aid)
		if err != nil {
			zap.L().Error("GetArticleService() service.article.GetUserByUsername err=", zap.Error(err))
			return nil, err
		}
	} else {
		// 若为用户
		uid, err := mysql.QueryUserIdByUsername(username)
		if err != nil {
			zap.L().Error("GetArticleService() service.article.QueryUserIdByUsername err=", zap.Error(err))
			return nil, err
		}

		// 查询用户是否为管理员
		IsManager, err := mysql.QueryUserIsManager(uid)
		if err != nil {
			zap.L().Error("GetArticleService() service.article.QueryUserIsManager err=", zap.Error(err))
			return nil, err
		}

		// 查找文章作者username
		user, err := mysql.QueryUserByArticleId(aid)
		if err != nil {
			zap.L().Error("GetArticleService() service.article.QueryUserByArticleId err=", zap.Error(err))
			return nil, err
		}
		if IsManager || user.Username == username {
			// 若为管理员或者作者本人
			article, err = mysql.QueryArticleByIdOfManager(aid)
			if err != nil {
				zap.L().Error("GetArticleService() service.article.QueryArticleByIdOfManager err=", zap.Error(err))
				return nil, err
			}
		} else {
			// 若为普通用户
			err, article = mysql.QueryArticleById(aid, user.ID)
			if err != nil {
				zap.L().Error("GetArticleService() service.article.QueryArticleById err=", zap.Error(err))
				return nil, err
			}
		}
	}

	err = mysql.DB.Transaction(func(tx *gorm.DB) error {
		if username != "passenger" {
			// 查询是否点赞或收藏
			liked, err := redis.IsUserLiked(strconv.Itoa(aid), username, 0)
			if err != nil {
				zap.L().Error("GetArticleService() service.article.IsUserLiked err=", zap.Error(err))
				return err
			}
			article.IsLike = liked
			selected, err := redis.IsUserCollected(username, strconv.Itoa(aid))
			if err != nil {
				zap.L().Error("GetArticleService() service.article.IsUserCollected err=", zap.Error(err))
				return err
			}
			article.IsCollect = selected

			// 存储到浏览记录
			uid, err := mysql.SelectUserByUsername(username)
			if err != nil {
				zap.L().Error("GetArticleService() service.article.GetIdByUsername err=", zap.Error(err))
				return err
			}
			err = mysql.InsertReadRecord(uid, aid, tx)
			if err != nil {
				zap.L().Error("GetArticleService() service.article.InsertReadRecord err=", zap.Error(err))
				return err
			}
		}

		// 计算发布时间
		article.PostTime = timeConverter.IntervalConversion(article.CreatedAt)

		// 该文章阅读量+1
		err = UpdateArticleReadNumService(aid, 1, tx)
		if err != nil {
			zap.L().Error("GetArticleService() service.article.UpdateArticleReadNumService err=", zap.Error(err))
			return err
		}
		return nil
	})
	if err != nil {
		zap.L().Error("GetArticleService() service.article.Transaction err=", zap.Error(err))
		return nil, err
	}

	return article, err
}

// GetArticleListService 后台获取文章列表
func GetArticleListService(page, limit int, sortType, order, startAtString, endAtString, topic, keyWords, name string, isBan bool, role, username string) (articles []model.Article, count int, err error) {
	user, err := mysql.GetUserByUsername(username)
	if err != nil {
		zap.L().Error("GetArticleListService() service.article.GetUserByUsername err=", zap.Error(err))
		return nil, -1, err
	}

	switch role {
	case "class":
		// 查询该用户班级
		articles, err = mysql.QueryArticleAndUserListByPageForClass(page, limit, sortType, order, startAtString, endAtString, topic, keyWords, name, isBan, user.Class)
		count, err = mysql.QueryArticleNumByPageForClass(sortType, order, startAtString, endAtString, topic, keyWords, name, isBan, user.Class)
	case "grade1":
		articles, err = mysql.QueryArticleAndUserListByPageForGrade(page, limit, sortType, order, startAtString, endAtString, topic, keyWords, name, isBan, 1)
		count, err = mysql.QueryArticleNumByPageForGrade(sortType, order, startAtString, endAtString, topic, keyWords, name, isBan, 1)
	case "grade2":
		articles, err = mysql.QueryArticleAndUserListByPageForGrade(page, limit, sortType, order, startAtString, endAtString, topic, keyWords, name, isBan, 2)
		count, err = mysql.QueryArticleNumByPageForGrade(sortType, order, startAtString, endAtString, topic, keyWords, name, isBan, 2)
	case "grade3":
		articles, err = mysql.QueryArticleAndUserListByPageForGrade(page, limit, sortType, order, startAtString, endAtString, topic, keyWords, name, isBan, 3)
		count, err = mysql.QueryArticleNumByPageForGrade(sortType, order, startAtString, endAtString, topic, keyWords, name, isBan, 3)
	case "grade4":
		articles, err = mysql.QueryArticleAndUserListByPageForGrade(page, limit, sortType, order, startAtString, endAtString, topic, keyWords, name, isBan, 4)
		count, err = mysql.QueryArticleNumByPageForGrade(sortType, order, startAtString, endAtString, topic, keyWords, name, isBan, 4)
	case "college":
		articles, err = mysql.QueryArticleAndUserListByPageForSuperman(page, limit, sortType, order, startAtString, endAtString, topic, keyWords, name, isBan)
		count, err = mysql.QueryArticleNumByPageForSuperman(sortType, order, startAtString, endAtString, topic, keyWords, name, isBan)
	case "superman":
		articles, err = mysql.QueryArticleAndUserListByPageForSuperman(page, limit, sortType, order, startAtString, endAtString, topic, keyWords, name, isBan)
		count, err = mysql.QueryArticleNumByPageForSuperman(sortType, order, startAtString, endAtString, topic, keyWords, name, isBan)
	}

	if err != nil {
		zap.L().Error("GetArticleListService() service.article.SelectArticleAndUserListByPage err=", zap.Error(err))
		return nil, -1, err
	}
	return articles, count, nil
}

// BannedArticleService 解封或封禁文章
func BannedArticleService(j *jsonvalue.V, role string, username string) error {
	// 获取文章id和封禁状态
	aid, err := j.GetInt("article_id")
	if err != nil {
		zap.L().Error("AddTopicsService() service.article.GetInt err=", zap.Error(err))
		return err
	}
	isBan, err := j.GetBool("article_ban")
	if err != nil {
		zap.L().Error("AddTopicsService() service.article.GetBool err=", zap.Error(err))
		return err
	}

	// 权限验证
	if role == "" {
		return myErr.OverstepCompetence
	}

	err = mysql.DB.Transaction(func(tx *gorm.DB) error {
		switch role {
		case "class":
			err = mysql.BannedArticleByIdForClass(aid, isBan, username, tx)
		case "grade1":
			err = mysql.BannedArticleByIdForGrade(aid, 1, tx)
		case "grade2":
			err = mysql.BannedArticleByIdForGrade(aid, 2, tx)
		case "grade3":
			err = mysql.BannedArticleByIdForGrade(aid, 3, tx)
		case "grade4":
			err = mysql.BannedArticleByIdForGrade(aid, 4, tx)
		case "college":
			err = mysql.BannedArticleByIdForSuperman(aid, tx)
		case "superman":
			err = mysql.BannedArticleByIdForSuperman(aid, tx)
		default:
			return myErr.ErrNotFoundError
		}

		if err != nil {
			zap.L().Error("BannedArticleService() service.article.GetIdByUsername err=", zap.Error(err))
			return err
		}

		/*
			若举报信箱存在该文章id，则标记该信息已读
		*/

		ok, err := mysql.QueryIsExistArticleIdByReportMsg(aid)

		if ok {
			err = mysql.DeleteArticleReportMsg(aid, tx)
			if err != nil {
				zap.L().Error("BannedArticleService() service.article.DeleteArticleReportMsg err=", zap.Error(err))
				return err
			}
		}

		// 将封禁信息通过系统通知发送给作者本人
		content, err := mysql.QueryContentByArticleId(aid)
		if err != nil {
			zap.L().Error("BannedArticleService() service.article.DeleteArticleReportMsg err=", zap.Error(err))
			return err
		}

		tarUser, err := mysql.QueryUserByArticleId(aid)
		if err != nil {
			zap.L().Error("BannedArticleService() service.article.DeleteArticleReportMsg err=", zap.Error(err))
			return err
		}

		ownUserId, err := mysql.QueryUserIdByUsername(username)
		if err != nil {
			zap.L().Error("BannedArticleService() service.article.QueryUserIdByUsername err=", zap.Error(err))
			return err
		}

		msg := fmt.Sprintf("您的内容为:<br/>%s<br/>已被封禁!", content)
		err = mysql.AddBanSystemNotification(msg, ownUserId, tx, int(tarUser.ID))
		if err != nil {
			zap.L().Error("BannedArticleService() service.article.DeleteArticleReportMsg err=", zap.Error(err))
			return err
		}

		/*
			回滚分数
		*/
		// 回滚分数的条件：1. point>0;文章处于封禁或私密的状态	2. point<0;文章处于非封禁或私密的状态

		// 通过文章id获取文章
		article, err := mysql.QueryArticleByIdOfManager(aid)
		if err != nil {
			zap.L().Error("BannedArticleService() service.article.QueryArticleByIdOfManager err=", zap.Error(err))
			return err
		}

		curPoint, err := mysql.QueryArticlePoint(aid)
		if err != nil {
			zap.L().Error("BannedArticleService() service.article.DeleteArticleReportMsg err=", zap.Error(err))
			return err
		}

		// 文章且文章状态为公开,则减去分数
		if isBan && article.Status {
			curPoint = -curPoint
			err = UpdatePointByUsernamePointAid(username, curPoint, aid, tx)
			if err != nil {
				zap.L().Error("BannedArticleService() service.article.UpdatePointByUsernamePointAid err=", zap.Error(err))
				return err
			}
		}
		// 解封文章且文章状态为公开，则加上分数
		if !isBan && article.Status {
			err = UpdatePointByUsernamePointAid(username, curPoint, aid, tx)
			if err != nil {
				zap.L().Error("BannedArticleService() service.article.UpdatePointByUsernamePointAid err=", zap.Error(err))
				return err
			}
		}

		// 其他情况，无需理会

		return nil
	})
	if err != nil {
		zap.L().Error("BannedArticleService() service.article.Transaction err=", zap.Error(err))
		return err
	}

	return nil
}

// DeleteArticleService 删除文章
func DeleteArticleService(aid int, role string, username string) error {

	article, err := mysql.QueryArticleByIdOfManager(aid)
	if err != nil {
		zap.L().Error("DeleteArticleService() service.article.QueryArticleByIdOfManager err=", zap.Error(err))
		return err
	}

	curPoint, err := mysql.QueryArticlePoint(aid)
	if err != nil {
		zap.L().Error("DeleteArticleService() service.article.QueryArticlePoint err=", zap.Error(err))
		return err
	}

	// 权限验证
	if role == "" && article.User.Username != username {
		return myErr.OverstepCompetence
	}

	err = mysql.DB.Transaction(func(tx *gorm.DB) error {
		/*
			回滚分数
		*/

		// 条件: 删除文章且文章状态为公开且未被www

		if !article.Ban && article.Status {
			err = UpdatePointByUsernamePointAid(username, -curPoint, aid, tx)
			if err != nil {
				zap.L().Error("DeleteArticleService() service.article.UpdatePointByUsernamePointAid err=", zap.Error(err))
				return err
			}
		}

		/*
			删除文章
		*/

		// 若为本人删除
		if article.User.Username == username {
			err = mysql.DeleteArticleByIdForSuperman(aid, tx)
		} else {
			// 若为管理员删除
			switch role {
			case "class":
				err = mysql.DeleteArticleByIdForClass(aid, username, tx)
			case "grade1":
				err = mysql.DeleteArticleByIdForGrade(aid, 1, tx)
			case "grade2":
				err = mysql.DeleteArticleByIdForGrade(aid, 2, tx)
			case "grade3":
				err = mysql.DeleteArticleByIdForGrade(aid, 3, tx)
			case "grade4":
				err = mysql.DeleteArticleByIdForGrade(aid, 4, tx)
			case "college":
				err = mysql.DeleteArticleByIdForSuperman(aid, tx)
			case "superman":
				err = mysql.DeleteArticleByIdForSuperman(aid, tx)
			default:
				return myErr.ErrNotFoundError
			}
		}

		if err != nil {
			zap.L().Error("DeleteArticleService() service.article.DeleteArticleByIdForClass err=", zap.Error(err))
			return err
		}

		/*
			删除文章的关联图片表
		*/
		err = mysql.DeleteArticlePicByArticleId(aid)
		if err != nil {
			zap.L().Error("DeleteArticleService() service.article.DeleteArticlePicByArticleId err=", zap.Error(err))
			return err
		}

		/*
			删除文章的标签关联表
		*/
		err = mysql.DeleteArticleTagByArticleId(aid)
		if err != nil {
			zap.L().Error("DeleteArticleService() service.article.DeleteArticleTagByArticleId err=", zap.Error(err))
			return err
		}

		return nil
	})
	if err != nil {
		zap.L().Error("DeleteArticleService() service.article.Transaction err=", zap.Error(err))
		return err
	}

	// 删除redis
	redis.RDB.HDel("article", strconv.Itoa(aid))
	redis.RDB.HDel("collect", strconv.Itoa(aid))

	return nil
}

// ReportArticleService 举报文章
func ReportArticleService(j *jsonvalue.V, username string) error {
	// 获取文章id和举报者用户id,举报信息
	aid, err := j.GetInt("article_id")
	if err != nil {
		zap.L().Error("AddTopicsService() service.article.SearchHotArticlesOfDay.GetInt err=", zap.Error(err))
		return err
	}

	reportMsg, err := j.GetString("report_msg")

	// 通过username查询id
	uid, err := mysql.SelectUserByUsername(username)
	if err != nil {
		zap.L().Error("AddTopicsService() service.article.SearchHotArticlesOfDay.SelectUserByUsername err=", zap.Error(err))
		return err
	}

	err = mysql.ReportArticleById(aid, uid, reportMsg)
	if err != nil {
		zap.L().Error("AddTopicsService() service.article.SearchHotArticlesOfDay.ReportArticleById err=", zap.Error(err))
		return err
	}
	return nil
}

// SearchHotArticlesOfDayService 获取今日十条热帖
func SearchHotArticlesOfDayService(j *jsonvalue.V) (model.Articles, error) {

	// 获取热帖条数
	count, err := j.GetInt("article_count")
	if err != nil {
		zap.L().Error("AddTopicsService() service.article.ArticleService.GetInt err=", zap.Error(err))
		return nil, err
	}
	// 计算今日的始末时间
	startOfDay := time.Now().Truncate(24 * time.Hour) // 今天的开始时间
	endOfDay := startOfDay.Add(24 * time.Hour)        // 明天的开始时间

	articles, err := mysql.SearchHotArticlesOfDay(startOfDay, endOfDay)
	if err != nil {
		zap.L().Error("AddTopicsService() service.article.SearchHotArticlesOfDay.GetInt err=", zap.Error(err))
		return nil, err
	}
	// 排序
	sort.Sort(articles)

	var list model.Articles

	// 获取热度前count条数据
	for i := 0; i < count && i < len(articles); i++ {
		list = append(list, articles[i])
	}

	return list, nil
}

// QueryArticleStatusAndBanById 查询文章的私密、封禁状态
func QueryArticleStatusAndBanById(aid int) (bool, error) {
	isBan, err := mysql.QueryIsBanByArticleId(aid)
	if err != nil {
		zap.L().Error("QueryArticleStatusAndBanById() service.article.QueryIsBanByArticleId.GetInt err=", zap.Error(err))
		return false, err
	}
	status, err := mysql.QueryArticleStatusById(aid)
	if err != nil {
		zap.L().Error("QueryArticleStatusAndBanById() service.article.QueryArticleStatusById.GetInt err=", zap.Error(err))
		return false, err
	}
	if !isBan && status {
		return true, nil
	}
	return false, nil
}

// SelectArticleAndUserListByPageFirstPageService 前台首页模糊查询文章列表
func SelectArticleAndUserListByPageFirstPageService(username, keyWords, topic, SortWay string, limit, page int) ([]model.Article, error) {
	// 查询符合模糊搜索的文章集合
	var articles model.Articles
	var err error
	if topic == "全部话题" {
		articles, err = mysql.QueryArticleAndUserListByPageFirstPage(keyWords, limit, page)
		if err != nil {
			zap.L().Error("SelectArticleAndUserListByPageFirstPageService() service.article.QueryArticleAndUserListByPageFirstPage err=", zap.Error(err))
			return nil, err
		}
	} else {
		articles, err = mysql.QueryArticleAndUserListByPageFirstPageByTopic(keyWords, topic, limit, page)
		if err != nil {
			zap.L().Error("SelectArticleAndUserListByPageFirstPageService() service.article.QueryArticleAndUserListByPageFirstPageByTopic err=", zap.Error(err))
			return nil, err
		}
	}

	if SortWay == "hot" {
		sort.Sort(articles)
	}

	// 遍历文章集合并判断当前用户是否点赞或收藏该文章
	for i := 0; i < len(articles); i++ {
		okSelect, err := redis.IsUserCollected(strconv.Itoa(int(articles[i].ID)), username)
		okLike, err := redis.IsUserLiked(strconv.Itoa(int(articles[i].ID)), username, 0)
		if err != nil {
			zap.L().Error("SelectArticleAndUserListByPageFirstPageService() service.article.IsUserLiked err=", zap.Error(err))
			return nil, err
		}
		articles[i].IsCollect = okSelect
		articles[i].IsLike = okLike

		// 计算发布时间
		articles[i].PostTime = timeConverter.IntervalConversion(articles[i].CreatedAt)
	}

	return articles, nil
}

// PublishArticleService 发布文章
func PublishArticleService(username, content, topic string, wordCount int, tags []string, pics []*multipart.FileHeader, video []*multipart.FileHeader, status bool) error {
	// 检查文本内容字数
	if wordCount < constant.WordLimitMin || wordCount > constant.WordLimitMax {
		zap.L().Error("PublishArticleService() service.article.ArticleService err=", zap.Error(myErr.DataFormatError()))
		return myErr.DataFormatError()
	}

	// 检查标签
	if len(tags) <= 0 {
		zap.L().Error("PublishArticleService() service.article.ArticleService err=", zap.Error(myErr.DataFormatError()))
		return myErr.DataFormatError()
	}

	//  检查视频数量
	if len(video) > 1 {
		zap.L().Error("PublishArticleService() service.article.ArticleService err=", zap.Error(myErr.DataFormatError()))
		return myErr.DataFormatError()
	}

	//  将图片上传至oss
	var picPath []string
	if len(pics) > 0 && len(pics) < 10 {
		for _, pic := range pics {
			url, err := fileProcess.UploadFile("image", pic)
			fmt.Println(url)
			if err != nil {
				zap.L().Error("PublishArticleService() service.article.UploadFile err=", zap.Error(err))
				return err
			}
			picPath = append(picPath, url)
		}
	}

	// 将视频上传至oss
	var videoPath string
	if len(video) > 0 {
		url, err := fileProcess.UploadFile("video", video[0])
		if err != nil {
			zap.L().Error("PublishArticleService() service.article.UploadFile err=", zap.Error(err))
			return err
		}
		videoPath = url
	}

	// 检查本日发表相应话题的文章数
	startOfDay := time.Now().Truncate(24 * time.Hour) // 今天的开始时间
	endOfDay := startOfDay.Add(24 * time.Hour)        // 明天的开始时间

	uid, err := mysql.QueryUserIdByUsername(username)
	if err != nil {
		zap.L().Error("PublishArticleService() service.article.QueryUserIdByUsername err=", zap.Error(err))
		return err
	}

	count, err := mysql.QueryArticleNumByDay(topic, startOfDay, endOfDay, uid)
	if err != nil {
		zap.L().Error("PublishArticleService() service.article.QueryArticleNumByDay err=", zap.Error(err))
		return err
	}
	if count >= constant.ArticlePublishLimit {
		return errors.New("文章发布次数已达上限")
	}

	// 计算文章分数
	point := len(pics)*constant.ImagePointConstant + len(video)*constant.VideoPointConstant + constant.TextPointConstant

	var aid int
	err = mysql.DB.Transaction(func(tx *gorm.DB) error {
		// 插入新文章
		aid, err = mysql.InsertArticleContent(content, topic, uid, wordCount, tags, picPath, videoPath, status, tx)
		if err != nil {
			zap.L().Error("PublishArticleService() service.article.InsertArticleContent err=", zap.Error(err))
			return err
		}

		topicId, err := mysql.QueryTopicIdByTopicName(topic)
		if err != nil {
			zap.L().Error("PublishArticleService() service.article.QueryTagIdByTagName err=", zap.Error(err))
			return err
		}

		if status {
			// 添加分数记录
			err = UpdatePointService(uid, point, topicId, tx)
			if err != nil {
				zap.L().Error("PublishArticleService() service.article.UpdatePointService err=", zap.Error(err))
				return err
			}
		}

		// 将文章更新到redis点赞、收藏
		redis.RDB.HSet("article", strconv.Itoa(aid), 0)
		redis.RDB.HSet("collect", strconv.Itoa(aid), 0)
		return nil
	})
	if err != nil {
		zap.L().Error("PublishArticleService() service.article.Transaction err=", zap.Error(myErr.DataFormatError()))
		return err
	}

	if status {
		// 增加改文章的分数
		err = mysql.UpdateArticlePoint(aid, point)
		if err != nil {
			zap.L().Error("PublishArticleService() service.article.QueryTagIdByTagName err=", zap.Error(err))
			return err
		}
	}

	return nil
}

// GetArticlesByClassService 班级分类查询文章
func GetArticlesByClassService(keyWords, username, sortWay string, limit, page int, class string) ([]model.Article, error) {

	articles, err := mysql.QueryArticleByClass(limit, page, class, keyWords)
	if err != nil {
		zap.L().Error("GetArticlesByClassService() service.article.QueryArticleByClass err=", zap.Error(err))
		return nil, err
	}

	if sortWay == "hot" {
		sort.Sort(articles)
	}

	// 遍历文章集合并判断当前用户是否点赞或收藏该文章
	for i := 0; i < len(articles); i++ {
		okSelect, err := redis.IsUserCollected(strconv.Itoa(int(articles[i].ID)), username)
		okLike, err := redis.IsUserLiked(strconv.Itoa(int(articles[i].ID)), username, 0)
		if err != nil {
			fmt.Println("SelectArticleAndUserListByPageFirstPageService() service.article.IsUserSelectedService err=", err)
			return nil, err
		}
		articles[i].IsCollect = okSelect
		articles[i].IsLike = okLike

		// 计算发布时间
		articles[i].PostTime = timeConverter.IntervalConversion(articles[i].CreatedAt)
	}
	return articles, nil
}

// ReviseArticleStatusService 修改文章的私密状态
func ReviseArticleStatusService(aid int, status bool) error {
	// 检查私密状态
	curStatus, err := mysql.QueryArticleStatusById(aid)
	if err != nil {
		zap.L().Error("ReviseArticleStatus() service.article.QueryArticleStatusById", zap.Error(err))
		return err
	}
	if curStatus == status {
		zap.L().Error("ReviseArticleStatus() service.article", zap.Error(myErr.HasExistError()))
		return myErr.HasExistError()
	}

	// 若修改状态为公开，则检查本日的帖子限额
	if status {
		// 检查本日发表相应话题的文章数
		startOfDay := time.Now().Truncate(24 * time.Hour) // 今天的开始时间
		endOfDay := startOfDay.Add(24 * time.Hour)        // 明天的开始时间

		article, err := mysql.QueryArticleByIdOfManager(aid)
		if err != nil {
			zap.L().Error("PublishArticleService() service.article.QueryArticleByIdOfManager err=", zap.Error(err))
			return err
		}

		count, err := mysql.QueryArticleNumByDay(article.Topic, startOfDay, endOfDay, int(article.UserID))
		if err != nil {
			zap.L().Error("PublishArticleService() service.article.QueryArticleNumByDay err=", zap.Error(err))
			return err
		}
		if count >= constant.ArticlePublishLimit {
			return errors.New("文章发布次数已达上限")
		}
	}

	err = mysql.DB.Transaction(func(tx *gorm.DB) error {
		// 修改状态
		err = mysql.UpdateArticleStatusById(aid, status, tx)
		if err != nil {
			zap.L().Error("ReviseArticleStatus() service.article.UpdateArticleStatusById", zap.Error(err))
			return err
		}

		/*
			回滚积分
		*/

		// 条件：1. 私密&未被封禁	2. 公开&未被封禁

		article, err := mysql.QueryArticleByIdOfManager(aid)
		if err != nil {
			return err
		}

		if !article.Ban {

			curPoint, err := mysql.QueryArticlePoint(aid)
			if err != nil {
				zap.L().Error("ReviseArticleStatus() service.article.QueryArticlePoint", zap.Error(err))
				return err
			}
			if status == false {
				curPoint = -curPoint
			}
			user, err := mysql.QueryUserByArticleId(aid)
			if err != nil {
				zap.L().Error("DeleteArticleService() service.article.QueryUserByArticleId err=", zap.Error(myErr.DataFormatError()))
				return err
			}
			err = UpdatePointByUsernamePointAid(user.Username, curPoint, aid, tx)
			if err != nil {
				zap.L().Error("DeleteArticleService() service.article.UpdatePointByUsernamePointAid err=", zap.Error(myErr.DataFormatError()))
				return err
			}

		}

		return nil
	})
	if err != nil {
		zap.L().Error("DeleteArticleService() service.article.Transaction err=", zap.Error(myErr.DataFormatError()))
		return err
	}

	return nil
}

// AddTopicsService 添加话题
func AddTopicsService(j *jsonvalue.V) error {
	// 获取话题
	name, err := j.GetString("topic_name")
	if err != nil {
		zap.L().Error("AddTopicsService() service.article.ArticleService.GetString err=", zap.Error(err))
		return err
	}
	content, err := j.GetString("topic_name")
	if err != nil {
		zap.L().Error("AddTopicsService() service.article.ArticleService.GetString err=", zap.Error(err))
		return err
	}
	//添加话题
	err = mysql.CreateTopic(name, content)
	if err != nil {
		zap.L().Error("AddTopicsService() service.article.CreateTopic err=", zap.Error(err))
		return err
	}
	return nil
}

// GetAllTopicsService 获取所有话题
func GetAllTopicsService() ([]model.Topic, error) {
	// 获取所有话题
	topics, err := mysql.QueryAllTopics()
	if err != nil {
		zap.L().Error("GetAllTopicsService() service.article.QueryAllTopics err=", zap.Error(err))
		return nil, err
	}
	return topics, nil
}

// AddTagsByTopicService 添加话题标签
func AddTagsByTopicService(topic string, tags []string) error {
	//添加标签
	for _, v := range tags {
		err := mysql.CreateTagByTopic(topic, v)
		if err != nil {
			zap.L().Error("AddTagsByTopicService() service.article.CreateTagByTopic err=", zap.Error(err))
			return err
		}
	}
	return nil
}

// GetTagsByTopicService 获取话题对应的标签
func GetTagsByTopicService(topicId int) ([]map[string]any, error) {
	//获取想要添加标签的对应话题
	tags, err := mysql.QueryTagsByTopic(topicId)
	if err != nil {
		zap.L().Error("GetTagsByTopicService() service.article.QueryTagsByTopic err=", zap.Error(err))
		return nil, err
	}

	list := make([]map[string]any, 0)
	for i, v := range tags {
		list = append(list, map[string]any{
			"id":   i,
			"name": v.TagName,
		})
	}
	return list, nil
}

// AdvancedArticleFilteringService 文章高级筛选
func AdvancedArticleFilteringService(page, limit int, sortField, order, startAt, endAt, topic, keyWords, name, username string, class []string, grade int, role string) (articles []model.Article, err error) {
	if topic == "全部话题" {
		topic = ""
	}
	if role == "teacher" {
		articles, err = mysql.QueryTeacherAndArticleByAdvancedFilter(startAt, endAt, topic, keyWords, sortField, order, name, false, page, limit)
	} else {
		articles, err = mysql.QueryStuAndArticleByAdvancedFilter(startAt, endAt, topic, keyWords, sortField, order, name, grade, class, false, page, limit)
	}
	if err != nil {
		zap.L().Error("AdvancedArticleFilteringService() service.article.QueryUserAndArticleByAdvancedFilter err=", zap.Error(err))
		return nil, err
	}
	// 遍历文章集合并判断当前用户是否点赞或收藏该文章
	for i := 0; i < len(articles); i++ {
		okSelect, err := redis.IsUserCollected(strconv.Itoa(int(articles[i].ID)), username)
		okLike, err := redis.IsUserLiked(strconv.Itoa(int(articles[i].ID)), username, 0)
		if err != nil {
			zap.L().Error("AdvancedArticleFilteringService() service.article.IsUserLiked err=", zap.Error(err))
			return nil, err
		}
		articles[i].IsCollect = okSelect
		articles[i].IsLike = okLike

		// 计算发布时间
		articles[i].PostTime = timeConverter.IntervalConversion(articles[i].CreatedAt)
	}

	return articles, nil
}

// SelectGoodArticleService 评选优秀帖子
func SelectGoodArticleService(articleId, quality int, role, username string) (err error) {
	// 身份验证
	switch role {
	case "class":
		// 查询用户所属班级
		class, err := mysql.QueryClassByUsername(username)
		if err != nil {
			zap.L().Error("SelectGoodArticleService() service.article.QueryTagsByTopic err=", zap.Error(err))
			return err
		}
		if quality > 1 {
			// 权限越界
			return myErr.OverstepCompetence
		}

		err = mysql.UpdateArticleQualityForClass(class, articleId, quality)
		if err != nil {
			zap.L().Error("SelectGoodArticleService() service.article.UpdateArticleQualityForClass err=", zap.Error(err))
			return err
		}

	case "grade1":
		if quality > 2 {
			return myErr.OverstepCompetence
		}
		err = mysql.UpdateArticleQualityForGrade(1, articleId, quality)
	case "grade2":
		if quality > 2 {
			return myErr.OverstepCompetence
		}
		err = mysql.UpdateArticleQualityForGrade(2, articleId, quality)
	case "grade3":
		if quality > 2 {
			return myErr.OverstepCompetence
		}
		err = mysql.UpdateArticleQualityForGrade(3, articleId, quality)
	case "grade4":
		if quality > 2 {
			return myErr.OverstepCompetence
		}
		err = mysql.UpdateArticleQualityForGrade(4, articleId, quality)
	case "college":
		err = mysql.UpdateArticleQualityForSuperMan(articleId, quality)
	case "superman":
		err = mysql.UpdateArticleQualityForSuperMan(articleId, quality)
	}

	if err != nil {
		zap.L().Error("SelectGoodArticleService() service.article.UpdateArticleQualityForSuperMan err=", zap.Error(err))
		return err
	}

	return nil
}

// QueryGoodArticlesByRoleService 根据管理员身份查询对应的优秀帖子
func QueryGoodArticlesByRoleService(role, startAt, endAt, topic, keyWords, sortField, order, name, username string, page, limit int) (articles []model.Article, count int, err error) {
	grade := -1
	switch role {
	case "class":
		// 获取班级
		var class string
		class, err = mysql.QueryClassByUsername(username)
		articles, err = mysql.QueryArticlesByClass(page, limit, startAt, endAt, topic, keyWords, sortField, order, name, class)
		count, err = mysql.QueryClassCommonArticleNum(startAt, endAt, topic, keyWords, sortField, order, name, class)
	case "grade1":
		grade = 1
		articles, err = mysql.QueryClassGoodArticles(1, startAt, endAt, topic, keyWords, sortField, order, name, page, limit)
		count, err = mysql.QueryClassGoodArticleNum(grade, startAt, endAt, topic, keyWords, sortField, order, name)
	case "grade2":
		grade = 2
		articles, err = mysql.QueryClassGoodArticles(2, startAt, endAt, topic, keyWords, sortField, order, name, page, limit)
		count, err = mysql.QueryClassGoodArticleNum(grade, startAt, endAt, topic, keyWords, sortField, order, name)
	case "grade3":
		grade = 3
		articles, err = mysql.QueryClassGoodArticles(3, startAt, endAt, topic, keyWords, sortField, order, name, page, limit)
		count, err = mysql.QueryClassGoodArticleNum(grade, startAt, endAt, topic, keyWords, sortField, order, name)
	case "grade4":
		grade = 4
		articles, err = mysql.QueryClassGoodArticles(4, startAt, endAt, topic, keyWords, sortField, order, name, page, limit)
		count, err = mysql.QueryClassGoodArticleNum(grade, startAt, endAt, topic, keyWords, sortField, order, name)
	case "college":
		articles, err = mysql.QueryGradeGoodArticles(page, limit, startAt, endAt, topic, keyWords, sortField, order, name)
		count, err = mysql.QueryGradeGoodArticleNum(startAt, endAt, topic, keyWords, sortField, order, name)
	case "superman":
		articles, err = mysql.QueryGradeGoodArticles(page, limit, startAt, endAt, topic, keyWords, sortField, order, name)
		count, err = mysql.QueryGradeGoodArticleNum(startAt, endAt, topic, keyWords, sortField, order, name)
	}
	if err != nil {
		zap.L().Error("QueryGoodArticlesByRoleService() service.article.QueryGradeGoodArticles err=", zap.Error(err))
		return nil, -1, err
	}

	return articles, count, nil
}
