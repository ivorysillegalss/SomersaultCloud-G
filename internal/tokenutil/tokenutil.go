package tokenutil

import (
	"SomersaultCloud/constant/common"
	"SomersaultCloud/domain"
	"fmt"
	jwt "github.com/golang-jwt/jwt/v4"
	"time"
)

func CreateAccessToken(user *domain.User, secret string, expiry int) (accessToken string, err error) {
	exp := time.Now().Add(time.Hour * time.Duration(expiry)).Unix()
	claims := &domain.JwtCustomClaims{
		Name: user.Name,
		ID:   user.ID.Hex(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return t, err
}

func CreateRefreshToken(user *domain.User, secret string, expiry int) (refreshToken string, err error) {
	claimsRefresh := &domain.JwtCustomRefreshClaims{
		ID: user.ID.Hex(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(expiry)).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)
	rt, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return rt, err
}

func IsAuthorized(requestToken string, secret string) (bool, error) {
	_, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func ExtractIDFromToken(requestToken string, secret string) (string, error) {
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok && !token.Valid {
		return "", fmt.Errorf("Invalid Token")
	}

	return claims["id"].(string), nil
}

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
	signature := []byte(env.JwtSecretToken)
	//从配置文件中读取数字签名 可以放在全局变量那 但是会读不到值
	return signature, nil
}

func DecodeToId(tokenString string) (int, error) {
	token, err := parseJWT(tokenString)
	if err != nil {
		fmt.Println("Error parsing token:", err)
		return common.FalseInt, err
	}

	var intId int
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 可以直接访问claims里的信息，例如用户ID
		if id, ok := claims["uid"].(float64); ok {
			intId = int(id)
		}
	} else {
		fmt.Println("Invalid token")
		return common.FalseInt, err
	}
	return intId, nil
}
