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
  denom TEXT NOT NULL REFERENCES token_unit (denom),
  src_chain   TEXT NOT NULL,
  dst_chain   TEXT NOT NULL,
  channel     TEXT NOT NULL,
  ibc_denom   TEXT NOT NULL UNIQUE,
  CONSTRAINT unique_token_ibc_denom UNIQUE (denom, ibc_denom)
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
