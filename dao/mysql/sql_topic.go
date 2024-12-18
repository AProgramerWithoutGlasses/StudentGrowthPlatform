package mysql

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"studentGrow/models/gorm_model"
	myErr "studentGrow/pkg/error"
)

// QueryIsExistByTopicName 查询话题是否存在通过话题名字
func QueryIsExistByTopicName(topicName string) (bool, error) {
	var count int64
	if err := DB.Model(&gorm_model.Topic{}).Where("topic_name = ?", topicName).Count(&count).Error; err != nil {
		zap.L().Error("QueryIsExistByTopicName() dao.mysql.sql_topic err=", zap.Error(err))
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

// CreateTopic 添加话题
func CreateTopic(topicName string, topicContent string) error {
	topic := gorm_model.Topic{
		TopicName:    topicName,
		TopicContent: topicContent,
	}
	// 检查是否存在该话题
	ok, err := QueryIsExistByTopicName(topicName)
	if err != nil {
		zap.L().Error("CreateTopic() dao.mysql.sql_topic.QueryIsExistByTopicName err=", zap.Error(err))
		return err
	}
	// 若存在
	if ok {
		return myErr.HasExistError()
	} else {
		if err = DB.Create(&topic).Error; err != nil {
			zap.L().Error("CreateTopic() dao.mysql.sql_topic.Create err=", zap.Error(myErr.HasExistError()))
			return err
		}
	}
	return nil
}

// CreateTagByTopic 添加话题所对应的标签
func CreateTagByTopic(topicName string, tagName string) error {
	topicId, err := QueryTopicIdByTopicName(topicName)
	if err != nil {
		zap.L().Error("CreateTopic() dao.mysql.sql_topic.QueryTopicIdByTopicName err=", zap.Error(err))
		return err
	}
	tag := gorm_model.Tag{
		TopicID: uint(topicId),
		TagName: tagName,
	}
	// 检查是否存在该话题
	ok, err := QueryIsExistByTopicName(topicName)
	if err != nil {
		zap.L().Error("CreateTopic() dao.mysql.sql_topic.QueryIsExistByTopicName err=", zap.Error(err))
		return err
	}
	if ok {
		if err = DB.Create(&tag).Error; err != nil {
			zap.L().Error("CreateTopic() dao.mysql.sql_topic.Create err=", zap.Error(err))
			return err
		}
	}
	return nil
}

// DeleteTagById 删除标签
func DeleteTagById(id int) error {
	if err := DB.Where("id = ?", id).Delete(&gorm_model.Tag{}).Error; err != nil {
		zap.L().Error("DeleteTagById() dao.mysql.sql_topic.Delete err=", zap.Error(err))
		return err
	}
	return nil
}

// QueryAllTopics 获取所有话题
func QueryAllTopics() ([]gorm_model.Topic, error) {
	var topics []gorm_model.Topic
	if err := DB.Find(&topics).Error; err != nil {
		zap.L().Error("QueryAllTopics() dao.mysql.sql_topic err=", zap.Error(err))
		return nil, err
	}
	return topics, nil
}

// QueryTagsByTopic 获取话题对应的标签
func QueryTagsByTopic(id int) ([]gorm_model.Tag, error) {
	var tags []gorm_model.Tag
	if err := DB.Where("topic_id = ?", id).Find(&tags).Error; err != nil {
		zap.L().Error("QueryAllTopics() dao.mysql.sql_topic err=", zap.Error(err))
		return nil, err
	}

	return tags, nil
}

// QueryTagIdByTagName 根据标签名字查找标签ID
func QueryTagIdByTagName(name string) (int, error) {
	var tag gorm_model.Tag
	if err := DB.Where("tag_name = ?", name).First(&tag).Error; err != nil {
		zap.L().Error("QueryTagIdByTagName() dao.mysql.sql_topic err=", zap.Error(myErr.ErrNotFoundError))
		return -1, myErr.ErrNotFoundError
	}
	return int(tag.ID), nil
}

// InsertArticleTags 添加文章标签
func InsertArticleTags(tags []string, articleId int, db *gorm.DB) error {
	for _, tag := range tags {
		tagId, err := QueryTagIdByTagName(tag)
		if err != nil {
			zap.L().Error("InsertArticleTags() dao.mysql.sql_topic.QueryTagIdByTagName err=", zap.Error(err))
			return err
		}
		articleTag := gorm_model.ArticleTag{
			ArticleID: uint(articleId),
			TagID:     uint(tagId),
		}
		if err = db.Create(&articleTag).Error; err != nil {
			zap.L().Error("InsertArticleTags() dao.mysql.sql_topic.Create err=", zap.Error(err))
			return err
		}
	}
	return nil
}

// QueryTopicByArticleId 通过文章id查询文章的话题
func QueryTopicByArticleId(aid int) (string, error) {
	var article gorm_model.Article
	if err := DB.Where("id = ?", aid).First(&article).Error; err != nil {
		zap.L().Error("QueryTopicByArticleId() dao.mysql.sql_topic.First err=", zap.Error(err))
		return "", err
	}
	return article.Topic, nil
}

// QueryTopicIdByTopicName 通过topicName查询topicId
func QueryTopicIdByTopicName(topicName string) (int, error) {
	var topic gorm_model.Topic
	if err := DB.Where("topic_name = ?", topicName).First(&topic).Error; err != nil {
		zap.L().Error("QueryTopicByArticleId() dao.mysql.sql_topic.First err=", zap.Error(err))
		return -1, err
	}
	return int(topic.ID), nil
}
