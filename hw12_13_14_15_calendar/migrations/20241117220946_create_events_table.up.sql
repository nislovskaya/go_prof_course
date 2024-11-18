CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    date_time TIMESTAMP NOT NULL,
    duration INTERVAL NOT NULL,
    description TEXT,
    user_id VARCHAR(50) NOT NULL,
    notify_in TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);