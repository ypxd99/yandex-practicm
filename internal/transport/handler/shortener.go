package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ypxd99/yandex-practicm/internal/model"
	"github.com/ypxd99/yandex-practicm/internal/service"
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
		if errors.Is(err, service.ErrURLExist) {
			responseTextPlain(c, http.StatusConflict, nil, []byte(resp))
		}
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

	//err = c.ShouldBindJSON(&req)
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

	resp, err := h.service.ShorterLink(c.Request.Context(), req.URL)
	if err != nil {
		if errors.Is(err, service.ErrURLExist) {
			response(c, http.StatusConflict, nil, model.ShortenResponse{Result: resp})
		}
		response(c, http.StatusInternalServerError, err, model.ShortenResponse{Result: ""})
		return
	}

	response(c, http.StatusCreated, nil, model.ShortenResponse{Result: resp})
}

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

    responses, err := h.service.BatchShorten(c.Request.Context(), req)
    if err != nil {
        response(c, http.StatusInternalServerError, err, nil)
        return
    }

    response(c, http.StatusCreated, nil, responses)
}