-- Filename: migrations/000001_add_todo_indexes.down.sql

DROP INDEX IF EXISTS todo_name_idx;
DROP INDEX IF EXISTS todo_task_idx;