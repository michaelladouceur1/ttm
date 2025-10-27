CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY,
    title TEXT,
    description TEXT,
    category TEXT,
    priority TEXT,
    status TEXT,
    opened_at TIMESTAMP,
    closed_at TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sessions (
    id INTEGER PRIMARY KEY,
    task_id INTEGER,
    start_time TIMESTAMP,
    end_time TIMESTAMP
);