-- +goose Up
-- +goose StatementBegin
ALTER TABLE shortener.links ADD COLUMN is_deleted BOOLEAN DEFAULT FALSE NOT NULL;
CREATE INDEX idx_links_is_deleted ON shortener.links (is_deleted);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS shortener.idx_links_is_deleted;
ALTER TABLE shortener.links DROP COLUMN is_deleted;
-- +goose StatementEnd 