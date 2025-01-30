package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
)

var mySigningKey = []byte("secret") // 设置密钥用于签名

type UserAuth struct {
	Id      uint   `json:"id"`
	Name    string `json:"name"`
	IsAdmin bool   `json:"is_admin"`
	jwt.StandardClaims
}

// GenerateToken 生成Token
func GenerateToken(id uint, name string) (string, error) {
	userAuth := UserAuth{
		Id:             id,
		Name:           name,
		StandardClaims: jwt.StandardClaims{},
	}

	// 从声明创建token令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userAuth)

	// 生成字符串形式的token令牌
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		log.Println("generate token string error: ", err)
		return "", err
	}
	return tokenString, nil
}

// AnalyseToken 解析Token字符串，返回用户声明
func AnalyseToken(tokenString string) (*UserAuth, error) {
	userAuth := new(UserAuth)

	claims, err := jwt.ParseWithClaims(tokenString, userAuth, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	if err != nil {
		log.Println("analyse token string error: ", err)
		return nil, err
	}

	if !claims.Valid {
		return nil, fmt.Errorf("analyse token string error: %v", err)
	}
	return userAuth, nil
}
