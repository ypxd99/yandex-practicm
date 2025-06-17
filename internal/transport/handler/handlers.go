package handler

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/gin-gonic/gin"
	"github.com/ypxd99/yandex-practicm/internal/repository/middleware"
	"github.com/ypxd99/yandex-practicm/internal/service"
	"github.com/ypxd99/yandex-practicm/util"
)

// Handler представляет HTTP-обработчик сервиса сокращения URL.
// Содержит бизнес-логику сервиса и предоставляет HTTP-эндпоинты
// для операций сокращения URL.
type Handler struct {
	service service.LinkService
}

// InitHandler создает и возвращает новый экземпляр Handler с предоставленным сервисом.
// Инициализирует обработчик с заданной реализацией LinkService.
func InitHandler(service service.LinkService) *Handler {
	return &Handler{service: service}
}

// InitRoutes настраивает все HTTP-маршруты для сервиса сокращения URL.
// Настраивает middleware, аутентификацию и все API-эндпоинты.
// Маршруты включают:
// - Эндпоинты отладки для профилирования
// - Эндпоинты метрик и проверки работоспособности
// - Эндпоинты сокращения URL
// - Эндпоинты управления URL пользователя
func (h *Handler) InitRoutes(r *gin.Engine) {
	r.GET("/debug/pprof/", gin.WrapF(http.HandlerFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.DefaultServeMux.ServeHTTP(w, r)
	}))))
	r.GET("/debug/pprof/:profile", gin.WrapF(http.HandlerFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.DefaultServeMux.ServeHTTP(w, r)
	}))))

	util.GetMetricsRoute(r)
	util.GetHealthcheckRoute(r)
	util.GetRouteList(r)

	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.GzipMiddleware())
	r.Use(middleware.AuthMiddleware())

	r.POST("/", h.shorterLink)
	r.GET("/:id", h.getLinkByID)
	r.GET("/ping", h.getStorageStatus)

	rAPI := r.Group("/api")
	rAPI.POST("/shorten", h.shorten)
	rAPI.POST("/shorten/batch", h.batchShorten)

	userAPI := rAPI.Group("/user")
	userAPI.Use(middleware.RequireAuth())
	userAPI.GET("/urls", h.getUserURLs)
	userAPI.DELETE("/urls", h.deleteURLs)
}
