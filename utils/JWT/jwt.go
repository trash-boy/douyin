package JWT

import (
	"TinyTolk/config"
	"github.com/gin-gonic/gin"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

//自定义一个字符串
var jwtkey = config.JWT_KEY


type Claims struct {
	UserId uint
	jwt.StandardClaims
}


//颁发token
func SetToken(Id uint) (string, error){
	//log.Println("id为：", Id)
	expireTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserId: Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间
			IssuedAt:  time.Now().Unix(),
			Issuer:    "127.0.0.1",  // 签名颁发者
			Subject:   "user token", //签名主题
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//fmt.Println(token)
	tokenString, err := token.SignedString([]byte(jwtkey))
	if err != nil {
		return "", err
	}
	return tokenString,err

}


func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	Claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, Claims, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(jwtkey), nil
	})
	return token, Claims, err
}

//todo 验证token是否有效
func VerifyToken(ctx *gin.Context, tokenString string)bool{

	//vcalidate token formate
	if tokenString == "" {
		return false
	}

	token, claim, err := ParseToken(tokenString)
	if err != nil || !token.Valid {
		log.Println("verfi error", err)
		return false
	}
	if claim.Valid() != nil{
		log.Println("verfi error", err)
		return false
	}
	ctx.Set("id",claim.UserId)
	return true
}
