-- 005_quests.up.sql

CREATE TABLE daily_quests (
    id          BIGSERIAL PRIMARY KEY,
    quest_type  VARCHAR(24) NOT NULL,
    target      INT NOT NULL DEFAULT 0,
    reward_item INT NOT NULL DEFAULT 0,
    reward_qty  INT NOT NULL DEFAULT 1,
    date        DATE NOT NULL,
    UNIQUE (quest_type, date)
);

CREATE TABLE player_quests (
    player_id   BIGINT REFERENCES players(id) ON DELETE CASCADE,
    quest_id    BIGINT REFERENCES daily_quests(id) ON DELETE CASCADE,
    progress    INT NOT NULL DEFAULT 0,
    completed   BOOLEAN NOT NULL DEFAULT FALSE,
    claimed     BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (player_id, quest_id)
);

CREATE TABLE growtokens (
    player_id   BIGINT REFERENCES players(id) ON DELETE CASCADE,
    source      VARCHAR(24) NOT NULL,
    earned_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (player_id, source)
);
