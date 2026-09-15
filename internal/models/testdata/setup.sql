CREATE TABLE snippets (
	id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	title VARCHAR(100) NOT NULL,
	content TEXT NOT NULL,
	created DATETIME NOT NULL,
	expires DATETIME NOT NULL
);

CREATE INDEX idx_snippets_created ON snippets(created);

INSERT INTO snippets (id, title, content, created, expires) VALUES (
	1,
	'Expiring soon',
	'An old silent pond...',
	UTC_TIMESTAMP(),
	DATE_ADD(UTC_TIMESTAMP(), INTERVAL 1 DAY)
);

INSERT INTO snippets (id, title, content, created, expires) VALUES (
	2,
	'Long lived',
	'Over the wintry forest...',
	UTC_TIMESTAMP(),
	DATE_ADD(UTC_TIMESTAMP(), INTERVAL 10 DAY)
);

INSERT INTO snippets (id, title, content, created, expires) VALUES (
	3,
	'Already expired',
	'First autumn morning...',
	DATE_SUB(UTC_TIMESTAMP(), INTERVAL 10 DAY),
	DATE_SUB(UTC_TIMESTAMP(), INTERVAL 1 DAY)
);

CREATE TABLE users (
	id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
	name VARCHAR(255) NOT NULL,
	email VARCHAR(255) NOT NULL,
	hashed_password CHAR(60) NOT NULL,
	created DATETIME NOT NULL
);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);

INSERT INTO users (name, email, hashed_password, created) VALUES (
	'Alice Jones',
	'alice@example.com',
	'$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG',
	'2022-01-01 09:18:24'
);

