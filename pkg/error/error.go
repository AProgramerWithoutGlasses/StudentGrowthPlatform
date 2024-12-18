package error

import (
	"fmt"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)
import res "studentGrow/pkg/response"

type Error struct {
	Code res.Code
	Msg  string
}

// 实现 Error 方法
func (e *Error) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s", e.Code, e.Msg)
}

var ErrNotFoundError = errors.New("record not found")
var OverstepCompetence = errors.New("permission is out of bounds")

// HasExistError 错误--数据内容重复冲突
func HasExistError() *Error {
	return &Error{
		Code: res.ServerErrorCode,
		Msg:  "data conflict",
	}
}

// RejectRepeatSubmission 拒绝重复提交
func RejectRepeatSubmission() *Error {
	return &Error{
		Code: res.ServerErrorCode,
		Msg:  "reject repeat submission",
	}
}

// DataFormatError 数据格式错误
func DataFormatError() *Error {
	return &Error{
		Code: res.ServerErrorCode,
		Msg:  "data format error",
	}
}

// OverstepCompetence 权限越界
//func OverstepCompetence() *Error {
//	return &Error{
//		Code: res.UnprocessableEntity,
//		Msg:  "permission is out of bounds OR not found",
//	}
//}

// CheckErrors 一键检查错误,并返回给客户端msg
func CheckErrors(err error, c *gin.Context) {
	if errors.Is(err, DataFormatError()) {
		// 前端发送数据格式错误
		res.ResponseErrorWithMsg(c, res.ServerErrorCode, DataFormatError().Msg)
		return
	}
	if errors.Is(err, HasExistError()) {
		// 数据已存在，发生冲突
		res.ResponseErrorWithMsg(c, res.ServerErrorCode, HasExistError().Msg)
		return
	}

	if errors.Is(err, RejectRepeatSubmission()) {
		// 拒绝重复提交
		res.ResponseErrorWithMsg(c, res.ServerErrorCode, RejectRepeatSubmission().Msg)
		return
	}

	//if errors.Is(err, OverstepCompetence()) {
	//	// 权限越界
	//	res.ResponseErrorWithMsg(c, res.UnprocessableEntity, OverstepCompetence().Msg)
	//	return
	//}

	if errors.Is(err, ErrNotFoundError) {
		// 错误捕获-分页查询时找不到数据
		res.ResponseSuccessWithMsg(c, ErrNotFoundError.Error(), []struct{}{})
		return
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 错误捕获-gorm数据库找不到数据错误
		res.ResponseSuccessWithMsg(c, ErrNotFoundError.Error(), struct{}{})
		return
	}

	//if errors.Is(err, OverstepCompetence) {
	//	// 错误捕获-权限越界
	//	res.ResponseError(c, OverstepCompetence.Error(), []struct{}{})
	//	return
	//}
	// 其他错误
	res.ResponseErrorWithMsg(c, res.ServerErrorCode, err.Error())
}
