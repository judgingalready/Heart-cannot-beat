package controller

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const key = "secret"

func GenerateToken(userName string) (string, error) {
	// 生成 Token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = userName
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenString string, user *User) error {
	// 验证 Token
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// 返回签名的 key
		return []byte(key), nil
	})
	if err != nil {
		return err
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		// 验证过期时间
		exp := int64(claims["exp"].(float64))
		if exp < time.Now().Unix() {
			return errors.New("exp: exp < time.Now().Unix()")
		}
		// 获取自定义信息
		userName := claims["username"].(string)
		userExitErr := db.Where("name = ?", userName).Take(&user).Error
		return userExitErr
	} else {
		return errors.New("invalid: token is invalid")
	}
}
