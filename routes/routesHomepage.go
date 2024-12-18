package routes

import (
	"github.com/gin-gonic/gin"
	"studentGrow/controller/homepage"
	"studentGrow/utils/token"
)

// 前台个人主页
func routesHomepage(r *gin.Engine) {
	rh := r.Group("/user")

	rh.Use(token.AuthMiddleware())

	{
		rh.POST("/profiles_get", homepage.GetMesControl)
		rh.POST("/userHeadshot_update", homepage.UpdateHeadshotControl)
		rh.POST("/selfCotnent_get", homepage.GetSelfContentContro)
		rh.POST("/selfContent_update", homepage.UpdateSelfContentContro)
		rh.POST("/userMotto_update", homepage.UpdateHomepageMottoControl)
		rh.POST("/userPhone_update", homepage.UpdatePhoneNumberControl)
		rh.POST("/userEmail_update", homepage.UpdateEmailControl)
		rh.POST("/userData_get", homepage.GetUserDataControl)
		rh.POST("/fans_get", homepage.GetFansListControl)
		rh.POST("/concern_get", homepage.GetConcernListControl)
		rh.POST("/isConcern_get", homepage.GetIsConcernControl)
		rh.POST("/concern_change", homepage.ChangeConcernControl)
		rh.POST("/history_get", homepage.GetHistoryControl)
		// 我的足迹
		rh.POST("/tracks_get", homepage.GetTracksControl)
		rh.POST("/star_get", homepage.GetStarControl)
		rh.POST("/class_get", homepage.GetClassControl)
		rh.GET("/points_get", homepage.GetTopicPointsControl)
		rh.GET("/article_get", homepage.GetArticleControl)
		rh.POST("/ban", homepage.BanUserControl)
		rh.POST("/unban", homepage.UnbanUserControl)
		rh.POST("/pwd_update", homepage.ChangePasswordControl)
		rh.POST("/advice_get", homepage.GetAdviceControl)
	}
}
