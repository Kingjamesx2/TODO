--Filename: migrations/000001_create_todo_table.up.sql

CREATE TABLE IF NOT EXISTS todo (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    task text NOT NULL,
    version int NOT NULL DEFAULT 1

);