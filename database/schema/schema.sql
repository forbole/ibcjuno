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

CREATE TABLE token_ibc_denom
(
  denom          TEXT     NOT NULL,
  origin_chain   TEXT     NOT NULL,
  origin_denom   TEXT     NOT NULL,
  origin_type    TEXT     NOT NULL,
  symbol         TEXT     NOT NULL,
  enable         BOOLEAN  NOT NULL,
  path           TEXT     NOT NULL,
  channel        TEXT     NOT NULL,
  counter_party  JSONB    NOT NULL
--   CONSTRAINT unique_token_ibc_denom UNIQUE (denom, symbol)
);

/* ---- TOKEN PRICES ---- */

CREATE TABLE token_price
(
    id         SERIAL                      NOT NULL PRIMARY KEY,
    unit_name  TEXT                        NOT NULL UNIQUE,
    price      DECIMAL                     NOT NULL,
    market_cap BIGINT                      NOT NULL,
    timestamp  TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE INDEX token_price_timestamp_index ON token_price (timestamp);


CREATE TABLE token_ibc_denom_new
(
  denom          TEXT     NOT NULL,
  origin_chain   TEXT     NOT NULL,
  target_denom   TEXT     NOT NULL,
  target_chain    TEXT     NOT NULL,
  is_stale       BOOLEAN     NOT NULL,
  url           TEXT     NOT NULL,
);
