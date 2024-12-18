package readMessage

import (
	"fmt"
	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
	"github.com/gin-gonic/gin"
)

// GetJsonvalue 解析为json数据
func GetJsonvalue(c *gin.Context) (*jsonvalue.V, error) { //获取原始数据
	b, _ := c.GetRawData()
	//转化为jsonvalue
	j, err := jsonvalue.Unmarshal(b)
	if err != nil {
		fmt.Println("analyzeToMap() utils.readMessage.Unmarshal() err = ", err)
		return nil, err
	}
	//解析成功，则返回 *jsonvalue,V
	return j, nil
}
