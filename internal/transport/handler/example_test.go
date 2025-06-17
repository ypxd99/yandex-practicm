package handler_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/ypxd99/yandex-practicm/internal/mocks"
	"github.com/ypxd99/yandex-practicm/internal/model"
	"github.com/ypxd99/yandex-practicm/internal/transport/handler"
	"github.com/ypxd99/yandex-practicm/util"
	"gopkg.in/yaml.v3"
)

var (
	cfgPath = "configuration/config.yaml"
)

func parseConfig(st interface{}, cfgPath string) {
	f, err := os.Open(cfgPath)
	if err != nil {
		log.Fatal(errors.WithMessage(err, "error occurred while opening cfg file"))
	}

	fi, err := f.Stat()
	if err != nil {
		log.Fatal(errors.WithMessage(err, "error occurred while getting file stats"))
	}

	data := make([]byte, fi.Size())
	_, err = f.Read(data)
	if err != nil {
		log.Fatal(errors.WithMessage(err, "error occurred while reading data"))
	}

	err = yaml.Unmarshal(data, st)
	if err != nil {
		log.Fatal(errors.WithMessage(err, "error occurred while unmashaling data"))
	}
}

func decode(str string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", errors.WithMessage(err, "error occurred while decoding string(base64)")
	}
	res, err := util.GetRSA().Decrypt(data)
	if err != nil {
		return "", errors.WithMessage(err, "error occurred while decoding string(RSA)")
	}
	return string(res), err
}

func decodeCFG(cfg *util.Config) error {
	var err error
	cfg.Postgres.Address, err = decode(cfg.Postgres.Address)
	if err != nil {
		return errors.WithMessage(err, "error occurred while decode address")
	}
	cfg.Postgres.User, err = decode(cfg.Postgres.User)
	if err != nil {
		return errors.WithMessage(err, "error occurred while decode user")
	}
	cfg.Postgres.Password, err = decode(cfg.Postgres.Password)
	if err != nil {
		return errors.WithMessage(err, "error occurred while decode password")
	}

	return nil
}

func init() {
	var (
		conf util.Config
	)
	parseConfig(&conf, cfgPath)
	if conf.UseDecode {
		decodeCFG(&conf)
	}

	conf.Server.ServerAddress = fmt.Sprintf("%s:%d", conf.Server.Address, conf.Server.Port)
	conf.Server.BaseURL = fmt.Sprintf("http://%s:%d", conf.Server.Address, conf.Server.Port)
	conf.Postgres.UsePostgres = false

	util.InitLogger(conf.Logger)
}

// Example_shorten демонстрирует использование эндпоинта сокращения URL.
func Example_shorten() {
	// Создаем тестовый роутер
	router := gin.Default()

	// Инициализируем обработчик с мок-сервисом
	mockService := &mocks.MockLinkService{}
	mockService.On("ShorterLink", mock.Anything, "https://example.com/very/long/url", mock.Anything).
		Return("http://localhost:8080/abc123", nil)
	h := handler.InitHandler(mockService)
	h.InitRoutes(router)

	// Создаем тестовый запрос
	req := model.ShortenRequest{
		URL: "https://example.com/very/long/url",
	}
	body, _ := json.Marshal(req)

	// Создаем тестовый HTTP запрос
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Cookie", fmt.Sprintf("user_id=%s", uuid.New().String()))

	// Отправляем запрос
	router.ServeHTTP(w, httpReq)

	// Выводим ответ
	fmt.Println(w.Code)
	fmt.Println(w.Body.String())

	// Output:
	// 201
	// {"result":"http://localhost:8080/abc123"}
}

// Example_batchShorten демонстрирует использование эндпоинта пакетного сокращения URL.
func Example_batchShorten() {
	// Создаем тестовый роутер
	router := gin.Default()

	// Инициализируем обработчик с мок-сервисом
	mockService := &mocks.MockLinkService{}
	batchReq := []model.BatchRequest{
		{
			CorrelationID: "1",
			OriginalURL:   "https://example.com/url1",
		},
		{
			CorrelationID: "2",
			OriginalURL:   "https://example.com/url2",
		},
	}
	batchResp := []model.BatchResponse{
		{
			CorrelationID: "1",
			ShortURL:      "http://localhost:8080/abc123",
		},
		{
			CorrelationID: "2",
			ShortURL:      "http://localhost:8080/def456",
		},
	}
	mockService.On("BatchShorten", mock.Anything, batchReq, mock.Anything).
		Return(batchResp, nil)
	h := handler.InitHandler(mockService)
	h.InitRoutes(router)

	// Создаем тестовый запрос
	body, _ := json.Marshal(batchReq)

	// Создаем тестовый HTTP запрос
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/shorten/batch", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Cookie", fmt.Sprintf("user_id=%s", uuid.New().String()))

	// Отправляем запрос
	router.ServeHTTP(w, httpReq)

	// Выводим ответ
	fmt.Println(w.Code)
	fmt.Println(w.Body.String())

	// Output:
	// 201
	// [{"correlation_id":"1","short_url":"http://localhost:8080/abc123"},{"correlation_id":"2","short_url":"http://localhost:8080/def456"}]
}

// Example_getUserURLs демонстрирует использование эндпоинта получения URL пользователя.
func Example_getUserURLs() {
	// Создаем тестовый роутер
	router := gin.Default()

	// Инициализируем обработчик с мок-сервисом
	mockService := &mocks.MockLinkService{}
	urls := []model.UserURLResponse{
		{
			ShortURL:    "http://localhost:8080/abc123",
			OriginalURL: "https://example.com/url1",
		},
		{
			ShortURL:    "http://localhost:8080/def456",
			OriginalURL: "https://example.com/url2",
		},
	}
	mockService.On("GetUserURLs", mock.Anything, mock.Anything).
		Return(urls, nil)
	h := handler.InitHandler(mockService)
	h.InitRoutes(router)

	// Создаем тестовый HTTP запрос
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/user/urls", nil)
	httpReq.Header.Set("Cookie", fmt.Sprintf("user_id=%s", uuid.New().String()))

	// Отправляем запрос
	router.ServeHTTP(w, httpReq)

	// Выводим ответ
	fmt.Println(w.Code)
	fmt.Println(w.Body.String())

	// Output:
	// 200
	// [{"short_url":"http://localhost:8080/abc123","original_url":"https://example.com/url1"},{"short_url":"http://localhost:8080/def456","original_url":"https://example.com/url2"}]
}

// Example_deleteURLs демонстрирует использование эндпоинта удаления URL.
func Example_deleteURLs() {
	// Создаем тестовый роутер
	router := gin.Default()

	// Инициализируем обработчик с мок-сервисом
	mockService := &mocks.MockLinkService{}
	ids := []string{"abc123", "def456"}
	mockService.On("DeleteURLs", mock.Anything, ids, mock.Anything).
		Return(2, nil)
	h := handler.InitHandler(mockService)
	h.InitRoutes(router)

	// Создаем тестовый запрос
	req := model.DeleteRequest(ids)
	body, _ := json.Marshal(req)

	// Создаем тестовый HTTP запрос
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("DELETE", "/api/user/urls", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Cookie", fmt.Sprintf("user_id=%s", uuid.New().String()))

	// Отправляем запрос
	router.ServeHTTP(w, httpReq)

	// Выводим ответ
	fmt.Println(w.Code)

	// Output:
	// 202
}
