-- +goose Up

CREATE TABLE households (
    id         uuid        PRIMARY KEY DEFAULT uuidv7(),
    name       text        NOT NULL CHECK (char_length(name) BETWEEN 1 AND 100),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TRIGGER households_set_updated_at BEFORE UPDATE ON households
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE household_members (
    id           uuid        PRIMARY KEY DEFAULT uuidv7(),
    household_id uuid        NOT NULL REFERENCES households (id) ON DELETE CASCADE,
    user_id      uuid        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at   timestamptz NOT NULL DEFAULT now(),
    updated_at   timestamptz NOT NULL DEFAULT now()
);

-- One membership per user per household. Also serves "list members of a
-- household", because household_id is the leading column.
CREATE UNIQUE INDEX household_members_household_user_key
    ON household_members (household_id, user_id);

-- Every logged-in request asks "which households is this user in?".
CREATE INDEX household_members_user_id_idx ON household_members (user_id);

CREATE TRIGGER household_members_set_updated_at BEFORE UPDATE ON household_members
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE household_invites (
    id           uuid        PRIMARY KEY DEFAULT uuidv7(),
    household_id uuid        NOT NULL REFERENCES households (id) ON DELETE CASCADE,
    code         text        NOT NULL CHECK (char_length(code) = 8),
    created_by   uuid        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    expires_at   timestamptz NOT NULL,
    created_at   timestamptz NOT NULL DEFAULT now(),
    updated_at   timestamptz NOT NULL DEFAULT now()
);
-- Single use: sign-up deletes the invite in the same transaction that adds
-- the membership. A cleanup job deletes expired, unused invites.

-- Sign-up looks the invite up by code.
CREATE UNIQUE INDEX household_invites_code_key ON household_invites (code);

-- Cascade deletes from households filter by household_id.
CREATE INDEX household_invites_household_id_idx ON household_invites (household_id);

CREATE TRIGGER household_invites_set_updated_at BEFORE UPDATE ON household_invites
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- +goose Down
DROP TABLE household_invites;
DROP TABLE household_members;
DROP TABLE households;
