--filename: 000002_add_todo_check_constraint.down.sql

ALTER TABLE todo DROP CONSTRAINT IF EXISTS mode_length_check;
