package middleWare

import "github.com/gin-gonic/gin"

func SetHTTPHeaders(c *gin.Context) {
	// 设置内容类型为图片，根据实际情况选择合适的MIME类型
	c.Writer.Header().Set("Content-Type", "image/jpeg")

	// 设置内容展示方式为内联，而不是附件下载
	c.Writer.Header().Set("Content-Disposition", "inline")

	// 调用后续的处理函数
	c.Next()
}
