package routesJoinAudit

import (
	"github.com/gin-gonic/gin"
	"studentGrow/dao/mysql"
	res "studentGrow/pkg/response"
	"time"
)

type ResOpenActivityMsg struct {
	IsShow       string `json:"is_show"`
	ID           uint
	ActivityName string    `json:"activity_name"`
	StartTime    time.Time `json:"start_time"`
	StopTime     time.Time `json:"stop_time"`
}

// OpenMsg 判断当前入团申请是否开放
func OpenMsg(c *gin.Context) {
	isShow, Msg, ActivityMsg := mysql.OpenActivityStates()
	Response := ResOpenActivityMsg{
		IsShow:       isShow,
		ID:           ActivityMsg.ID,
		ActivityName: ActivityMsg.ActivityName,
		StartTime:    ActivityMsg.StartTime,
		StopTime:     ActivityMsg.StopTime,
	}
	res.ResponseSuccessWithMsg(c, Msg, Response)
	return
}
