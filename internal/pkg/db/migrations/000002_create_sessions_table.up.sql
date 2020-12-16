CREATE TABLE IF NOT EXISTS sessions(
  id uuid PRIMARY KEY,
  originalAgenda uuid,
  duration INTERVAL MINUTE,
  creation TIMESTAMP
)