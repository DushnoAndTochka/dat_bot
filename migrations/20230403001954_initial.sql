-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users
(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tg_id       BIGINT UNIQUE NOT NULL,
    name        VARCHAR DEFAULT FALSE NOT NULL
);

CREATE INDEX ON users (tg_id);

CREATE TABLE problems
(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR DEFAULT FALSE NOT NULL,
    source      VARCHAR DEFAULT FALSE NOT NULL,
    status      VARCHAR DEFAULT FALSE NOT NULL
);

CREATE INDEX ON problems (status);
CREATE INDEX ON problems (name);

CREATE TABLE suggestions
(
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id        UUID REFERENCES users (id) NOT NULL,
    problem_id     UUID REFERENCES problems (id) NOT NULL
);
CREATE INDEX ON suggestions (user_id);
CREATE INDEX ON suggestions (problem_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE suggestions;
DROP TABLE projects;
DROP TABLE users;
-- +goose StatementEnd
