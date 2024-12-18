package routesJoinAudit

import (
	"embed"
	_ "embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nguyenthenguyen/docx"
	"go.uber.org/zap"
	"io"
	"io/fs"
	"net/url"
	"os"
	"strconv"
	"strings"
	"studentGrow/dao/mysql"
	"studentGrow/pkg/response"
	token2 "studentGrow/utils/token"
)

//go:embed template.docx
var data embed.FS

type resList struct {
	ListName string `json:"list_name"`
	IsFinish bool   `json:"is_finish"`
}

type rec struct {
	ActivityID int    `json:"activity_id"`
	CurMenu    string `json:"cur_menu"`
}

func ExportFormMsg(c *gin.Context) {
	token := token2.NewToken(c)
	_, exist := token.GetUser()
	if !exist {
		response.ResponseError(c, response.TokenError)
		zap.L().Error("token错误")
		return
	}
	var cr rec
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		response.ResponseErrorWithMsg(c, response.ParamFail, "json解析失败")
		return
	}
	userMsgList, _ := mysql.UserListWithOrganizer(cr.ActivityID, cr.CurMenu)
	var isFinish = true
	if len(userMsgList) == 0 {
		isFinish = false
	}
	response.ResponseSuccess(c, resList{
		ListName: cr.CurMenu,
		IsFinish: isFinish,
	})
}
func ExportFormFile(c *gin.Context) {
	token := token2.NewToken(c)
	_, exist := token.GetUser()
	if !exist {
		response.ResponseError(c, response.TokenError)
		zap.L().Error("token错误")
		return
	}
	var cr rec
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		response.ResponseErrorWithMsg(c, response.ParamFail, "json解析失败")
		return
	}
	_, _, activityMsg := mysql.OpenActivityStates()
	title := activityMsg.ActivityName
	userMsgList, classList := mysql.UserListWithOrganizer(cr.ActivityID, cr.CurMenu)
	responseDOCX(userMsgList, c, title, classList)
	return
}

func responseDOCX(dynamicList []map[string]interface{}, c *gin.Context, title string, classList []string) {
	rewriteList, classified := classifyWordListWithClass(dynamicList)
	// 读取模板文件
	templateFile, _ := data.Open("template.docx")
	defer func() {
		_ = templateFile.Close()
	}()
	//复制文件到当前目录下
	copyFileToPath(templateFile, "./template.docx")
	//读取模版
	r, err := docx.ReadDocxFile("./template.docx")
	if err != nil {
		response.ResponseErrorWithMsg(c, response.ServerErrorCode, "文件打开失败")
		zap.L().Error(err.Error())
		return
	}
	defer func() {
		_ = r.Close()
	}()
	// 使文档可编辑
	docx1 := r.Editable()
	// 替换模板中的占位符
	err = docx1.Replace("{{LIST}}", rewriteList, -1)
	if err != nil {
		zap.L().Error(err.Error())
	}
	var classNameList []string

	//循环遍历学生名单填入对应的list中
	for k, v := range classList {
		classStr := "CLASS"
		nameStr := "NAME"
		classStr = "{{" + classStr + strconv.Itoa(k) + "}}"
		nameStr = "{{" + nameStr + strconv.Itoa(k) + "}}"
		err = docx1.Replace(classStr, v+":", -1)
		if err != nil {
			zap.L().Error(err.Error())
		}
		classNameList = make([]string, 0)
		for _, value := range classified[v] {
			userName, _ := value["name"].(string)
			classNameList = append(classNameList, userName)
		}
		rewriteNameList := strings.Join(classNameList, "\t")
		err = docx1.Replace(nameStr, rewriteNameList, -1)
		if err != nil {
			zap.L().Error(err.Error())
		}
	}
	err = docx1.Replace("{{Title}}", title, -1)
	if err != nil {
		zap.L().Error(err.Error())
	}
	filename := url.QueryEscape(title + ".docx")

	//添加响应头
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename*=utf-8''%s", filename))
	//将结果写到响应中
	err = docx1.Write(c.Writer)
	if err != nil {
		zap.L().Error(err.Error())
	}
	return
}

// 根据班级把查询到的名字分类
func classifyWordListWithClass(dynamicList []map[string]interface{}) (rewriteList string, classified map[string][]map[string]interface{}) {
	classified = make(map[string][]map[string]interface{})
	for _, v := range dynamicList {
		userClass, _ := v["user_class"].(string)
		classified[userClass] = append(classified[userClass], v)
	}
	n := len(classified)
	replaceNameList := make([]string, 0)
	for i := 0; i < n; i++ {
		classStr := "CLASS"
		nameStr := "NAME"
		classStr = "{{" + classStr + strconv.Itoa(i) + "}}"
		nameStr = "{{" + nameStr + strconv.Itoa(i) + "}}"
		replaceNameList = append(replaceNameList, "\n"+classStr)
		replaceNameList = append(replaceNameList, nameStr)
	}
	rewriteList = strings.Join(replaceNameList, "\n")
	return
}

// 复制打包到可执行文件中的文件到同级目录中
func copyFileToPath(templateFile fs.File, path string) {
	newFile, err := os.Create(path)
	if err != nil {
		zap.L().Error(err.Error())
	}
	defer func() {
		_ = newFile.Close()
	}()

	_, err = io.Copy(newFile, templateFile)
	if err != nil {
		zap.L().Error(err.Error())
	}
	return
}
