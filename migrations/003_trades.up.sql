-- 003_trades.up.sql

CREATE TABLE trades (
    id          BIGSERIAL PRIMARY KEY,
    player_a_id BIGINT REFERENCES players(id),
    player_b_id BIGINT REFERENCES players(id),
    status      VARCHAR(16) NOT NULL DEFAULT 'pending',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);

CREATE TABLE trade_items (
    trade_id    BIGINT REFERENCES trades(id) ON DELETE CASCADE,
    player_id   BIGINT REFERENCES players(id),
    item_id     INT NOT NULL,
    quantity    INT NOT NULL DEFAULT 1,
    PRIMARY KEY (trade_id, player_id, item_id)
);
