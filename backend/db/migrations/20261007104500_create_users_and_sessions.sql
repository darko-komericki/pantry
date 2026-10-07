-- +goose Up

-- Shared by every table: keeps updated_at correct even if a query forgets it.
-- +goose StatementBegin
CREATE FUNCTION set_updated_at() RETURNS trigger AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TABLE users (
    id            uuid        PRIMARY KEY DEFAULT uuidv7(),
    email         text        NOT NULL CHECK (char_length(email) BETWEEN 3 AND 254),
    password_hash text        NOT NULL,
    display_name  text        NOT NULL CHECK (char_length(display_name) BETWEEN 1 AND 100),
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now()
);

-- Login finds a user by email. lower() makes uniqueness case-insensitive,
-- so Darko@x.com and darko@x.com cannot both register.
CREATE UNIQUE INDEX users_email_lower_key ON users (lower(email));

CREATE TRIGGER users_set_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE sessions (
    id         uuid        PRIMARY KEY DEFAULT uuidv7(),
    user_id    uuid        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token_hash bytea       NOT NULL,  -- SHA-256 of the cookie token, never the token itself
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

-- Every logged-in request finds its session by token hash.
CREATE UNIQUE INDEX sessions_token_hash_key ON sessions (token_hash);

-- "Log out everywhere" and the cascade delete from users filter by user_id.
-- Postgres does not index foreign keys automatically.
CREATE INDEX sessions_user_id_idx ON sessions (user_id);

CREATE TRIGGER sessions_set_updated_at BEFORE UPDATE ON sessions
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- +goose Down
DROP TABLE sessions;
DROP TABLE users;
DROP FUNCTION set_updated_at();
