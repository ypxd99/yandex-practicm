package model

// ShortenRequest представляет запрос на сокращение URL.
// Используется в API для получения URL, который нужно сократить.
type ShortenRequest struct {
	// URL оригинальный URL, который нужно сократить
	URL string `json:"url"`
}

// ShortenResponse представляет ответ на запрос сокращения URL.
// Содержит сокращенную версию URL.
type ShortenResponse struct {
	// Result сокращенный URL
	Result string `json:"result"`
}

// BatchRequest представляет запрос на пакетное сокращение URL.
// Используется для массового сокращения ссылок.
type BatchRequest struct {
	// CorrelationID идентификатор корреляции для сопоставления запроса и ответа
	CorrelationID string `json:"correlation_id"`
	// OriginalURL оригинальный URL, который нужно сократить
	OriginalURL string `json:"original_url"`
}

// BatchResponse представляет ответ на пакетный запрос сокращения URL.
// Содержит сокращенные версии URL с их идентификаторами корреляции.
type BatchResponse struct {
	// CorrelationID идентификатор корреляции из запроса
	CorrelationID string `json:"correlation_id"`
	// ShortURL сокращенный URL
	ShortURL string `json:"short_url"`
}

// UserURLResponse представляет информацию о сокращенной ссылке пользователя.
// Используется для отображения списка сокращенных ссылок пользователя.
type UserURLResponse struct {
	// ShortURL сокращенный URL
	ShortURL string `json:"short_url"`
	// OriginalURL оригинальный URL
	OriginalURL string `json:"original_url"`
}

// DeleteRequest представляет запрос на удаление сокращенных ссылок.
// Содержит список идентификаторов ссылок для удаления.
type DeleteRequest []string
