package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ypxd99/yandex-practicm/util"
)

type Response struct {
	ErrorText string      `json:"error_text"`
	HasError  bool        `json:"has_error"`
	Resp      interface{} `json:"resp"`
}

func responseTextPlain(c *gin.Context,
	statusCode int,
	err error,
	message []byte) {

	if err != nil {
		util.GetLogger().Error(err)
		c.String(statusCode, err.Error())
		return
	}

	c.Data(statusCode, "text/plain; charset=utf-8", message)
}

func response(c *gin.Context,
	statusCode int,
	err error,
	message interface{}) {
	//resp := Response{
	//	Resp: message,
	//}

	if err != nil {
		//resp.ErrorText = err.Error()
		//resp.HasError = true
		util.GetLogger().Error(err)
	}

	c.AbortWithStatusJSON(statusCode, message)
	//c.AbortWithStatusJSON(statusCode, resp)
}
