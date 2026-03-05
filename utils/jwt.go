package utils

import (
	"time"

	"github.com/fakhri-rasyad/wpu_goreact/config"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// generate token
func GenerateJWTToken(userID int64, role , email string, publicID uuid.UUID) (string, error){
	secret := config.APPConfig.JWTSecret
	duration, _ := time.ParseDuration(config.APPConfig.JWTExpire)

	claims := jwt.MapClaims{"user_id" : userID, "role" : role, "pub_id" : publicID, "email" : email, "exp" : time.Now().Add(duration).Unix()}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func RefreshJWTToken(userID int64) (string, error) {
	secret := config.APPConfig.JWTSecret
	duration, _ := time.ParseDuration(config.APPConfig.JWTRefreshToken)

	claims := jwt.MapClaims{"user_id" : userID, "exp" : time.Now().Add(duration).Unix()}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ExtractClaims(JWTtoken string)(jwt.MapClaims, bool){
	secret := config.APPConfig.JWTSecret
	hmac := []byte(secret)
	token, err := jwt.Parse(JWTtoken, func(t *jwt.Token) (any, error) {
		return hmac, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid{
		return claims, true
	} else {
		return nil, false
	}
}