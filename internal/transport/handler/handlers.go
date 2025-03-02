package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ypxd99/yandex-practicm/internal/service"
	"github.com/ypxd99/yandex-practicm/util"
)

type Handler struct {
	service service.LinkService
}

func InitHandler(service service.LinkService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes(r *gin.Engine) {
	util.GetMetricsRoute(r)
	util.GetHealthcheckRoute(r)
	util.GetRouteList(r)

	r.POST("/", h.shorterLink)
	r.GET("/:id", h.getLinkByID)
}
