package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ypxd99/yandex-practicm/internal/model"
	"github.com/ypxd99/yandex-practicm/internal/repository/middleware"
	"github.com/ypxd99/yandex-practicm/internal/service"
	"github.com/ypxd99/yandex-practicm/util"
)

// shorterLink обрабатывает POST-запрос для сокращения URL.
// Принимает URL в теле запроса в виде текста.
// Возвращает сокращенный URL в теле ответа.
// Статусы ответа:
// - 201: URL успешно сокращен
// - 400: Неверный формат запроса
// - 401: Пользователь не авторизован
// - 409: URL уже существует
// - 500: Внутренняя ошибка сервера
func (h *Handler) shorterLink(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		responseTextPlain(c, http.StatusBadRequest, err, nil)
		return
	}
	if len(body) == 0 {
		responseTextPlain(c, http.StatusBadRequest, errors.New("empty data"), nil)
		return
	}
	userID, err := middleware.GetUserID(c)
	if err != nil {
		responseTextPlain(c, http.StatusUnauthorized, err, nil)
		return
	}

	resp, err := h.service.ShorterLink(c.Request.Context(), string(body), userID)
	if err != nil {
		if errors.Is(err, service.ErrURLExist) {
			responseTextPlain(c, http.StatusConflict, nil, []byte(resp))
			return
		}
		responseTextPlain(c, http.StatusInternalServerError, err, nil)
		return
	}

	responseTextPlain(c, http.StatusCreated, nil, []byte(resp))
}

// getLinkByID обрабатывает GET-запрос для получения оригинального URL по его сокращенному идентификатору.
// Принимает идентификатор в параметре пути.
// Выполняет редирект на оригинальный URL.
// Статусы ответа:
// - 307: Редирект на оригинальный URL
// - 400: Неверный формат запроса
// - 410: URL был удален
// - 500: Внутренняя ошибка сервера
func (h *Handler) getLinkByID(c *gin.Context) {
	req := c.Param("id")
	if req == "" {
		responseTextPlain(c, http.StatusBadRequest, errors.New("empty data"), nil)
		return
	}

	resp, err := h.service.FindLink(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrURLDeleted) {
			responseTextPlain(c, http.StatusGone, err, nil)
			return
		}
		responseTextPlain(c, http.StatusBadRequest, err, nil)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, resp)
}

// shorten обрабатывает POST-запрос для сокращения URL через API.
// Принимает JSON с полем "url".
// Возвращает JSON с полем "result", содержащим сокращенный URL.
// Статусы ответа:
// - 201: URL успешно сокращен
// - 400: Неверный формат запроса
// - 401: Пользователь не авторизован
// - 409: URL уже существует
// - 500: Внутренняя ошибка сервера
func (h *Handler) shorten(c *gin.Context) {
	var (
		err error
		req model.ShortenRequest
	)

	// err = c.ShouldBindJSON(&req)
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response(c, http.StatusBadRequest, err, model.ShortenResponse{Result: ""})
		return
	}
	defer c.Request.Body.Close()

	err = json.Unmarshal(body, &req)
	if err != nil || req.URL == "" {
		response(c, http.StatusBadRequest, err, model.ShortenResponse{Result: ""})
		return
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		response(c, http.StatusUnauthorized, err, model.ShortenResponse{Result: ""})
		return
	}

	resp, err := h.service.ShorterLink(c.Request.Context(), req.URL, userID)
	if err != nil {
		if errors.Is(err, service.ErrURLExist) {
			response(c, http.StatusConflict, nil, model.ShortenResponse{Result: resp})
			return
		}
		response(c, http.StatusInternalServerError, err, model.ShortenResponse{Result: ""})
		return
	}

	response(c, http.StatusCreated, nil, model.ShortenResponse{Result: resp})
}

// getStorageStatus обрабатывает GET-запрос для проверки состояния хранилища.
// Возвращает статус доступности хранилища.
// Статусы ответа:
// - 200: Хранилище доступно
// - 500: Хранилище недоступно
func (h *Handler) getStorageStatus(c *gin.Context) {
	status, err := h.service.StorageStatus(c.Request.Context())
	if err != nil {
		response(c, http.StatusInternalServerError, err, nil)
		return
	}
	if !status {
		response(c, http.StatusInternalServerError, errors.New("bad storage status"), nil)
		return
	}

	response(c, http.StatusOK, nil, nil)
}

// batchShorten обрабатывает POST-запрос для пакетного сокращения URL.
// Принимает массив JSON-объектов с полями "correlation_id" и "original_url".
// Возвращает массив JSON-объектов с полями "correlation_id" и "short_url".
// Статусы ответа:
// - 201: URLs успешно сокращены
// - 400: Неверный формат запроса
// - 401: Пользователь не авторизован
// - 500: Внутренняя ошибка сервера
func (h *Handler) batchShorten(c *gin.Context) {
	var (
		err error
		req []model.BatchRequest
	)

	err = c.ShouldBindJSON(&req)
	if err != nil {
		response(c, http.StatusBadRequest, err, nil)
		return
	}
	if len(req) == 0 {
		response(c, http.StatusBadRequest, errors.New("empty batch"), nil)
		return
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		response(c, http.StatusUnauthorized, err, model.ShortenResponse{Result: ""})
		return
	}

	resp, err := h.service.BatchShorten(c.Request.Context(), req, userID)
	if err != nil {
		response(c, http.StatusInternalServerError, err, nil)
		return
	}

	response(c, http.StatusCreated, nil, resp)
}

// getUserURLs обрабатывает GET-запрос для получения списка URL пользователя.
// Возвращает массив JSON-объектов с полями "short_url" и "original_url".
// Статусы ответа:
// - 200: Список URL успешно получен
// - 204: У пользователя нет сохраненных URL
// - 401: Пользователь не авторизован
// - 500: Внутренняя ошибка сервера
func (h *Handler) getUserURLs(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		responseTextPlain(c, http.StatusInternalServerError, err, nil)
		return
	}

	urls, err := h.service.GetUserURLs(c.Request.Context(), userID)
	if err != nil {
		responseTextPlain(c, http.StatusInternalServerError, err, nil)
		return
	}

	if len(urls) == 0 {
		responseTextPlain(c, http.StatusNoContent, errors.New("no content"), nil)
		return
	}

	c.JSON(http.StatusOK, urls)
}

// deleteURLs обрабатывает DELETE-запрос для удаления URL пользователя.
// Принимает массив идентификаторов URL в теле запроса.
// Выполняет мягкое удаление (помечает URL как удаленные).
// Статусы ответа:
// - 202: Запрос на удаление принят
// - 400: Неверный формат запроса
// - 401: Пользователь не авторизован
// - 500: Внутренняя ошибка сервера
func (h *Handler) deleteURLs(c *gin.Context) {
	var (
		err error
		req model.DeleteRequest
	)

	err = c.ShouldBindJSON(&req)
	if err != nil {
		response(c, http.StatusBadRequest, err, nil)
		return
	}

	if len(req) == 0 {
		response(c, http.StatusBadRequest, errors.New("empty url list"), nil)
		return
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		response(c, http.StatusUnauthorized, err, nil)
		return
	}

	count, err := h.service.DeleteURLs(c.Request.Context(), req, userID)
	if err != nil {
		response(c, http.StatusInternalServerError, err, nil)
		return
	}

	if count > 0 {
		util.GetLogger().Infof("successfully marked %d URLs as deleted", count)
	}

	response(c, http.StatusAccepted, nil, nil)
}
