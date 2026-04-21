package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 使用BCrypt算法哈希密码
// 返回哈希后的密码字符串和错误信息
func HashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), 12)
	return string(hash), err
}

// GenerateJWT 生成JSON Web Token
// 输入用户名，返回包含username声明的JWT Token，有效期72小时
func GenerateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // 72小时有效期
	})

	signedToken, err := token.SignedString([]byte("secret"))
	return "Bearer " + signedToken, err
}

// CheckPassword 验证密码是否正确
// 使用BCrypt比较明文密码和哈希密码，返回是否匹配
func CheckPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ParseJWT 解析JSON Web Token
// 从Token字符串中提取用户名声明
// 如果Token以"Bearer "开头，会自动去除
func ParseJWT(tokenString string) (string, error) {
	// 去除"Bearer "前缀
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	// 解析Token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected Signing Method")
		}
		return []byte("secret"), nil
	})

	if err != nil {
		return "", err
	}

	// 提取用户名
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username, ok := claims["username"].(string)
		if !ok {
			return "", errors.New("username claim is not a string")
		}
		return username, nil
	}

	return "", err
}
