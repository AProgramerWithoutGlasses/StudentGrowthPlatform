package article

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"studentGrow/dao/mysql"
)

func UpdatePointService(uid, point, topicId int, db *gorm.DB) error {
	// 查询point表中是否存在该用户对应的话题分数
	ok, err := mysql.QueryUserPointOfTopicIsExist(uid, topicId)
	if err != nil {
		zap.L().Error("UpdatePointService() service.article.QueryUserPointOfTopicIsExist err=", zap.Error(err))
		return err
	}

	// 若不存在，则创建该用户的话题分数
	if !ok {
		err = mysql.CreateUserPointOfTopic(uid, topicId)
		if err != nil {
			zap.L().Error("UpdatePointService() service.article.CreateUserPointOfTopic err=", zap.Error(err))
			return err
		}
	}

	// 获取当前分数
	curPoint, err := mysql.QueryUserPointByTopic(topicId, uid)
	if err != nil {
		zap.L().Error("UpdatePointService() service.article.QueryUserPointByTopic err=", zap.Error(err))
		return err
	}

	// 修改后的分数
	aftarPoint := curPoint + point
	if aftarPoint < 0 {
		aftarPoint = 0
	}

	if curPoint >= 0 {
		// 修改分数
		err = mysql.UpdateUserPointByTopic(aftarPoint, uid, topicId, db)
		if err != nil {
			zap.L().Error("UpdatePointService() service.article.UpdateUserPointByTopic err=", zap.Error(err))
			return err
		}
	}
	return nil
}

// UpdatePointByUsernamePointAid 已知username,point,aid修改分数
func UpdatePointByUsernamePointAid(username string, point, aid int, db *gorm.DB) error {
	uid, err := mysql.GetIdByUsername(username)
	if err != nil {
		zap.L().Error("UpdatePointByUsernamePointAid() service.article.GetIdByUsername err=", zap.Error(err))
		return err
	}
	//  获取topicId
	topic, err := mysql.QueryTopicByArticleId(aid)
	if err != nil {
		zap.L().Error("UpdatePointByUsernamePointAid() service.article.QueryTopicByArticleId err=", zap.Error(err))
		return err
	}
	topicId, err := mysql.QueryTopicIdByTopicName(topic)
	if err != nil {
		zap.L().Error("UpdatePointByUsernamePointAid() service.article.QueryTopicIdByTopicName err=", zap.Error(err))
		return err
	}

	err = UpdatePointService(uid, point, topicId, db)
	if err != nil {
		zap.L().Error("UpdatePointByUsernamePointAid() service.article.UpdatePointService err=", zap.Error(err))
		return err
	}
	return nil
}
