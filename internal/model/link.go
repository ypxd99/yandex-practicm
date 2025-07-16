package model

import (
	"github.com/google/uuid"
)

// Link представляет структуру для хранения информации о сокращенной ссылке.
// Используется для работы с базой данных и API.
type Link struct {
	// ID уникальный идентификатор сокращенной ссылки
	ID string `json:"id" bun:"id,pk"`
	// Link оригинальный URL, который был сокращен
	Link string `bun:",notnull" json:"link"`
	// UserID идентификатор пользователя, создавшего сокращенную ссылку
	UserID uuid.UUID `bun:",notnull" json:"user_id"`
	// IsDeleted флаг, указывающий, была ли ссылка помечена как удаленная
	IsDeleted bool `bun:",default:false" json:"is_deleted"`
	// TimeCreated time.Time `bun:",default:now()" json:"time_created"`
}
