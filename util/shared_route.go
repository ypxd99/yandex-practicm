package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type UpdateRoutesReq struct {
	ServiceName string `json:"service_name"`
	URL         string `json:"URL"`
}

type route struct {
	Method  string `json:"method"`
	Path    string `json:"path"`
	Handler string `json:"handler"`
}

type RoutesLog struct {
	Routes map[string]bool `yaml:"Routes"`
}

func GetMetricsRoute(r *gin.Engine) {
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

func GetHealthcheckRoute(r *gin.Engine) {
	r.GET("/healthcheck", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
}

func GetRouteList(r *gin.Engine) {
	r.GET("/routes", func(c *gin.Context) {
		resp := make([]route, 0)
		for _, r := range r.Routes() {
			resp = append(resp, route{
				Method:  r.Method,
				Path:    r.Path,
				Handler: r.Handler,
			})
		}

		c.JSON(http.StatusOK, resp)
	})
}
