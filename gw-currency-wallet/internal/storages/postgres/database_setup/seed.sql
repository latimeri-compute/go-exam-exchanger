CREATE DATABASE wallets;

\c wallets

CREATE TABLE IF NOT EXISTS wallets(
    id SERIAL NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    usd_balance bigint NOT NULL DEFAULT 0,
    eur_balance bigint NOT NULL DEFAULT 0,
    rub_balance bigint NOT NULL DEFAULT 0,
    PRIMARY KEY(id)
);
CREATE INDEX IF NOT EXISTS idx_wallets_deleted_at ON wallets USING btree (deleted_at);


CREATE TABLE IF NOT EXISTS users(
    id SERIAL NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    email varchar(255) NOT NULL,
    password_hash bytea NOT NULL,
    wallet_id bigint,
    username varchar(255) NOT NULL,
    PRIMARY KEY(id),
    CONSTRAINT fk_users_wallet FOREIGN key(wallet_id) REFERENCES wallets(id)
);
CREATE UNIQUE INDEX IF NOT EXISTS uni_users_email ON users USING btree (email);
CREATE UNIQUE INDEX IF NOT EXISTS uni_users_username ON users USING btree (username);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users USING btree (deleted_at);
CREATE UNIQUE INDEX IF NOT EXISTS uni_users_username ON users USING btree (username);
