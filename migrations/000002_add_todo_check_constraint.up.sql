--filename: 000002_add_todo_check_constraint.up.sql

ALTER TABLE todo ADD CONSTRAINT mode_length_check CHECK (array_length(mode, 1) BETWEEN 1 AND 5);
