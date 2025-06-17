package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ypxd99/yandex-practicm/util"
)

func responseTextPlain(c *gin.Context, statusCode int, err error, message []byte) {
	if err != nil {
		util.GetLogger().Error(err)
		c.String(statusCode, err.Error())
		return
	}

	c.Data(statusCode, "text/plain; charset=utf-8", message)
}

func response(c *gin.Context, statusCode int, err error, message interface{}) {
	if err != nil {
		util.GetLogger().Error(err)
	}

	c.AbortWithStatusJSON(statusCode, message)
}
