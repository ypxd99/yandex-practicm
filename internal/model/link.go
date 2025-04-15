package model

import (
	"github.com/google/uuid"
)

type Link struct {
	ID     string    `bun:",pk" json:"id"`
	Link   string    `bun:",notnull" json:"link"`
	UserID uuid.UUID `bun:",notnull" json:"user_id"`
	//TimeCreated time.Time `bun:",default:now()" json:"time_created"`
}
