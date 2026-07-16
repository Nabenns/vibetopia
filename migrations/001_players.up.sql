-- 001_players.up.sql
-- Core player table for VIBETOPIA GTPS

CREATE TABLE players (
    id          BIGSERIAL PRIMARY KEY,
    growid      VARCHAR(24) UNIQUE NOT NULL,
    password    VARCHAR(255) NOT NULL,
    display     VARCHAR(32) NOT NULL DEFAULT '',
    role        VARCHAR(16) NOT NULL DEFAULT 'player',
    gems        INT NOT NULL DEFAULT 0,
    level       INT NOT NULL DEFAULT 1,
    xp          BIGINT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_login  TIMESTAMPTZ,
    last_ip     INET,
    is_banned   BOOLEAN NOT NULL DEFAULT FALSE,
    account_notes TEXT NOT NULL DEFAULT '[]',
    mac_addr    VARCHAR(32) DEFAULT '',
    device_id   VARCHAR(64) DEFAULT ''
);

CREATE INDEX idx_players_growid ON players(growid);

CREATE TABLE player_inventory (
    player_id   BIGINT REFERENCES players(id) ON DELETE CASCADE,
    item_id     INT NOT NULL,
    quantity    INT NOT NULL DEFAULT 1,
    PRIMARY KEY (player_id, item_id)
);

CREATE TABLE player_bank (
    player_id   BIGINT REFERENCES players(id) ON DELETE CASCADE,
    item_id     INT NOT NULL,
    quantity    INT NOT NULL DEFAULT 1,
    deposited_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (player_id, item_id)
);

CREATE TABLE player_favorites (
    player_id   BIGINT REFERENCES players(id) ON DELETE CASCADE,
    item_id     INT NOT NULL,
    sort_order  INT NOT NULL DEFAULT 0,
    PRIMARY KEY (player_id, item_id)
);
