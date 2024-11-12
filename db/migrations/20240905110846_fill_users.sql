-- +goose Up
-- +goose StatementBegin
INSERT INTO myusers (id, login, password, role, token) VALUES
(1, 'cmd@cmd.ru', '$2a$10$42ivYqjjpyKLmg8hnV7XROYJHOlhyW4lpndm7CCjee/VVQbLkYuz6', 'admin', ''),
(2, 'cmd@cmd.org', '$2a$10$42ivYqjjpyKLmg8hnV7XROYJHOlhyW4lpndm7CCjee/VVQbLkYuz6', 'regular',''),
(3, 'cmd@cmd.com', '$2a$10$42ivYqjjpyKLmg8hnV7XROYJHOlhyW4lpndm7CCjee/VVQbLkYuz6', 'regular','');
		

SELECT setval(pg_get_serial_sequence('myusers', 'id'), max(id)) FROM myusers;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DELETE FROM myusers;
-- +goose StatementEnd