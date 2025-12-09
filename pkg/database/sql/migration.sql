CREATE TABLE IF NOT EXISTS hourly_data (
    device_id     VARCHAR(50) NOT NULL,
    ts            TIMESTAMPTZ NOT NULL,
    usage_kwh     DECIMAL(10,3) NOT NULL,
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
    usage_kwh    DECIMAL(10,3) NOT NULL,
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