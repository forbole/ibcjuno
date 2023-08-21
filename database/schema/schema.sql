/* ---- TOKENS ---- */

CREATE TABLE token
(
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE token_unit
(
    token_name TEXT NOT NULL REFERENCES token (name),
    symbol     TEXT NOT NULL,
    denom      TEXT NOT NULL UNIQUE,
    exponent   BIGINT  NOT NULL,
    price_id   TEXT
);

/* ---- TOKEN PRICES ---- */

CREATE TABLE token_price
(
    id                       SERIAL                      NOT NULL PRIMARY KEY,
    price_id                 TEXT                        NOT NULL UNIQUE,
    name                     TEXT                        NOT NULL,
    image                    TEXT,
    price                    DECIMAL                     NOT NULL,
    market_cap               TEXT                        NOT NULL,
    market_cap_rank          BIGINT                      NOT NULL,
    fully_diluted_valuation  DECIMAL                     NOT NULL,
    total_volume             DECIMAL                     NOT NULL,
    high_24h                 DECIMAL                     NOT NULL,
    low_24h                  DECIMAL                     NOT NULL,
    circulating_supply       DECIMAL                     NOT NULL,
    total_supply             DECIMAL                     NOT NULL,
    max_supply               DECIMAL                     NOT NULL,
    ath                      DECIMAL                     NOT NULL,
    atl                      DECIMAL                     NOT NULL,
    timestamp                TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE INDEX token_price_timestamp_index ON token_price (timestamp);

/* ---- IBC TOKEN ---- */

CREATE TABLE token_ibc
(
  origin_denom   TEXT     NOT NULL,
  origin_chain_price_id   TEXT     NOT NULL,
  target_denom   TEXT     NOT NULL,
  target_chain_price_id   TEXT     NOT NULL,
  trade_url      TEXT     NOT NULL,
  timestamp      TIMESTAMP WITHOUT TIME ZONE NOT NULL 
);

