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
-- Description: Coffee products, represent price with [S,M,L,X]
CREATE TABLE coffee (
	coffee_id INT GENERATED ALWAYS AS IDENTITY,
	coffee_name TEXT,
	coffee_description TEXT,
	coffee_price NUMERIC(5,2),
	PRIMARY KEY (coffee_id)
);

-- Version: 1.3
-- Description: Cart will map to session token. 
-- Each cart will contain user session and coffees added.
-- If session is deleted, cart will also be deleted.
CREATE TABLE shoppingcart (
	cart_id INT GENERATED ALWAYS AS IDENTITY,
	sessions_token TEXT,
	customs_coffees_id INT,
	PRIMARY KEY (cart_id),
	CONSTRAINT fk_session
		FOREIGN KEY(sessions_token)
			REFERENCES sessions(token)
			ON DELETE CASCADE
);