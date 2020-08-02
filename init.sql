CREATE TABLE metrics (
  timestamp INT NOT NULL,
  cpuLoad FLOAT NULL DEFAULT NULL,
  concurrency INT NULL DEFAULT NULL,
  PRIMARY KEY (timestamp)
)