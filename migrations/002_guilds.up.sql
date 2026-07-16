-- 002_guilds.up.sql

CREATE TABLE guilds (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(24) UNIQUE NOT NULL,
    tag         VARCHAR(6) UNIQUE NOT NULL,
    owner_id    BIGINT REFERENCES players(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    level       INT NOT NULL DEFAULT 1,
    description TEXT NOT NULL DEFAULT '',
    member_count INT NOT NULL DEFAULT 1
);

CREATE TABLE guild_members (
    guild_id   BIGINT REFERENCES guilds(id) ON DELETE CASCADE,
    player_id  BIGINT REFERENCES players(id) ON DELETE CASCADE,
    role       VARCHAR(16) NOT NULL DEFAULT 'member',
    joined_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (guild_id, player_id)
);

CREATE TABLE guild_chat (
    id         BIGSERIAL PRIMARY KEY,
    guild_id   BIGINT REFERENCES guilds(id) ON DELETE CASCADE,
    player_id  BIGINT REFERENCES players(id),
    message    TEXT NOT NULL,
    sent_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_guild_chat_guild ON guild_chat(guild_id, sent_at);
