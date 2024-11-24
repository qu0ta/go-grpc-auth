package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/qu0ta/go-grpc-auth/internal/domain/models"
	"time"
)

func NewToken(user models.User, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func ParseToken(token string, app models.App) (models.User, error) {
	tokenData := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, tokenData, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(app.Secret), nil
	})
	if err != nil {
		return models.User{}, err
	}
	return models.User{
		ID:    int64(tokenData["uid"].(float64)),
		Email: tokenData["email"].(string),
	}, nil
}
