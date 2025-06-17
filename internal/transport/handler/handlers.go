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
// Принимает реализацию интерфейса LinkService.
// Возвращает инициализированный обработчик.
func InitHandler(service service.LinkService) *Handler {
	return &Handler{service: service}
}

// InitRoutes настраивает все HTTP-маршруты для сервиса сокращения URL.
// Настраивает middleware, аутентификацию и все API-эндпоинты.
// Маршруты включают:
// - Эндпоинты отладки для профилирования (/debug/pprof/*)
// - Эндпоинты метрик и проверки работоспособности (/metrics, /health)
// - Эндпоинты сокращения URL (/, /api/shorten)
// - Эндпоинты управления URL пользователя (/api/user/urls)
// Принимает экземпляр gin.Engine для настройки маршрутов.
func (h *Handler) InitRoutes(r *gin.Engine) {
	// Настройка эндпоинтов профилирования
	r.GET("/debug/pprof/", gin.WrapF(http.HandlerFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.DefaultServeMux.ServeHTTP(w, r)
	}))))
	r.GET("/debug/pprof/:profile", gin.WrapF(http.HandlerFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.DefaultServeMux.ServeHTTP(w, r)
	}))))

	// Настройка эндпоинтов метрик и проверки работоспособности
	util.GetMetricsRoute(r)
	util.GetHealthcheckRoute(r)
	util.GetRouteList(r)

	// Настройка middleware
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.GzipMiddleware())
	r.Use(middleware.AuthMiddleware())

	// Настройка основных эндпоинтов
	r.POST("/", h.shorterLink)
	r.GET("/:id", h.getLinkByID)
	r.GET("/ping", h.getStorageStatus)

	// Настройка API эндпоинтов
	rAPI := r.Group("/api")
	rAPI.POST("/shorten", h.shorten)
	rAPI.POST("/shorten/batch", h.batchShorten)

	// Настройка эндпоинтов для работы с URL пользователя
	userAPI := rAPI.Group("/user")
	userAPI.Use(middleware.RequireAuth())
	userAPI.GET("/urls", h.getUserURLs)
	userAPI.DELETE("/urls", h.deleteURLs)
}
