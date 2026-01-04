CREATE TABLE IF NOT EXISTS hourly_data (
    device_id     VARCHAR(50) NOT NULL,
    ts            TIMESTAMPTZ NOT NULL,
    energy     DECIMAL(10,3) NOT NULL,
    total_cost    DECIMAL(15,2) NOT NULL,
    avg_voltage   DECIMAL(10,2),
    avg_current   DECIMAL(10,3),
    avg_power     DECIMAL(10,2),
    min_power     DECIMAL(10,2),
    max_power     DECIMAL(10,2),
    created_at    TIMESTAMPTZ DEFAULT NOW(),

    PRIMARY KEY (device_id, ts)
) PARTITION BY RANGE (ts);

CREATE TABLE IF NOT EXISTS daily_data (
    device_id    VARCHAR(50) NOT NULL,
    day          DATE NOT NULL,
    energy    DECIMAL(10,3) NOT NULL,
    total_cost   DECIMAL(15,2) NOT NULL,
    avg_voltage  DECIMAL(10,2),
    avg_current  DECIMAL(10,3),
    avg_power    DECIMAL(10,2),
    min_power    DECIMAL(10,2),
    max_power    DECIMAL(10,2),
    created_at   TIMESTAMPTZ DEFAULT NOW(),

    PRIMARY KEY (device_id, day)
) PARTITION BY RANGE (day);

CREATE TABLE hourly_data_2025
PARTITION OF hourly_data
FOR VALUES FROM ('2025-01-01') TO ('2026-01-01');

CREATE INDEX idx_hourly_2025_device ON hourly_data_2025(device_id);
CREATE INDEX idx_hourly_2025_ts ON hourly_data_2025(ts);

CREATE TABLE daily_data_2025
PARTITION OF daily_data
FOR VALUES FROM ('2025-01-01') TO ('2026-01-01');

CREATE INDEX idx_daily_2025_device ON daily_data_2025(device_id);
CREATE INDEX idx_daily_2025_day ON daily_data_2025(day);

CREATE TABLE tariffs (
    id BIGSERIAL PRIMARY KEY,

    type_tarrif VARCHAR(20) NOT NULL,
    power_va INTEGER NOT NULL,

    price_per_kwh NUMERIC(10,2),

    effective_from DATE NOT NULL,
    effective_to DATE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS monthly_data (
    device_id    VARCHAR(50) NOT NULL,
    month        DATE NOT NULL, 
    energy       DECIMAL(10,3) NOT NULL,
    total_cost   DECIMAL(15,2) NOT NULL,
    created_at   TIMESTAMPTZ DEFAULT NOW(),

    PRIMARY KEY (device_id, month)
) PARTITION BY RANGE (month);

CREATE TABLE monthly_data_2025
PARTITION OF monthly_data
FOR VALUES FROM ('2025-01-01') TO ('2026-01-01');

CREATE INDEX idx_monthly_2025_device ON monthly_data_2025(device_id);
CREATE INDEX idx_monthly_2025_month ON monthly_data_2025(month);

CREATE TABLE IF NOT EXISTS users_information (
    user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    birth DATE,
    birth_place VARCHAR(100),
    phone VARCHAR(30),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS devices (
  id BIGSERIAL PRIMARY KEY,

  device_id VARCHAR(64) UNIQUE NOT NULL,
  device_name VARCHAR(100),

  device_secret TEXT NOT NULL,

  paired BOOLEAN DEFAULT FALSE,
  paired_at TIMESTAMPTZ,

  created_at TIMESTAMPTZ DEFAULT NOW(),
  last_seen TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS device_users (
  id BIGSERIAL PRIMARY KEY,

  device_id BIGINT REFERENCES devices(id) ON DELETE CASCADE,
  user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,

  role VARCHAR(20) DEFAULT 'owner',

  created_at TIMESTAMPTZ DEFAULT NOW(),

  UNIQUE (device_id, user_id)
);