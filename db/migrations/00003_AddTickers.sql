-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE tickers (
  id VARCHAR NOT NULL CHECK (length(id) < 10), -- ex LTC,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  deleted_at timestamptz,
  name VARCHAR NOT NULL CHECK (length(name) < 50), -- ex Litecoin
  unit VARCHAR CHECK (length(curency_symbol) < 10), -- ex Ł
  marketcap BIGINT CHECK (marketcap >= 0),
  volume_24h BIGINT CHECK (volume_24h >= 0),
  max_supply BIGINT CHECK (max_supply >= 0),
  current_supply BIGINT CHECK (current_supply >= 0),
  logo_url VARCHAR CHECK (length(logo_url) < 255),
  website VARCHAR CHECK (length(website) < 255),
  price_usd REAL NOT NULL CHECK (price >= 0),
  percent_change_1h REAL CHECK (percent_change_1h >= 0),
  percent_change_24h REAL CHECK (percent_change_24h >= 0),
  percent_change_7d REAL CHECK (percent_change_7d >= 0),
  coinmarketcap_id VARCHAR NOT NULL UNIQUE CHECK (length(coinmarketcap_id) < 50),
  PRIMARY KEY (id)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE tickers;
