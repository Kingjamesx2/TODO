-- Filename: migrations/000001_add_todo_indexes.up.sql
CREATE INDEX IF NOT EXISTS todo_name_idx ON todo USING GIN(to_tsvector('simple', name));
CREATE INDEX IF NOT EXISTS todo_task_idx ON todo USING GIN(to_tsvector('simple', task));