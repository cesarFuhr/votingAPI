CREATE TABLE IF NOT EXISTS agendas(
  id uuid PRIMARY KEY,
  description VARCHAR(500),
  creation TIMESTAMP
)