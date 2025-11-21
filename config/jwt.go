package config

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	jwtKey         []byte
	jwtExpireHours int
)

func init() {
	jwtKey = []byte(GetEnv("JWT_SECRET_KEY", ""))
	expStr := GetEnv("JWT_EXPIRE_HOURS", "")

	hours, err := strconv.Atoi(expStr)
	if err != nil || hours <= 0 {
		hours = 24
	}
	jwtExpireHours = hours
}

func GenerateJWT(userID uuid.UUID, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * time.Duration(jwtExpireHours)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ParseJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return token.Claims.(jwt.MapClaims), nil
}
