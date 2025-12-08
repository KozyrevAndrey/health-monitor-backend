-- Initial schema for Health Monitor

-- Targets table
CREATE TABLE IF NOT EXISTS targets (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    config TEXT NOT NULL, -- JSON
    interval INTEGER NOT NULL, -- in nanoseconds
    timeout INTEGER NOT NULL, -- in nanoseconds
    enabled BOOLEAN NOT NULL DEFAULT 1,
    tags TEXT, -- JSON array
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_targets_enabled ON targets(enabled);
CREATE INDEX IF NOT EXISTS idx_targets_type ON targets(type);

-- Check results table
CREATE TABLE IF NOT EXISTS check_results (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    target_id TEXT NOT NULL,
    status TEXT NOT NULL,
    response_time_ms INTEGER NOT NULL,
    status_code INTEGER,
    error TEXT,
    message TEXT,
    metadata TEXT, -- JSON
    checked_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (target_id) REFERENCES targets(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_check_results_target_id ON check_results(target_id);
CREATE INDEX IF NOT EXISTS idx_check_results_checked_at ON check_results(checked_at);
CREATE INDEX IF NOT EXISTS idx_check_results_status ON check_results(status);
CREATE INDEX IF NOT EXISTS idx_check_results_target_checked ON check_results(target_id, checked_at DESC);

-- Incidents table
CREATE TABLE IF NOT EXISTS incidents (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    target_id TEXT NOT NULL,
    target_name TEXT NOT NULL,
    status TEXT NOT NULL,
    started_at TIMESTAMP NOT NULL,
    resolved_at TIMESTAMP,
    duration INTEGER, -- in nanoseconds
    failure_count INTEGER NOT NULL DEFAULT 0,
    last_error TEXT,
    alerts_sent INTEGER NOT NULL DEFAULT 0,
    severity TEXT NOT NULL,
    first_check_result_id INTEGER,
    last_check_result_id INTEGER,
    FOREIGN KEY (target_id) REFERENCES targets(id) ON DELETE CASCADE,
    FOREIGN KEY (first_check_result_id) REFERENCES check_results(id),
    FOREIGN KEY (last_check_result_id) REFERENCES check_results(id)
);

CREATE INDEX IF NOT EXISTS idx_incidents_target_id ON incidents(target_id);
CREATE INDEX IF NOT EXISTS idx_incidents_status ON incidents(status);
CREATE INDEX IF NOT EXISTS idx_incidents_started_at ON incidents(started_at DESC);
CREATE INDEX IF NOT EXISTS idx_incidents_target_status ON incidents(target_id, status);

-- Alerts table (for tracking sent alerts)
CREATE TABLE IF NOT EXISTS alerts (
    id TEXT PRIMARY KEY,
    target_id TEXT NOT NULL,
    target_name TEXT NOT NULL,
    type TEXT NOT NULL,
    severity TEXT NOT NULL,
    message TEXT NOT NULL,
    description TEXT,
    metadata TEXT, -- JSON
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    resolved_at TIMESTAMP,
    resolved BOOLEAN NOT NULL DEFAULT 0,
    FOREIGN KEY (target_id) REFERENCES targets(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_alerts_target_id ON alerts(target_id);
CREATE INDEX IF NOT EXISTS idx_alerts_created_at ON alerts(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_alerts_resolved ON alerts(resolved);
