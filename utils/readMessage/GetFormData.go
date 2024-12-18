package readMessage

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	res "studentGrow/pkg/response"
)

// GetFormData 获取formData
func GetFormData(c *gin.Context) (m map[string]any, err error) {
	form, err := c.MultipartForm()
	if err != nil {
		fmt.Println("GetFormData() utils.readMessage.MultipartForm() err = ", err)
		return nil, err
	}
	//如果文本为空，则不允许提交
	if form.Value["word_count"][0] == "0" {
		res.ResponseErrorWithMsg(c, res.UnprocessableEntity, "文本为空，不允许提交")
		return nil, errors.New("文本为空，不允许提交")
	}

	// 解析到map
	for key, val := range form.Value {
		m[key] = val[0]
	}

	//获取文件
	files := form.File

	// 若没有文件，则返回错误"noFile"
	if len(files) == 0 {
		return m, errors.New("noFile")
	}

	//否则将文件返回	---files["files"]是一个包含多个文件的切片
	//---注意：客户端传来的文件字段必须是"files"，可上传多个文件
	m["files"] = files["files"]

	//将所有数据解析到map完成，返回map
	return m, nil
}
