package common

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/wcc4869/ginessential/model"
	"time"
)

var jwtKey = []byte("a_secret_create")

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

func ReleaseToken(user model.User) (string, error) {
	// 有效期 7 天
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), //过期时间
			IssuedAt:  time.Now().Unix(),
			Issuer:    "pywcc.top",
			Subject:   "user token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
	//	eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9
	//	.eyJVc2VySWQiOjQsImV4cCI6MTU5MDMzMzU3NCwiaWF0IjoxNTg5NzI4Nzc0LCJpc3MiOiJ3Y2M0ODY5LmNvbSIsInN1YiI6InVzZXIgdG9rZW4ifQ
	//	.-jKmkVdLFPUgDZpQ1oO2Kb0SzIX-4B37hMHt_o-SgG4
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})

	return token, claims, err

}
