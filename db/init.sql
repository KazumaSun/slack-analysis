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

-- -- 初期化
-- INSERT INTO users (slack_id, name, attributes)
-- VALUES
--   ('U12345678', 'Alice', '{"age": 30, "location": "Tokyo"}'),
--   ('U87654321', 'Bob', '{"age": 25, "location": "Osaka"}');
