package services

import (
	"emoji-maker-backend/config"
	"emoji-maker-backend/models"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService JWT服务接口
type JWTService interface {
	GenerateToken(user *models.User) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
}

// jwtServiceImpl JWT服务实现
type jwtServiceImpl struct {
	secretKey string
}

// NewJWTService 创建JWT服务实例
func NewJWTService() JWTService {
	return &jwtServiceImpl{
		secretKey: config.AppConfig.JWT.Secret,
	}
}

// Claims 自定义JWT claims
type Claims struct {
	ID    int64  `json:"id"`
	Phone string `json:"phone"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT
func (s *jwtServiceImpl) GenerateToken(user *models.User) (string, error) {
	claims := &Claims{
		ID:    user.ID,
		Phone: user.Phone,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)), // Token有效期72小时
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "emoji-maker-backend",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

// ValidateToken 验证JWT
func (s *jwtServiceImpl) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secretKey), nil
	})
}
