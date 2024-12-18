package redis

import (
	"fmt"
	"go.uber.org/zap"
	"strconv"
)

var Selection = "collect"

// AddArticleToCollectSet 添加文章到用户收藏集合中
func AddArticleToCollectSet(uid string, aid string) error {
	if err := RDB.SAdd(Selection+uid, aid).Err(); err != nil {
		zap.L().Error("AddArticleToCollectSet() dao.redis.redis_collect.SAdd err=", zap.Error(err))
		return err
	}
	return nil
}

// IsUserCollected 检查用户是否已经收藏
func IsUserCollected(uid string, aid string) (bool, error) {
	result, err := RDB.SIsMember(Selection+uid, aid).Result()
	if err != nil {
		zap.L().Error("IsUserCollected() dao.redis.redis_collect.SIsMember err=", zap.Error(err))
		return false, nil
	}

	return result, nil
}

// SetArticleCollections 设置文章收藏数
func SetArticleCollections(aid string, selectNum int) error {
	if err := RDB.HIncrBy(Selection, aid, int64(selectNum)).Err(); err != nil {
		zap.L().Error("SetArticleCollections() dao.redis.redis_collect.HIncrBy err=", zap.Error(err))
		return err
	}
	return nil
}

// GetArticleCollections 获取文章收藏数
func GetArticleCollections(aid string) (int, error) {
	selectNum, err := RDB.HGet(Selection, aid).Result()
	fmt.Println(Selection, aid)
	if err != nil {
		zap.L().Error("GetArticleCollections() dao.redis.redis_collect.HGet err=", zap.Error(err))
		return -1, err
	}
	res, err := strconv.Atoi(selectNum)
	if err != nil {
		zap.L().Error("GetArticleCollections() dao.redis.redis_collect.Atoi err=", zap.Error(err))
		return -1, err
	}
	return res, nil
}

// GetUserCollectionSet 获取用户的收藏集合
func GetUserCollectionSet(uid string) ([]string, error) {
	slice, err := RDB.SMembers(Selection + uid).Result()
	if err != nil {
		zap.L().Error("GetUserCollectionSet() dao.redis.redis_collect.SMembers err=", zap.Error(err))
		return nil, err
	}

	return slice, nil
}

// RemoveUserCollectionSet 将文章从用户收藏集合中移除
func RemoveUserCollectionSet(aid, uid string) error {
	fmt.Println("uid", uid)
	err := RDB.SRem(Selection+uid, aid).Err()
	if err != nil {
		zap.L().Error("RemoveUserCollectionSet() dao.redis.redis_collect.SRem err=", zap.Error(err))
		return err
	}
	return nil
}
