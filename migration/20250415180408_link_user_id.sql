-- +goose Up
-- +goose StatementBegin
ALTER TABLE shortener.links
ADD COLUMN IF NOT EXISTS user_id UUID DEFAULT '00000000-0000-0000-0000-000000000000' NOT NULL;

CREATE INDEX IF NOT EXISTS idx_user_id ON shortener.links (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS shortener.idx_user_id;
ALTER TABLE shortener.links DROP COLUMN IF EXISTS user_id;
-- +goose StatementEnd 