package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"mini-gpt/constant"
	"mini-gpt/setting"
)

// 解析JWT
func parseJWT(tokenString string) (*jwt.Token, error) {
	// 解析并验证JWT。
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
	signature := []byte(setting.Conf.JwtSecretKey)
	//从配置文件中读取数字签名 可以放在全局变量那 但是会读不到值
	return signature, nil
}

func DecodeToId(tokenString string) (int, error) {
	token, err := parseJWT(tokenString)
	if err != nil {
		fmt.Println("Error parsing token:", err)
		return constant.FalseInt, err
	}

	var intId int
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 可以直接访问claims里的信息，例如用户ID
		if id, ok := claims["uid"].(float64); ok {
			intId = int(id)
		}
	} else {
		fmt.Println("Invalid token")
		return constant.FalseInt, err
	}
	return intId, nil
}
