-- +goose Up
-- +goose StatementBegin
ALTER TABLE shortener.links 
ADD CONSTRAINT links_link_unique UNIQUE (link);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE shortener.links 
DROP CONSTRAINT IF EXISTS links_link_unique;
-- +goose StatementEnd