package ossProject

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/viper"
)

var Bucket *oss.Bucket

func Init() error {
	endpoint := viper.GetString("aliyun.oss.file.endpoint")
	accessKeyId := viper.GetString("aliyun.oss.file.keyid")
	accessKeySecret := viper.GetString("aliyun.oss.file.keysecret")
	bucketName := viper.GetString("aliyun.oss.file.bucketname")

	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	Bucket, err = client.Bucket(bucketName)
	return nil
}
