CREATE TYPE user_role AS ENUM ('user', 'admin');
CREATE TYPE subscription_tier AS ENUM ('free', 'premium');

CREATE TABLE users (
    id                          bigserial primary key,
    first_name                  varchar not null,
    last_name                   varchar not null,
    email                       varchar unique not null,
    password_hash               varchar not null,
    role                        user_role not null default 'user',
    tier                        subscription_tier not null default 'free',
    password_reset_token        varchar,
    password_reset_expires_at   timestamptz,

    CONSTRAINT chk_password_reset_paired
        CHECK (
            (password_reset_token IS NULL) = (password_reset_expires_at IS NULL)
        ),

    created_at                  timestamptz not null default now(),
    updated_at                  timestamptz not null default now()
);

CREATE INDEX idx_users_password_reset_token ON users(password_reset_token)
    WHERE password_reset_token IS NOT NULL;

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_set_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();