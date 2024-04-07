package auth

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/zondaf12/ecom-api/config"
	"github.com/zondaf12/ecom-api/types"
	"github.com/zondaf12/ecom-api/utils"
	"golang.org/x/net/context"
)

type contextKey string

const UserKey contextKey = "userId"

func CreateJWT(secret []byte, userID int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    strconv.Itoa(userID),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func WithJWTAuth(next echo.HandlerFunc, store types.UserStore) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get Token from request
		tokenString := getTokenFromRequest(c)

		// Validate JWT Token
		token, err := validateToken(tokenString)
		if err != nil {
			log.Printf("error validating token: %v", err)
			return permissionDenied(c)
		}

		if !token.Valid {
			log.Printf("token is invalid")
			return permissionDenied(c)
		}

		// Get User ID from JWT Token if valid
		claims := token.Claims.(jwt.MapClaims)
		str := claims["userId"].(string)

		userId, _ := strconv.Atoi(str)

		u, err := store.GetUserByID(userId)
		if err != nil {
			log.Printf("error getting user: %v", err)
			return permissionDenied(c)
		}

		// set context with user ID
		ctx := c.Request().Context()
		ctx = context.WithValue(ctx, UserKey, u.ID)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

func getTokenFromRequest(c echo.Context) string {
	tokenAuth := c.Request().Header.Get("Authorization")
	if tokenAuth != "" {
		return tokenAuth
	}

	return ""
}

func validateToken(t string) (*jwt.Token, error) {
	return jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Envs.JWTSecret), nil
	})
}

func permissionDenied(c echo.Context) error {
	return utils.WriteError(c, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func GetUserIDFromContext(ctx context.Context) int {
	userId, ok := ctx.Value(UserKey).(int)
	if !ok {
		return -1
	}

	return userId
}
