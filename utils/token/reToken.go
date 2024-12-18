package token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
	"studentGrow/dao/mysql"
	"studentGrow/models"
	"studentGrow/models/gorm_model"
	"studentGrow/pkg/response"
	"studentGrow/service/userService"
)

type Token struct {
	C *gin.Context
}

// AuthMiddleware 中间件检验token是否合法
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取前端传过来的信息
		tokenString := c.GetHeader("token")
		//验证前端传过来的token格式，不为空，开头为Bearer
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			response.ResponseErrorWithMsg(c, 400, "token不合法")
			c.Abort()
			return
		}
		//验证通过，提取有效部分（除去Bearer)
		tokenString = tokenString[7:] //截取字符
		//解析token
		token, _, err := ParseToken(tokenString)
		//解析失败||解析后的token无效
		if err != nil || !token.Valid {
			response.ResponseErrorWithMsg(c, 400, "token解析失败")
			c.Abort()
			return
		}
		c.Set("claim", token.Claims)
		c.Next()
	}
}

// ParseToken 解析从前端获取到的token值
func ParseToken(tokenString string) (*jwt.Token, *models.Claims, error) {
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	return token, claims, err
}

func NewToken(c *gin.Context) *Token {
	return &Token{C: c}
}

func (this *Token) GetUser() (gorm_model.User, bool) {
	claim, exist := this.C.Get("claim")
	if !exist {
		return gorm_model.User{}, false
	}
	user := claim.(*models.Claims).User
	return user, true
}

func (this *Token) GetRole() (string, error) {
	user, exist := this.GetUser()
	var role string
	if !exist {
		response.ResponseError(this.C, response.TokenError)
		zap.L().Error("token错误")
		return "", fmt.Errorf("token错误")
	}
	//验证用户是否是管理员
	newOk := userService.BVerifyExit(user.Username)
	if newOk {
		cId, err := mysql.SelCasId(user.Username)
		role, err = mysql.SelRole(cId)
		if err != nil {
			return "", err
		}
	} else {
		role = ""
	}
	return role, nil
}
