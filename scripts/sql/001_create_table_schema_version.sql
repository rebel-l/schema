CREATE TABLE IF NOT EXISTS schema_script (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, -- TODO: AUTOINCREMENT is not available in every database
  script_name TEXT NOT NULL,
  executed_at DATETIME NOT NULL,
  execution_status VARCHAR(100) NOT NULL,
  app_version CHAR(30) NULL,
  error_msg TEXT NULL
);
