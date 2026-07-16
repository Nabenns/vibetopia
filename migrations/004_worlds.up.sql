-- 004_worlds.up.sql

CREATE TABLE worlds (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(24) UNIQUE NOT NULL,
    owner_id    BIGINT REFERENCES players(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    width       INT NOT NULL DEFAULT 100,
    height      INT NOT NULL DEFAULT 60,
    is_locked   BOOLEAN NOT NULL DEFAULT FALSE,
    visits      INT NOT NULL DEFAULT 0,
    wotd_score  INT NOT NULL DEFAULT 0,
    category    INT NOT NULL DEFAULT 0
);

CREATE TABLE world_tiles (
    world_id   BIGINT REFERENCES worlds(id) ON DELETE CASCADE,
    x          INT NOT NULL,
    y          INT NOT NULL,
    fg_item    INT NOT NULL DEFAULT 0,
    bg_item    INT NOT NULL DEFAULT 0,
    PRIMARY KEY (world_id, x, y)
);

CREATE TABLE world_blocks (
    world_id   BIGINT REFERENCES worlds(id) ON DELETE CASCADE,
    x          INT NOT NULL,
    y          INT NOT NULL,
    item_id    INT NOT NULL,
    label      TEXT DEFAULT '',
    PRIMARY KEY (world_id, x, y)
);
