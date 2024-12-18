package NotificationPush

import (
	"studentGrow/dao/mysql"
	"studentGrow/models/gorm_model"
)

// BuildLikeNotification 构建点赞消息
func BuildLikeNotification(username, tarUsername string, objId, likeType int) (*gorm_model.InterNotification, error) {
	ownUserId, err := mysql.QueryUserIdByUsername(username)
	if err != nil {
		return nil, err
	}
	tarUserId, err := mysql.QueryUserIdByUsername(tarUsername)
	if err != nil {
		return nil, err
	}

	user, err := mysql.QueryUserByUserId(ownUserId)
	if err != nil {
		return nil, err
	}

	notification := gorm_model.InterNotification{
		TarUserId:  uint(tarUserId),
		OwnUserId:  uint(ownUserId),
		NoticeType: 0,
		SuperType:  likeType,
		SuperId:    objId,
		IsRead:     false,
		Content:    "",
		OwnUser:    *user,
	}

	return &notification, nil
}

// BuildCollectNotification 构建收藏消息
func BuildCollectNotification(username, tarUsername string, aid int) (*gorm_model.InterNotification, error) {
	ownUserId, err := mysql.QueryUserIdByUsername(username)
	if err != nil {
		return nil, err
	}
	tarUserId, err := mysql.QueryUserIdByUsername(tarUsername)
	if err != nil {
		return nil, err
	}

	user, err := mysql.QueryUserByUserId(ownUserId)
	if err != nil {
		return nil, err
	}

	notification := gorm_model.InterNotification{
		TarUserId:  uint(tarUserId),
		OwnUserId:  uint(ownUserId),
		NoticeType: 1,
		SuperType:  0,
		SuperId:    aid,
		IsRead:     false,
		Content:    "",
		OwnUser:    *user,
	}

	return &notification, nil
}

// BuildCommentNotification 构建评论消息
func BuildCommentNotification(username, tarUsername string, objId, comType int) (*gorm_model.InterNotification, error) {
	ownUserId, err := mysql.QueryUserIdByUsername(username)
	if err != nil {
		return nil, err
	}
	tarUserId, err := mysql.QueryUserIdByUsername(tarUsername)
	if err != nil {
		return nil, err
	}

	user, err := mysql.QueryUserByUserId(ownUserId)
	if err != nil {
		return nil, err
	}

	notification := gorm_model.InterNotification{
		TarUserId:  uint(tarUserId),
		OwnUserId:  uint(ownUserId),
		NoticeType: 2,
		SuperType:  comType,
		SuperId:    objId,
		IsRead:     false,
		Content:    "",
		OwnUser:    *user,
	}

	return &notification, nil
}
