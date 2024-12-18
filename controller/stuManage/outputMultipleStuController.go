package stuManage

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"studentGrow/models/jrx_model"
	"studentGrow/pkg/response"
	"studentGrow/service"
)

// 批量导出
func OutputMultipleStuControl(c *gin.Context) {
	// 接收请求
	var selectedStuMesStruct jrx_model.SelectedStuMesStruct
	err := c.Bind(&selectedStuMesStruct)
	if err != nil {
		response.ResponseErrorWithMsg(c, 500, "stuManager.OutputMultipleStuControl() c.Bind() failed : "+err.Error())
		zap.L().Error("stuManager.OutputMultipleStuControl() c.Bind() failed : ", zap.Error(err))
		return
	}

	// 处理业务
	excelData, err := service.GetSelectedStuExcel(selectedStuMesStruct)
	if err != nil {
		response.ResponseErrorWithMsg(c, 500, "stuManager.OutputMultipleStuControl() service.GetSelectedStuExcel() failed : "+err.Error())
		zap.L().Error("stuManager.OutputMultipleStuControl() service.GetSelectedStuExcel() failed : ", zap.Error(err))
		return
	}

	// 响应
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", url.QueryEscape("批量导出学生信息.xlsx")))
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelData.Bytes())

}
