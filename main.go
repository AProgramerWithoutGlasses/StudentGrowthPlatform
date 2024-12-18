package main

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"studentGrow/aliyun/oss"
	"studentGrow/dao/mysql"
	"studentGrow/dao/redis"
	"studentGrow/logger"
	"studentGrow/routes"
	"studentGrow/service/article"
	"studentGrow/settings"
	"syscall"
	"time"
)

func main() {
	// 1. 加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("settings.Init() viper.ReadInConfig() err : %v\n", err)
		return
	}

	// 2. 初始化日志
	if err := logger.Init(); err != nil {
		fmt.Printf("logger.Init() l.UnmarshalText() err : %v\n", err)
		return
	}
	defer zap.L().Sync()
	zap.Open()

	// 3. 初始化Mysql
	if err := mysql.Init(); err != nil {
		fmt.Printf("mysql.Init() gorm.Open() err : %v\n", err)
		return
	}
	//err := mysql.DB.AutoMigrate(
	//	&gorm_model.JoinAudit{},
	//	&gorm_model.JoinAuditDuty{},
	//	&gorm_model.User{},
	//	&gorm_model.JoinAuditFile{},
	//)
	//if err != nil {
	//	return
	//}

	// 4. 初始化redis
	if err := redis.Init(); err != nil {
		fmt.Printf("redis.Init() rdb.Ping().Result() err : %v\n", err)
		return
	}

	// redis读写mysql
	article.InitMyMQ()

	// 5. 注册路由
	r := routes.Setup()

	//6.初始化oss
	err := ossProject.Init()
	if err != nil {
		zap.L().Error("main() oss.Init err=", zap.Error(err))
		return
	}

	//// 初始化eventBus
	//eventBus.InitEventBus()

	// 7. 启动服务（优雅关机）
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shut down Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Shut down Server: ", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}
