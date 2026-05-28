package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// ==================== JWT 配置 ====================
// 密钥（生产环境请放到配置文件/环境变量）
var jwtSecret = []byte("your-strong-secret-key-2025")

// CustomClaims 自定义token内容

// AuthMiddleware ==================== 中间件：Token 认证 ====================
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取 token
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" || len(tokenStr) < 7 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "未登录或token不存在",
			})
			c.Abort() // 终止请求
			return
		}

		// 去掉 "Bearer " 前缀
		tokenStr = tokenStr[7:]

		// 解析 token
		claims, err := ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "token无效或已过期",
			})
			c.Abort()
			return
		}

		// 将用户信息存入上下文，后续接口可直接使用
		c.Set("userID", claims.UserID)
		c.Next()
	}
}

// CustomClaims ==================== JWT 工具函数 ====================
// 1. 把 UserID 类型从 uint 改成 string
type CustomClaims struct {
	UserID string `json:"user_id"` // 这里改成 string
	jwt.RegisteredClaims
}

// GenerateToken 2. 函数参数也改成 string
func GenerateToken(userID string) (string, error) {
	expireTime := time.Now().Add(24 * time.Hour)
	claims := CustomClaims{
		UserID: userID, // 直接传 UUID 字符串
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken 解析 Token
func ParseToken(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	// 断言
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}

// HashPassword 加密密码（注册和修改时候时用）
func HashPassword(pwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 校验密码（登录时用）
func CheckPassword(hashedPwd, inputPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(inputPwd))
	return err == nil
}
