package model

type Link struct {
	ID   string `bun:",pk" json:"id"`
	Link string `bun:",notnull" json:"link"`
	//TimeCreated time.Time `bun:",default:now()" json:"time_created"`
}
