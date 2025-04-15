package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/ypxd99/yandex-practicm/util"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookieName := util.GetConfig().Auth.CookieName
		cookie, err := c.Cookie(cookieName)
		logger := util.GetLogger()

		cfg := util.GetConfig()
		secretKey := []byte(cfg.Auth.SecretKey)

		if err != nil || !isValidCookie(cookie, secretKey) {
			newUserID := uuid.New()
			userIDStr := newUserID.String()
			signedUserID := signCookie(userIDStr, secretKey)

			c.SetCookie(
				cookieName,
				signedUserID,
				3600*24*30,
				"/",
				"",
				false,
				true,
			)

			c.Set(cookieName, newUserID)
			logger.Infof("created new user with ID: %s", userIDStr)
		} else {
			userIDStr := extractUserID(cookie)

			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				userID = uuid.New()
				userIDStr = userID.String()
				signedUserID := signCookie(userIDStr, secretKey)

				c.SetCookie(
					cookieName,
					signedUserID,
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

func signCookie(userID string, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(userID))
	signature := h.Sum(nil)
	return userID + "." + hex.EncodeToString(signature)
}

func isValidCookie(cookie string, key []byte) bool {
	if cookie == "" {
		return false
	}

	parts := splitMax(cookie, ".", 2)
	if len(parts) != 2 {
		return false
	}

	userID, signatureHex := parts[0], parts[1]

	if _, err := uuid.Parse(userID); err != nil {
		return false
	}

	h := hmac.New(sha256.New, key)
	h.Write([]byte(userID))
	expectedSignature := h.Sum(nil)

	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		return false
	}

	return hmac.Equal(signature, expectedSignature)
}

func extractUserID(cookie string) string {
	parts := splitMax(cookie, ".", 2)
	if len(parts) != 2 {
		return ""
	}
	return parts[0]
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

func splitMax(s, sep string, n int) []string {
	if n <= 0 {
		return []string{s}
	}

	parts := strings.SplitN(s, sep, n)
	return parts
}
