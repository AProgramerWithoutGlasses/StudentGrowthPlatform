package middleWare

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

//"http://192.168.10.7":   true,
//"http://192.168.10.7:81": true,
//"http://8.154.36.180":   true,
//"http://8.154.36.180:81": true,

// 跨域
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 定义一个允许跨域的域名列表
		var allowedOrigins = map[string]bool{
			"http://192.168.10.7":      true,
			"http://192.168.10.7:81":   true,
			"http://8.154.36.180:8904": true,
			"http://8.154.36.180:8905": true,
			"http://127.0.0.1:8881":    true,
		}
		method := c.Request.Method               //请求方法
		origin := c.Request.Header.Get("Origin") //请求头部
		var headerKeys []string                  // 声明请求头keys
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}

		fmt.Println("len map", len(allowedOrigins))
		fmt.Println("map", allowedOrigins)
		fmt.Println("origin-->", origin, "-->has Map", allowedOrigins[origin])
		// 检查请求的Origin是否在允许的域名列表中

		fmt.Println("mode:-->", viper.GetString("app.mode"))
		if allowedOrigins[origin] || viper.GetString("app.mode") == "dev" {

			fmt.Println("设置了权限--》", origin)
			// 如果是，则设置Access-Control-Allow-Origin为请求的Origin
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			// 如果请求的Origin不在允许的列表中，你可以选择拒绝请求或返回默认值
			// 例如，设置为一个固定的域名或返回空字符串（不推荐）
			// c.Header("Access-Control-Allow-Origin", "")
			// 为了安全起见，在此情况下最好直接返回错误或拒绝请求
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Origin not allowed"})
			return
		}

		//fmt.Println(c.ClientIP())
		//if !strings.Contains(c.ClientIP(), "192.168") {
		//	c.Abort()
		//	return
		//}

		if origin != "" {
			//c.Header("Access-Control-Allow-Origin", "https://8.154.36.180, https://192.168.10.7") // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//     允许跨域设置                                                                                                      可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "false")                                                                                                                                                  //  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")                                                                                                                                                              // 设置返回格式是json
			c.Set("content-disposition", "inline")

		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next() //  处理请求
	}
}
