package tokens

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Claims struct {
	UserID string `json:"_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

func NewJWTManager(secretKey string, duration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     secretKey,
		tokenDuration: duration,
	}
}

func (j *JWTManager) Generate(userID, email string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTManager) Validate(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, errors.New("неверный токен")
	}

	claims, ok := token.Claims.(*Claims)
	if ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("недействительный токен")
	}
}
