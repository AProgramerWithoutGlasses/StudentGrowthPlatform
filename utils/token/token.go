package token

import (
	"github.com/dgrijalva/jwt-go"
	"studentGrow/models"
	"studentGrow/models/gorm_model"
	"time"
)

// JwtKey 定义密钥
var JwtKey = []byte("xlszxjm")

// ReleaseToken 生成密钥
func ReleaseToken(user gorm_model.User) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour) //token的有效期是七天
	claims := &models.Claims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), //token的有效期
			IssuedAt:  time.Now().Unix(),     //token发放的时间
			Issuer:    "xyq",                 //作者
			Subject:   "user token",          //主题
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey) //根据前面自定义的Jwt秘钥生成token
	tokenString = "Bearer " + tokenString
	if err != nil {
		//返回生成的错误
		return "", err
	}
	//返回成功生成的字符换
	return tokenString, nil
}
