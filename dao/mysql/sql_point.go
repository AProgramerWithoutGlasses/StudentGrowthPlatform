package mysql

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"studentGrow/models/gorm_model"
	myErr "studentGrow/pkg/error"
)

// QueryUserPointByTopic 查询用户话题积分
func QueryUserPointByTopic(topicId, uid int) (int, error) {
	var point gorm_model.UserPoint
	if err := DB.Where("topic_id = ? and user_id = ?", topicId, uid).First(&point).Error; err != nil {
		zap.L().Error("QueryUserPointByTopic() dao.mysql.sql_point.First err=", zap.Error(err))
		return -1, err
	}
	return point.Point, nil
}

// UpdateUserPointByTopic 修改用户话题积分
func UpdateUserPointByTopic(point, uid, topicId int, db *gorm.DB) error {
	if err := db.Model(&gorm_model.UserPoint{}).Where("topic_id = ? and user_id = ?", topicId, uid).Update("point", point).Error; err != nil {
		zap.L().Error("UpdateUserPointByTopic() dao.mysql.sql_point.Update err=", zap.Error(err))
		return err
	}
	return nil
}

// QueryUserAllPoint 查询用户所有话题积分
func QueryUserAllPoint(uid int) ([]gorm_model.UserPoint, error) {
	var points []gorm_model.UserPoint
	if err := DB.Where("user_id = ?", uid).First(&points).Error; err != nil {
		zap.L().Error("QueryUserAllPoint() dao.mysql.sql_point.First err=", zap.Error(err))
		return nil, err
	}

	if len(points) == 0 {
		zap.L().Error("QueryUserAllPoint() dao.mysql.sql_point.First err=", zap.Error(myErr.ErrNotFoundError))
		return nil, myErr.ErrNotFoundError
	}
	return points, nil
}

// QueryUserPointOfTopicIsExist 查询用户是否存在该话题的积分
func QueryUserPointOfTopicIsExist(uid, topicId int) (bool, error) {
	var count int64
	if err := DB.Model(&gorm_model.UserPoint{}).Where("topic_id = ? and user_id = ?", topicId, uid).Count(&count).Error; err != nil {
		zap.L().Error("QueryUserAllPoint() dao.mysql.sql_point.First err=", zap.Error(err))
		return false, err
	}
	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

// CreateUserPointOfTopic 创建用户话题分数
func CreateUserPointOfTopic(uid, topicId int) error {
	point := gorm_model.UserPoint{
		UserID:  uint(uid),
		TopicID: uint(topicId),
		Point:   0,
	}
	if err := DB.Create(&point).Error; err != nil {
		zap.L().Error("QueryUserAllPoint() dao.mysql.sql_point.First err=", zap.Error(err))
		return err
	}
	return nil
}
