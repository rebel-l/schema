CREATE TABLE IF NOT EXISTS schema_version (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  script_name TEXT NOT NULL,
  executed_at DATETIME NOT NULL,
  execution_status VARCHAR(100) NOT NULL,
  error_msg TEXT NULL
);
