package db

import (
	"errors"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	// SignKey 签名
	SignKey string = "douyin-douyin"
)

// ErrTokenExpired 过期
var ErrTokenExpired = errors.New("Token is expired")

// CustomClaims 载荷，可以加一些自己需要的信息
type CustomClaims struct {
	// UserID   string `json:"user_id"`
	// Username string `json:"username"`
	// IsVip    int    `json:"is_vip"`

	Uid      int
	Username string
	Token    string
	jwt.StandardClaims
}

// JWT 签名结构
type JWT struct {
	SigningKey []byte
}

// NewJWT 新建一个jwt实例
func NewJWT() *JWT {
	return &JWT{
		[]byte(SignKey),
	}
}

// CreateToken 生成一个token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// ParseToken 解析Token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	//fmt.Println(err)
	if err != nil {
		log.Println(err)
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("That's even not  a token")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New("Token not active yet")
			} else {
				return nil, errors.New("Couldn't handle this token 1")
			}
		}
		return nil, errors.New("Couldn't handle this token 2")
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("Couldn't handle this token 3")
}

// RefreshToken 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now

		// 1year 过期
		claims.StandardClaims.ExpiresAt = time.Now().Add(24 * 30 * 12 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", errors.New("Couldn't handle this token")
}
