package handler

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ypxd99/yandex-practicm/internal/model"
)

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

	resp, err := h.service.ShorterLink(c.Request.Context(), string(body))
	if err != nil {
		responseTextPlain(c, http.StatusInternalServerError, err, nil)
		return
	}

	responseTextPlain(c, http.StatusCreated, nil, []byte(resp))
}

func (h *Handler) getLinkByID(c *gin.Context) {
	req := c.Param("id")
	if req == "" {
		responseTextPlain(c, http.StatusBadRequest, errors.New("empty data"), nil)
		return
	}

	resp, err := h.service.FindLink(c.Request.Context(), req)
	if err != nil {
		responseTextPlain(c, http.StatusInternalServerError, err, nil)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, resp)
}

func (h *Handler) shorten(c *gin.Context) {
	var (
		err error
		req model.ShortenRequest
	)

	err = c.ShouldBindJSON(&req)
	if err != nil || req.URL == ""{
		response(c, http.StatusBadRequest, err, model.ShortenResponse{Result: ""})
		return
	}

	resp, err := h.service.ShorterLink(c.Request.Context(), req.URL)
	if err != nil {
		response(c, http.StatusInternalServerError, err, model.ShortenResponse{Result: ""})
		return
	}

	response(c, http.StatusCreated, nil, model.ShortenResponse{Result: resp})
}
