package message

import (
	"fmt"
	"go.uber.org/zap"
	"studentGrow/dao/mysql"
	"studentGrow/models/gorm_model"
	myErr "studentGrow/pkg/error"
)

// GetUnreadReportsForService 获取举报信息列表
func GetUnreadReportsForService(username, role string, limit, page int) (reports []gorm_model.UserReportArticleRecord, err error) {

	// 选择管理员角色
	switch role {
	case "class":
		reports, err = mysql.GetUnreadReportsForClass(username, limit, page)
	case "grade1":
		reports, err = mysql.GetUnreadReportsForGrade(1, limit, page)
	case "grade2":
		reports, err = mysql.GetUnreadReportsForGrade(2, limit, page)
	case "grade3":
		reports, err = mysql.GetUnreadReportsForGrade(3, limit, page)
	case "grade4":
		reports, err = mysql.GetUnreadReportsForGrade(4, limit, page)
	case "college":
		reports, err = mysql.GetUnreadReportsForSuperman(limit, page)
	case "superman":
		reports, err = mysql.GetUnreadReportsForSuperman(limit, page)
	default:
		return nil, myErr.ErrNotFoundError
	}
	if err != nil {
		zap.L().Error("GetUnreadReportsForClassService() service.article.GetUnreadReports err=", zap.Error(err))
		return nil, err
	}

	return reports, nil
}

// GetUnreadReportNumForService 获取未读举报信息的数目
func GetUnreadReportNumForService(username, role string) (count int, err error) {
	fmt.Println(role)
	switch role {
	case "class":
		count, err = mysql.GetUnreadReportNumForClass(username)
	case "grade1":
		count, err = mysql.GetUnreadReportNumForGrade(1)
	case "grade2":
		count, err = mysql.GetUnreadReportNumForGrade(2)
	case "grade3":
		count, err = mysql.GetUnreadReportNumForGrade(3)
	case "grade4":
		count, err = mysql.GetUnreadReportNumForGrade(4)
	case "college":
		count, err = mysql.GetUnreadReportNumForSuperman()
	case "superman":
		count, err = mysql.GetUnreadReportNumForSuperman()
	default:
		return -1, myErr.ErrNotFoundError
	}
	if err != nil {
		zap.L().Error("GetUnreadRoportNumForService() service.article.GetUnreadReports err=", zap.Error(err))
		return -1, err
	}
	return count, nil
}

// AckUnreadReportsService 确认举报信息
func AckUnreadReportsService(reportId int, username string, role string) (err error) {
	// 选择管理员角色
	switch role {
	case "class":
		err = mysql.AckUnreadReportsForClass(reportId, username)
	case "grade1":
		err = mysql.AckUnreadReportsForGrade(reportId, 1)
	case "grade2":
		err = mysql.AckUnreadReportsForGrade(reportId, 2)
	case "grade3":
		err = mysql.AckUnreadReportsForGrade(reportId, 3)
	case "grade4":
		err = mysql.AckUnreadReportsForGrade(reportId, 4)
	case "college":
		err = mysql.AckUnreadReportsForSuperman(reportId)
	case "superman":
		err = mysql.AckUnreadReportsForSuperman(reportId)
	default:
		return myErr.DataFormatError()
	}
	if err != nil {
		zap.L().Error("AckUnreadReportsService() service.article.GetUnreadReports err=", zap.Error(err))
		return err
	}

	return nil
}
