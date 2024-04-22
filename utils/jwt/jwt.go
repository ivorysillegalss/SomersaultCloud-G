package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"mini-gpt/constant"
	"mini-gpt/setting"
)

var mySigningKey = []byte(setting.Conf.JwtSecretKey)

// 解析JWT
func parseJWT(tokenString string) (*jwt.Token, error) {
	// 解析并验证JWT。注意：确保提供一个key function来验证签名算法
	//token, err := jwt.Parse(tokenString, verifySignature)
	token, err := jwt.Parse(tokenString, verifySignature)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func verifySignature(token *jwt.Token) (interface{}, error) {
	// 确保token的签名算法是我们期望的
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return mySigningKey, nil
}

func DecodeToId(tokenString string) (int, error) {
	token, err := parseJWT(tokenString)
	if err != nil {
		fmt.Println("Error parsing token:", err)
		return constant.FalseInt, err
	}

	var intId int

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
		// 可以直接访问claims里的信息，例如用户ID
		if id, ok := claims["uid"].(float64); ok {
			intId = int(id)
			fmt.Println(intId)
		}
	} else {
		fmt.Println("Invalid token")
		return constant.FalseInt, err
	}
	return intId, nil
}
