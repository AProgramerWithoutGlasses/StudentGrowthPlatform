package JoinAudit

import (
	"fmt"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"studentGrow/dao/mysql"
	"studentGrow/models/gorm_model"
	"studentGrow/utils/fileIsInList"
	"studentGrow/utils/fileProcess"
	"studentGrow/utils/hashMd5"
)

// StuFileUpload 返回值
type StuFileUpload struct {
	ID        uint
	FileName  string `json:"file_name"`
	FilePath  string `json:"file_path"`
	IsSuccess bool   `json:"is_success"`
	Note      string `json:"note"`
	Msg       string `json:"msg"`
}

// WhiteImageList 允许上传的文件格式白名单
var (
	WhiteImageList = []string{
		"jpg",
		"png",
		"jpeg",
		"doc",
		"docx",
		"pdf",
		"webp",
	}
)

func FileUpload(file *multipart.FileHeader, user gorm_model.User, ActivityMsg gorm_model.JoinAuditDuty, note string) (res StuFileUpload) {
	fileName := file.Filename
	res.FilePath = ""
	res.FileName = fileName
	res.IsSuccess = false
	res.Note = note
	//判断文件是否在白名单内
	if !fileIsInList.FileIsInList(fileName, WhiteImageList) {
		zap.L().Info("文件格式不合法")
		res.Msg = "文件格式不合法"
		return
	}

	//设置上传文件最大大小，判断文件大小是否合法
	const maxSize = float64(50)
	fileSize := float64(file.Size) / float64(1024*1024) //计算文件大小
	if fileSize > maxSize {
		zap.L().Info("文件不合法")
		res.Msg = "当前图片大小为" + fmt.Sprintf("%f MB 超出限制 %f MB", fileSize, maxSize)
		return
	}

	fileObj, err := file.Open()
	if err != nil {
		zap.L().Error(err.Error())
		res.Msg = err.Error()
		return
	}
	byteData, err := io.ReadAll(fileObj)
	if err != nil {
		zap.L().Error(err.Error())
		res.Msg = err.Error()
		return
	}
	imageHash := hashMd5.HashMd5(byteData)
	//根据hash值判断照片是否存在于库中
	var FileMsg gorm_model.JoinAuditFile
	////查询当前开启活动的信息
	//
	//err = mysql.DB.Take(&FileMsg, "file_hash = ? AND join_audit_duty_id = ? AND  username = ?", imageHash, ActivityMsg.ID, user.Username).Error
	//if !errors.Is(err, gorm.ErrRecordNotFound) {
	//	//当前活动中用户存在该照片，重复上传
	//	res.ID = FileMsg.ID
	//	res.FileName = FileMsg.FileName
	//	res.IsSuccess = true
	//	res.FilePath = FileMsg.FilePath
	//	res.Note = FileMsg.Note
	//	res.Msg = "用户照片存在，重复上传"
	//	return
	//}

	//上传阿里云
	imagePath, err := fileProcess.UploadFile("image", file)
	if err != nil {
		res.Msg = "阿里云上传出错"
		zap.L().Info(err.Error())
		return
	}
	//入库
	FileMsg = gorm_model.JoinAuditFile{
		Username:      user.Username,
		FilePath:      imagePath,
		FileName:      fileName,
		FileHash:      imageHash,
		JoinAuditDuty: ActivityMsg, //照片归属的活动
		Note:          note,
	}
	err = mysql.DB.Create(&FileMsg).Error
	if err != nil {
		res.ID = FileMsg.ID
		res.FilePath = imagePath
		res.Msg = "文件链接入库出现错误"
		return
	}
	res.ID = FileMsg.ID
	res.IsSuccess = true
	res.FilePath = imagePath
	res.Note = FileMsg.Note
	res.Msg = "文件上传成功"
	return
}
