package article

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	myErr "studentGrow/pkg/error"
	res "studentGrow/pkg/response"
	"studentGrow/service/article"
	"studentGrow/utils/token"
)

// LikeController 点赞
func LikeController(c *gin.Context) {

	in := struct {
		Id          int    `json:"id"`
		LikeType    int    `json:"like_type"`
		TarUsername string `json:"tar_username"`
	}{}

	err := c.ShouldBindJSON(&in)
	if err != nil {
		return
	}

	// 通过token获取username
	aToken := token.NewToken(c)
	user, exist := aToken.GetUser()
	if !exist {
		res.ResponseError(c, res.TokenError)
		zap.L().Error("token错误")
		return
	}

	username := user.Username

	//点赞
	err = article.LikeObjOrNot(strconv.Itoa(in.Id), username, in.TarUsername, in.LikeType)
	if err != nil {
		zap.L().Error("Like() controller.article.LikeObjOrNot err=", zap.Error(err))
		myErr.CheckErrors(err, c)
		return
	}

	//返回数据
	res.ResponseSuccess(c, struct{}{})
}

// CancelLikeController 取消点赞
//func CancelLikeController(c *gin.Context) {
//	//解析数据
//	j, e := utils.GetJsonvalue(c)
//	if e != nil {
//		fmt.Println("Like() controller.article.GetJsonvalue err=", e)
//		return
//	}
//
//	objId, userId, likeType, err := GetParam(j, "objId", "userId", "likeType")
//	if err != nil {
//		return
//	}
//
//	//取消点赞
//	err = article.CancelLike(objId, userId, likeType)
//	if err != nil {
//		fmt.Println("Like() controller.article.AnalyzeToMap err=", err)
//		return
//	}
//
//	//返回数据
//	res.ResponseSuccess(c, nil)
//}

// CheckLikeOrNotController 检查是否点赞
//func CheckLikeOrNotController(c *gin.Context) {
//	//解析数据
//	j, e := utils.GetJsonvalue(c)
//	if e != nil {
//		fmt.Println("Like() controller.article.GetJsonvalue err=", e)
//		return
//	}
//
//	objId, userId, likeType, err := GetParam(j, "objId", "userId", "likeType")
//	if err != nil {
//		fmt.Println("Like() controller.article.GetParam err=", err)
//		return
//	}
//
//	//检查用户是否已经点赞
//	result, err := redis.IsUserLiked(objId, userId, likeType)
//	if err != nil {
//		fmt.Println("Like() controller.article.IsUserLiked err=", err)
//		return
//	}
//
//	//返回数据
//	if result {
//		res.ResponseSuccess(c, "已点赞")
//	} else {
//		res.ResponseSuccess(c, "未点赞")
//	}
//}
//
//// GetObjLikeNumController 获取当前文章或评论的点赞数量
//func GetObjLikeNumController(c *gin.Context) {
//	//解析数据
//	j, e := utils.GetJsonvalue(c)
//	if e != nil {
//		fmt.Println("Like() controller.article.GetJsonvalue err=", e)
//		return
//	}
//	objId, _ := j.GetString("objId")
//	likeType, _ := j.GetInt("likeType")
//
//	//获取点赞数
//	likeSum, err := redis.GetObjLikes(objId, likeType)
//	if err != nil {
//		fmt.Println("Like() controller.article.GetObjLikes err=", err)
//		return
//	}
//	//返回数据
//	res.ResponseSuccess(c, likeSum)
//}

// 获取文章或评论的点赞集合
//func GetObjLikedUsersController() {
//
//}
