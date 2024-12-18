package fileProcess

import (
	"fmt"
	uuid2 "github.com/google/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mime/multipart"
	"studentGrow/aliyun/oss"
	"time"
)

// UploadFile 将文件上传至oss
func UploadFile(fileType string, fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		zap.L().Error("UploadFile() utils.fileProcess.Open err=", zap.Error(err))
		return "", err
	}
	defer func(file multipart.File) {
		err = file.Close()
		if err != nil {
			zap.L().Error("UploadFile() utils.fileProcess.Close err=", zap.Error(err))
			return
		}
	}(file)
	// 获取文件名称
	fileName := fileHeader.Filename
	// 随机生成唯一值，防止文件覆盖
	uuid := uuid2.NewString()
	fileName = uuid + fileName
	// 获取当前时间
	datePath := time.Now().Format("2006/01/02")
	// 拼接，按日期文件夹分类
	fileName = fileType + "/" + datePath + "/" + fileName
	// 创建请求
	err = ossProject.Bucket.PutObject(fileName, file)
	if err != nil {
		zap.L().Error("UploadFile() utils.fileProcess.PutObject err=", zap.Error(err))
		return "", err
	}
	url := fmt.Sprintf("https://%s.%s/%s", viper.GetString("aliyun.oss.file.bucketname"), viper.GetString("aliyun.oss.file.endpoint"), fileName)
	return url, nil
}

// 软删除
//func DelOssFile(src string) (error, http.Header) {
//	var retHeader http.Header
//	err := ossProject.Bucket.DeleteObject(src, oss.GetResponseHeader(&retHeader))
//	if err != nil {
//		log.Printf("Failed to delete object '%s': %v", src, err)
//	}
//	// 打印删除标记信息。
//	log.Printf("x-oss-version-id: %s", oss.GetVersionId(retHeader))
//	log.Printf("x-oss-delete-marker: %t", oss.GetDeleteMark(retHeader))
//	log.Println("Object deleted successfully.")
//	err = ossProject.Bucket.DeleteObject(src, oss.VersionId(oss.GetVersionId(retHeader)))
//	if err != nil {
//		log.Fatalf("Failed to delete object '%s' with version ID '%s': %v", src, oss.GetVersionId(retHeader), err)
//	}
//
//	log.Println("Object deleted successfully.")
//	return err, retHeader
//}

//彻底删除
