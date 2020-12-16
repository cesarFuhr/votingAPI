CREATE TABLE IF NOT EXISTS votes(
  associateID Varchar(50),
  sessionID uuid,
  document Varchar(50),
  vote Varchar(1),
  creation TIMESTAMP,
  UNIQUE (associateID, sessionID)
)
