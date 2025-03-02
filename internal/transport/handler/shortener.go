package handler

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
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

	responseTextPlain(c, http.StatusOK, nil, []byte(resp))
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
