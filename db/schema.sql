-- Schema applied by the Upsun deploy hook. Every statement is idempotent so
-- that re-running it on each deploy never drops or duplicates data.
-- Mirrors internal/models/testdata/setup.sql (without its seed user).

CREATE TABLE IF NOT EXISTS snippets (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created DATETIME NOT NULL,
    expires DATETIME NOT NULL,
    INDEX idx_snippets_created (created)
);

CREATE TABLE IF NOT EXISTS users (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created DATETIME NOT NULL,
    CONSTRAINT users_uc_email UNIQUE (email)
);

-- Required by github.com/alexedwards/scs/mysqlstore.
CREATE TABLE IF NOT EXISTS sessions (
    token CHAR(43) PRIMARY KEY,
    data BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL,
    INDEX sessions_expiry_idx (expiry)
);
