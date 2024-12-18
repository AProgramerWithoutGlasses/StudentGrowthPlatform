package redis

import (
	"fmt"
	"go.uber.org/zap"
	"strconv"
)

var List = []string{"article", "comment"}

// AddUserToLikeSet 添加用户到文章或评论点赞集合中
func AddUserToLikeSet(objId, userId string, likeType int) error {

	err := RDB.SAdd(List[likeType]+objId, userId).Err()

	if err != nil {
		zap.L().Error("AddUserToLikeSet() dao.redis.redis_like.SAdd err=", zap.Error(err))
		return err
	}

	return nil
}

// IsUserLiked 检查用户是否已经点赞
func IsUserLiked(objId, userId string, likeType int) (bool, error) {
	res, err := RDB.SIsMember(List[likeType]+objId, userId).Result()
	if err != nil {
		zap.L().Error("IsUserLiked() dao.redis.redis_like.Result err=", zap.Error(err))
		return false, err
	}
	return res, err
}

// SetObjLikes 设置文章或评论的点赞数
func SetObjLikes(objId string, likeNum int, likeType int) error {
	err := RDB.HIncrBy(List[likeType], objId, int64(likeNum)).Err()
	if err != nil {
		zap.L().Error("SetObjLikes() dao.redis.redis_like.HIncrBy err=", zap.Error(err))
		return err
	}
	return nil
}

// GetObjLikes 获取文章或评论点赞数
func GetObjLikes(objId string, likeType int) (int, error) {
	likesNumResult, err := RDB.HGet(List[likeType], objId).Result()
	fmt.Println(List[likeType], objId)
	fmt.Println(likesNumResult)
	result, _ := RDB.HKeys(List[likeType]).Result()
	fmt.Println(result)

	if err != nil {
		zap.L().Error("GetObjLikes() dao.redis.redis_like.Result err=", zap.Error(err))
		return -1, err
	}
	res, err := strconv.Atoi(likesNumResult)
	if err != nil {
		zap.L().Error("GetObjLikes() dao.redis.redis_like.Atoi err=", zap.Error(err))
		return -1, err
	}
	return res, nil
}

// GetObjLikedUsers 获取文章或评论点赞的用户username集合
func GetObjLikedUsers(objId string, likeType int) (result []string, err error) {
	slice, err := RDB.SMembers(List[likeType] + objId).Result()

	if err != nil {
		zap.L().Error("GetObjLikedUsers() dao.redis.redis_like.SMembers err=", zap.Error(err))
		return nil, err
	}
	return slice, nil
}

// RemoveUserFromLikeSet 移除用户从文章或评论的点赞集合中
func RemoveUserFromLikeSet(objId, userId string, likeType int) error {
	err := RDB.SRem(List[likeType]+objId, userId).Err()
	if err != nil {
		zap.L().Error("RemoveUserFromLikeSet() dao.redis.redis_like.v err=", zap.Error(err))
		return err
	}

	return nil
}
