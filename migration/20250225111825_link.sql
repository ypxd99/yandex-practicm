-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS shortener.links (
    id VARCHAR(8),
	link varchar(255) NOT NULL,
	time_created timestamptz DEFAULT now() NOT NULL,
	CONSTRAINT links_pkey PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS shortener.links;
-- +goose StatementEnd
