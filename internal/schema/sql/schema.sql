-- Version: 1.1
-- Description: Create Sessions table with index for expiry.
-- sessions expiry is managed by alexedwards/scs automatically.
CREATE TABLE sessions (
	token TEXT PRIMARY KEY,
	data BYTEA NOT NULL,
	expiry TIMESTAMPTZ NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);

-- Version: 1.2
-- Description: Coffee products
CREATE TABLE coffee (
	coffee_id INT GENERATED ALWAYS AS IDENTITY,
	coffee_name TEXT,
	coffee_description TEXT,
	coffee_price NUMERIC(5,2),
	coffee_caffeine TEXT,
	coffee_calories INT,
	PRIMARY KEY (coffee_id)
);

-- Version: 1.3
-- Description: Create User table, used for only admin currently
CREATE EXTENSION citext;

CREATE TABLE users (
id bigserial PRIMARY KEY,
created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(), 
name text NOT NULL,
email citext UNIQUE NOT NULL,
password_hash bytea NOT NULL
);

-- Version: 1.4
-- Description: Create Tokens for the privileged actions
CREATE TABLE tokens (
hash bytea PRIMARY KEY,
user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE, 
expiry timestamp(0) with time zone NOT NULL,
scope text NOT NULL
);

