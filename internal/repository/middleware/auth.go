package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/ypxd99/yandex-practicm/util"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookieName := util.GetConfig().Auth.CookieName
		cookie, err := c.Cookie(cookieName)
		logger := util.GetLogger()

		cfg := util.GetConfig()
		secretKey := []byte(cfg.Auth.SecretKey)

		if err != nil || !isValidToken(cookie, secretKey) {
			newUserID := uuid.New()
			userIDStr := newUserID.String()

			token, err := generateToken(userIDStr, secretKey)
			if err != nil {
				logger.Errorf("failed to generate token: %v", err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			c.SetCookie(
				cookieName,
				token,
				3600*24*30,
				"/",
				"",
				false,
				true,
			)

			c.Set(cookieName, newUserID)
			logger.Infof("created new user with ID: %s", userIDStr)
		} else {
			userIDStr, err := extractUserIDFromToken(cookie, secretKey)
			if err != nil {
				newUserID := uuid.New()
				userIDStr = newUserID.String()

				token, err := generateToken(userIDStr, secretKey)
				if err != nil {
					logger.Errorf("failed to generate token: %v", err)
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}

				c.SetCookie(
					cookieName,
					token,
					3600*24*30,
					"/",
					"",
					false,
					true,
				)

				c.Set(cookieName, newUserID)
				logger.Infof("regenerated token for user, new ID: %s", userIDStr)
				c.Next()
				return
			}

			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				userID = uuid.New()
				userIDStr = userID.String()

				token, err := generateToken(userIDStr, secretKey)
				if err != nil {
					logger.Errorf("failed to generate token: %v", err)
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}

				c.SetCookie(
					cookieName,
					token,
					3600*24*30,
					"/",
					"",
					false,
					true,
				)

				logger.Infof("regenerated UUID for user, new ID: %s", userIDStr)
			}

			c.Set(cookieName, userID)
			logger.Infof("authenticated user with ID: %s", userIDStr)
		}

		c.Next()
	}
}

func generateToken(userID string, key []byte) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func isValidToken(tokenString string, key []byte) bool {
	_, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return key, nil
	})

	return err == nil
}

func extractUserIDFromToken(tokenString string, key []byte) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return key, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, nil
	}

	return "", errors.New("invalid token claims")
}

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookieName := util.GetConfig().Auth.CookieName
		userID, exists := c.Get(cookieName)
		if !exists {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		_, ok := userID.(uuid.UUID)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}

func GetUserID(c *gin.Context) (uuid.UUID, error) {
	cookieName := util.GetConfig().Auth.CookieName
	userID, exists := c.Get(cookieName)
	if !exists {
		return uuid.Nil, errors.New("user ID not found")
	}
	return userID.(uuid.UUID), nil
}
