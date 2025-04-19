-- +goose Up
-- +goose StatementBegin
INSERT INTO job_queue(name, created_at, updated_at, scheduled_started_at, params)
VALUES ('prepare.data', NOW(), NOW(), round(EXTRACT(EPOCH FROM NOW())), '{}');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
