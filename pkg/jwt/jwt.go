package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 我们这里需要额外记录一个UserID字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	UserID uint64 `json:"user_id"`
	jwt.StandardClaims
}

var mySecret = []byte("golanger")

const TokenExpireDuration = time.Hour * 24 * 365

func keyFunc(_ *jwt.Token) (i interface{}, err error) {
	return mySecret, nil
}
func GenToken(userID uint64) (aToken, rToken string, err error) {
	// 创建一个我们自己的声明
	c := MyClaims{
		userID, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer:    "bluebell",                                 // 签发人
		},
	}
	// 加密并获得完整的编码后的字符串token
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(mySecret)

	// refresh token 不需要存任何自定义数据
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Second * 30).Unix(), // 过期时间
		Issuer:    "bluebell",                              // 签发人
	}).SignedString(mySecret)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return aToken, rToken, nil
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (claims *MyClaims, err error) {
	// 解析token
	var token *jwt.Token
	claims = new(MyClaims)
	token, err = jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		return
	}
	if !token.Valid { // 校验token
		err = errors.New("invalid token")
	}
	return
}
