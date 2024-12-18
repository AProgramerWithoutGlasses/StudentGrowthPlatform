package sse

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"studentGrow/dao/mysql"
	"studentGrow/models/gorm_model"
	"sync"
)

var isInitChannelsMap = false
var ChannelsMap sync.Map

func AddChannel(userId int) {
	if !isInitChannelsMap {
		ChannelsMap = sync.Map{}
		isInitChannelsMap = true
	}
	newChannel := make(chan string, 10)
	ChannelsMap.Store(userId, newChannel)
	fmt.Println("Build SSE connection for user = ", userId, newChannel)
}

func BuildNotificationChannel(username string, c *gin.Context) error {
	userId, err := mysql.QueryUserIdByUsername(username)
	if err != nil {
		return err
	}
	AddChannel(userId)
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	w := c.Writer
	flusher, _ := w.(http.Flusher)

	curChan, _ := ChannelsMap.Load(userId)

	// 监听客户端通道是否被关闭
	closeNotify := c.Request.Context().Done()

	go func() {
		<-closeNotify
		close(curChan.(chan string))
		ChannelsMap.Delete(userId)
		fmt.Println("SSE close for user = ", userId)
		return
	}()

	for msg := range curChan.(chan string) {
		_, err := fmt.Fprintf(w, "data:%s\n\n", msg)
		if err != nil {
			return err
		}
		flusher.Flush()
	}
	return nil
}

// SendInterNotification 互动消息推送
func SendInterNotification(n gorm_model.InterNotification) {

	// 若对方用户不在线，则不推送消息
	if val, ok := ChannelsMap.Load(int(n.TarUserId)); !ok {
		fmt.Println("Send interNotification to user = ", n.TarUserId, val)
		return
	}

	notification := map[string]any{
		"notice_type":   n.NoticeType,
		"content":       n.Content,
		"is_read":       n.IsRead,
		"time":          n.Time,
		"name":          n.OwnUser.Name,
		"user_headshot": n.OwnUser.HeadShot,
	}

	msg, err := json.Marshal(notification)
	if err != nil {
		return
	}

	ChannelsMap.Range(func(key, value any) bool {
		k := key.(int)
		fmt.Println("k", k)
		if k == int(n.TarUserId) {
			channel := value.(chan string)
			channel <- string(msg)
		}
		return true
	})
}

// SendSysNotification 广播消息推送
func SendSysNotification(n gorm_model.SysNotification) {
	fmt.Println("Send sysNotification is user = ", n.OwnUserId)
	notification := map[string]any{
		"notice_type":   n.NoticeType,
		"content":       n.Content,
		"is_read":       n.IsRead,
		"time":          n.Time,
		"name":          n.OwnUser.Name,
		"user_headshot": n.OwnUser.HeadShot,
	}

	msg, err := json.Marshal(notification)
	if err != nil {
		return
	}
	ChannelsMap.Range(func(key, value any) bool {
		channel := value.(chan string)
		channel <- string(msg)
		return true
	})
}
