-- +goose Up
-- +goose StatementBegin
CREATE TABLE solutions
(
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name           VARCHAR DEFAULT FALSE NOT NULL,
    problem_id     UUID REFERENCES problems (id) NOT NULL
    is_solved      BOOL BOOLEAN NOT NULL DEFAULT FALSE;
);
CREATE INDEX ON solutions (problem_id);
CREATE INDEX ON solutions (is_solved);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE solutions;
-- +goose StatementEnd
