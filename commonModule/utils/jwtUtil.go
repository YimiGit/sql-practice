package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

// 类似盐值的东西
var key = "yimiFlower"

func CreateToken(userId int64) (string, error) {
	claims := myClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			//ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Minute)),
			Issuer: key,
			//NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	// 生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 通过盐值签名
	signingKey := []byte(key)
	if tokenString, err := token.SignedString(signingKey); err != nil {
		//直接返回，后续判断err再log
		return "", err
	} else {
		return tokenString, nil
	}
}

func ParseToken(tokenString string) (int64, error) {
	claims := myClaims{}
	//解析token
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		return []byte(key), nil
	})
	if err != nil {
		return http.StatusForbidden, err
	}
	// 校验token
	if token.Valid {
		return claims.UserId, nil
	}
	return http.StatusForbidden, err
}

type myClaims struct {
	UserId int64 `json:"userId"`
	jwt.RegisteredClaims
}
