package helper

import (
	"errors"
	"esense/model"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var privateKey = []byte(os.Getenv("JWT_PRIVATE_KEY"))

func GenerateJWT(user model.User) (string, error) {
	// TTL stands for time-to-live; Duration of the validity of a token
	tokenTTL, _ := strconv.Atoi(os.Getenv("TOKEN_TTL"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"iat": time.Now().Unix(),
		"eat": time.Now().Add(time.Second * time.Duration(tokenTTL)).Unix(),
	})
	return token.SignedString(privateKey)
}

// ensures that the incoming request contains a valid token in the request header.
// This function will be used by the middleware to ensure that only authenticated requests are allowed past the middleware.
func ValidateJWT(context *gin.Context) error {
	token, err := extractToken(context)
	if err != nil {
		return err
	}
	_, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		return nil
	}
	return errors.New("invalid token provided")
}

// used to get the user associated with the provided JWT by retrieving the id key from the parsed JWT and retrieve the corresponding user from the database
func GetAuthorizedUserByJWT(context *gin.Context) (model.User, error) {
	err := ValidateJWT(context)
	if err != nil {
		return model.User{}, err
	}
	token, _ := extractToken(context)
	claims, _ := token.Claims.(jwt.MapClaims)
	userId := uint(claims["id"].(float64))

	user, err := model.FindUserById(userId)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

// uses the returned token string to parse the JWT, using the private key specified in .env.local.
func extractToken(context *gin.Context) (*jwt.Token, error) {
	tokenString := getTokenFromRequest(context)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return privateKey, nil
	})
	return token, err
}

func getTokenFromRequest(context *gin.Context) string {
	bearerToken := context.Request.Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}
