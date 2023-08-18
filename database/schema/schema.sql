/* ---- TOKENS ---- */

CREATE TABLE token
(
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE token_unit
(
    token_name TEXT NOT NULL REFERENCES token (name),
    denom      TEXT NOT NULL UNIQUE,
    exponent   INT  NOT NULL,
    price_id   TEXT
);

/* ---- TOKEN PRICES ---- */

CREATE TABLE token_price
(
    id         SERIAL                      NOT NULL PRIMARY KEY,
    unit_name  TEXT                        NOT NULL REFERENCES token_unit (denom) UNIQUE,
    price      DECIMAL                     NOT NULL,
    market_cap BIGINT                      NOT NULL,
    timestamp  TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE INDEX token_price_timestamp_index ON token_price (timestamp);

/* ---- IBC TOKEN ---- */

CREATE TABLE token_ibc
(
  origin_denom   TEXT     NOT NULL,
  origin_chain   TEXT     NOT NULL,
  target_denom   TEXT     NOT NULL,
  target_chain   TEXT     NOT NULL,
  is_stale       TEXT     NOT NULL,
  trade_url      TEXT     NOT NULL,
  timestamp      TIMESTAMP WITHOUT TIME ZONE NOT NULL 
);

