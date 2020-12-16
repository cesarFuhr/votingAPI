CREATE TABLE IF NOT EXISTS sessions(
  id uuid PRIMARY KEY,
  originalAgenda uuid,
  duration BIGINT,
  creation TIMESTAMP
)