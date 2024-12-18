package mysql

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func Init() error {
	username := viper.GetString("mysql.user")     //账号
	password := viper.GetString("mysql.password") //密码
	host := viper.GetString("mysql.host")         //数据库地址，可以是Ip或者域名
	port := viper.GetInt("mysql.port")            //数据库端口
	Dbname := viper.GetString("mysql.dbname")     //数据库名updated_at
	timeout := "5s"                               //连接超时，5秒
	//root:root@tcp(127.0.0.1:3306)/gorm?
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", username, password, host, port, Dbname, timeout)
	//连接MYSQL, 获得DB类型实例，用于后面的数据库读写操作。
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy:                           schema.NamingStrategy{},
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true, //禁用外键
	})
	if err != nil {
		return err
	}
	fmt.Println("连接成功, db:", DB)
	return nil
}
