-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE message (
    id bigint,
    from_user_id bigint,
    to_user_id bigint,
    text text NOT NULL,
    created_at timestamp,
    updated_at timestamp,
    PRIMARY KEY (id)
);
SELECT create_distributed_table('message',   'id');
CREATE SEQUENCE serial_message_id START 1;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE message;
DROP SEQUENCE serial_message_id;
