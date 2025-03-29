CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  slack_id TEXT UNIQUE,
  name TEXT,
  attributes JSONB,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS activity_logs (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users(id),
  timestamp TIMESTAMP,
  status TEXT -- online / away / etc.
);
