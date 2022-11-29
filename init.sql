create table entries(
  id          SERIAL PRIMARY KEY,
  created_at  TIMESTAMP,
  tags        JSONB,
  views       INT
  title       VARCHAR(50),
  body        TEXT,
)
