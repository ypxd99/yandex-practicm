package middleware

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ypxd99/yandex-practicm/util"
)

// TrustedSubnetMiddleware создает middleware для проверки доверенной подсети.
// Проверяет, что IP-адрес клиента входит в доверенную подсеть.
// Возвращает gin.HandlerFunc.
func TrustedSubnetMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := util.GetConfig()
		trustedSubnet := cfg.Server.TrustedSubnet

		if trustedSubnet == "" {
			c.Next()
			return
		}

		clientIP := c.GetHeader("X-Real-IP")
		if clientIP == "" {
			clientIP = c.ClientIP()
		}

		_, subnet, err := net.ParseCIDR(trustedSubnet)
		if err != nil {
			util.GetLogger().Errorf("invalid trusted subnet format: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ip := net.ParseIP(clientIP)
		if ip == nil {
			util.GetLogger().Errorf("invalid client IP: %s", clientIP)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if !subnet.Contains(ip) {
			util.GetLogger().Warnf("access denied for IP %s, not in trusted subnet %s", clientIP, trustedSubnet)
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		util.GetLogger().Infof("access granted for IP %s from trusted subnet %s", clientIP, trustedSubnet)
		c.Next()
	}
}
