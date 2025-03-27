-- +goose Up
-- +goose StatementBegin
CREATE TABLE prepared_report (
                        id BIGINT,
                        created_at TIMESTAMP WITH TIME ZONE NOT NULL,
                        updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
                        timestamp TIMESTAMP WITH TIME ZONE,
                        source VARCHAR(2048),
                        shortened VARCHAR(256),
                        user_email VARCHAR(72),
                        user_fullname VARCHAR(384),
                        PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE prepared_report_2025_03 PARTITION OF prepared_report
    FOR VALUES FROM ('2025-03-01') TO ('2025-04-01');
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION create_prepared_report_partition()
RETURNS TRIGGER AS $$
BEGIN
    DECLARE
partition_name TEXT;
        start_date TIMESTAMP WITH TIME ZONE;
        end_date TIMESTAMP WITH TIME ZONE;
BEGIN
        start_date := date_trunc('month', NEW.created_at);
        end_date := start_date + INTERVAL '1 month';
        partition_name := 'prepared_report_' || to_char(start_date, 'YYYY_MM');

        -- Проверяем, существует ли партиция
        IF NOT EXISTS (
            SELECT 1 FROM pg_class c
            JOIN pg_namespace n ON n.oid = c.relnamespace
            WHERE c.relname = partition_name AND n.nspname = 'public'
        ) THEN
            EXECUTE format(
                'CREATE TABLE %I PARTITION OF prepared_report_
                FOR VALUES FROM (%L) TO (%L);',
                partition_name, start_date, end_date
            );
END IF;
RETURN NEW;
END;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER prepared_report_partition_trigger
    BEFORE INSERT ON prepared_report
    FOR EACH ROW
    EXECUTE FUNCTION create_prepared_report_partition();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE prepared_report CASCADE;
DROP FUNCTION IF EXISTS create_prepared_report_partition;
-- +goose StatementEnd
