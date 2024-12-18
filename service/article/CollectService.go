package article

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"studentGrow/dao/mysql"
	"studentGrow/dao/redis"
	"studentGrow/models/nzx_model"
	"studentGrow/pkg/sse"
	NotificationPush "studentGrow/service/notificationPush"
)

// CollectService Collect 收藏
func CollectService(username, aid string) error {
	// 收藏文章
	err := redis.AddArticleToCollectSet(username, aid)
	if err != nil {
		zap.L().Error("CollectService() service.article.AddArticleToCollectSet err=", zap.Error(err))
		return err
	}
	// 获取文章收藏数
	selections, err := redis.GetArticleCollections(aid)
	if err != nil {
		zap.L().Error("CollectService() service.article.SetArticleCollections err=", zap.Error(err))
		return err
	}
	// 收藏数+1
	if selections >= 0 {
		err = redis.SetArticleCollections(aid, 1)
		if err != nil {
			zap.L().Error("CollectService() service.article.SetArticleCollections err=", zap.Error(err))
			return err
		}
	}
	if err != nil {
		zap.L().Error("CollectService() service.article.TxPipelined err=", zap.Error(err))
		return err
	}

	// 写入通道
	articleId, err := strconv.Atoi(aid)
	if err != nil {
		zap.L().Error("CollectService() service.article.Atoi err=", zap.Error(err))
		return err
	}
	ArticleCollectChan <- nzx_model.RedisCollectData{Aid: articleId, Username: username, Operator: "collect"}
	return nil
}

// CancelCollectService CancelCollect 取消收藏
func CancelCollectService(aid, username string) error {
	isExist, err := redis.IsUserCollected(username, aid)
	if err != nil {
		zap.L().Error("CancelCollectService() service.article.IsUserCollected err=", zap.Error(err))
		return err
	}

	if isExist {
		err = redis.RemoveUserCollectionSet(aid, username)
		if err != nil {
			zap.L().Error("CancelCollectService() service.article.RemoveUserCollectionSet err=", zap.Error(err))
			return err
		}
		selections, err := redis.GetArticleCollections(aid)
		if err != nil {
			zap.L().Error("CancelCollectService() service.article.GetArticleCollections err=", zap.Error(err))
			return err
		}
		if selections > 0 {
			err := redis.SetArticleCollections(aid, -1)
			if err != nil {
				zap.L().Error("CancelCollectService() service.SetArticleCollections.Atoi err=", zap.Error(err))
				return err
			}
		}
		if err != nil {
			zap.L().Error("CancelCollectService() service.TxPipelined.Atoi err=", zap.Error(err))
			return err
		}

		// 写入通道
		articleId, err := strconv.Atoi(aid)
		if err != nil {
			zap.L().Error("CancelCollectService() service.SetArticleCollections.Atoi err=", zap.Error(err))
			return err
		}
		ArticleCollectChan <- nzx_model.RedisCollectData{Aid: articleId, Username: username, Operator: "cancel_collect"}
	}

	return nil
}

// CollectOrNotService CollectOrNot 检查是否收藏并收藏或取消收藏
func CollectOrNotService(aid, username, tarUsername string) error {
	// 获取当前用户收藏列表
	slice, err := redis.GetUserCollectionSet(username)
	if err != nil {
		fmt.Println("CollectOrNot() service.article.GetUserCollectionSet err=", err)
		return err
	}

	// 若存在该文章,则取消收藏
	selectArticles := make(map[string]struct{})
	for _, s := range slice {
		selectArticles[s] = struct{}{}
	}
	_, ok := selectArticles[aid]

	if len(selectArticles) > 0 && ok {
		err = CancelCollectService(aid, username)
		if err != nil {
			fmt.Println("CollectOrNot() service.article.GetUserCollectionSet err=", err)
			return err
		}
	} else {
		// 反之，收藏
		err = CollectService(username, aid)
		if err != nil {
			fmt.Println("CollectOrNot() service.article.CollectService err=", err)
			return err
		}
		id, err := strconv.Atoi(aid)
		if err != nil {
			fmt.Println("CollectOrNot() service.article.Atoi err=", err)
			return err
		}
		notification, err := NotificationPush.BuildCollectNotification(username, tarUsername, id)
		if err != nil {
			fmt.Println("CollectOrNot() service.article.Atoi err=", err)
			return err
		}
		sse.SendInterNotification(*notification)
	}
	return nil
}

/*
mysql
*/

// CollectToMysql 收藏
func CollectToMysql(aid int, username string) error {
	uid, err := mysql.GetIdByUsername(username)
	if err != nil {
		zap.L().Error("CollectToMysql() service.article.GetIdByUsername err=", zap.Error(err))
		return err
	}
	err = mysql.DB.Transaction(func(tx *gorm.DB) error {
		// 添加收藏记录
		err = mysql.InsertCollectRecord(aid, uid, tx)
		if err != nil {
			zap.L().Error("CollectToMysql() service.article.InsertCollectRecord err=", zap.Error(err))
			return err
		}

		// 获取收藏数
		num, err := mysql.QueryCollectNum(aid)
		if err != nil {
			zap.L().Error("CollectToMysql() service.article.QueryCollectNum err=", zap.Error(err))
			return err
		}
		fmt.Println("mysql读入收藏，收藏后收藏数为:", num)

		// 收藏数+1
		err = mysql.UpdateCollectNum(aid, num+1, tx)
		if err != nil {
			zap.L().Error("CollectToMysql() service.article.UpdateCollectNum err=", zap.Error(err))
			return err
		}
		return nil
	})

	afterNum, err := mysql.QueryCollectNum(aid)
	if err != nil {
		zap.L().Error("CollectToMysql() service.article.QueryCollectNum err=", zap.Error(err))
		return err
	}

	fmt.Println("mysql读入收藏，收藏后收藏数为:", afterNum)
	if err != nil {
		zap.L().Error("CollectToMysql() service.article.Transaction err=", zap.Error(err))
		return err
	}
	return nil
}

// CancelCollectToMysql 取消收藏
func CancelCollectToMysql(aid int, username string) error {
	uid, err := mysql.GetIdByUsername(username)
	if err != nil {
		zap.L().Error("CollectToMysql() service.article.Transaction err=", zap.Error(err))
		return err
	}
	err = mysql.DB.Transaction(func(tx *gorm.DB) error {
		// 删除收藏记录
		err = mysql.DeleteCollectRecord(aid, uid, tx)
		if err != nil {
			zap.L().Error("CollectToMysql() service.article.DeleteCollectRecord err=", zap.Error(err))
			return err
		}

		// 获取收藏数
		num, err := mysql.QueryCollectNum(aid)
		if err != nil {
			zap.L().Error("CollectToMysql() service.article.QueryCollectNum err=", zap.Error(err))
			return err
		}
		fmt.Println("mysql读入取消收藏，收藏前收藏数为:", num)
		// 收藏数-1
		err = mysql.UpdateCollectNum(aid, num-1, tx)
		if err != nil {
			zap.L().Error("CollectToMysql() service.article.UpdateCollectNum err=", zap.Error(err))
			return err
		}
		return nil
	})
	afterNum, err := mysql.QueryCollectNum(aid)
	if err != nil {
		zap.L().Error("CollectToMysql() service.article.QueryCollectNum err=", zap.Error(err))
		return err
	}

	fmt.Println("mysql读入取消收藏，取消收藏后收藏数为:", afterNum)
	if err != nil {
		zap.L().Error("CancelCollectToMysql() service.article.Transaction err=", zap.Error(err))
		return err
	}
	return nil
}
