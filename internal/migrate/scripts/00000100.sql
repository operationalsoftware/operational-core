CREATE TABLE app_user (
	user_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	is_api_user BOOLEAN DEFAULT FALSE NOT NULL,
	username TEXT NOT NULL UNIQUE,
	email TEXT UNIQUE,
	first_name TEXT,
	last_name TEXT,
	created TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
	last_login TIMESTAMPTZ DEFAULT NULL,
	hashed_password TEXT NOT NULL,
	failed_login_attempts INTEGER DEFAULT 0 NOT NULL,
	login_blocked_until TIMESTAMPTZ DEFAULT NULL,
	permissions JSONB DEFAULT '{}'::JSONB NOT NULL,
	user_data JSONB DEFAULT '{}'::JSONB NOT NULL,
	session_duration_minutes INT
);


CREATE TABLE team (
	team_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	team_name TEXT NOT NULL,
	is_archived BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE user_team (
	user_team_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	user_id INT NOT NULL REFERENCES app_user(user_id) ON DELETE CASCADE,
	team_id INT NOT NULL REFERENCES team(team_id) ON DELETE CASCADE,
	role TEXT NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	UNIQUE (user_id, team_id)
);

CREATE TABLE recent_search (
	recent_search_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	search_term TEXT NOT NULL,
	search_entities TEXT[] NOT NULL,
	user_id INT REFERENCES app_user(user_id),
	last_searched_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

	CONSTRAINT unique_search_per_user UNIQUE (search_term, search_entities, user_id)
);


CREATE TABLE stock_item (
	stock_code TEXT PRIMARY KEY,
	description TEXT NOT NULL,
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE stock_item_change (
	stock_item_change_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	stock_code TEXT NOT NULL REFERENCES stock_item(stock_code) ON UPDATE CASCADE,
	stock_code_history TEXT,
	description TEXT,
	change_by INT REFERENCES app_user(user_id),
	changed_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE stock_transaction (
	stock_transaction_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	transaction_type TEXT NOT NULL,
	stock_code TEXT NOT NULL REFERENCES stock_item(stock_code),
	transaction_by INT NOT NULL REFERENCES app_user(user_id),
	transaction_note TEXT NOT NULL,
	timestamp TIMESTAMPTZ NOT NULL
);


CREATE TABLE stock_transaction_entry (
	stock_transaction_entry_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	account TEXT NOT NULL,
	location TEXT NOT NULL,
	bin TEXT NOT NULL,
	lot_number TEXT NOT NULL,
	quantity NUMERIC NOT NULL,
	running_total NUMERIC NOT NULL,
	stock_transaction_id INT NOT NULL REFERENCES stock_transaction(stock_transaction_id)
);


CREATE TABLE andon_issue (
	andon_issue_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	issue_name TEXT NOT NULL,
	is_archived BOOLEAN NOT NULL DEFAULT FALSE,
	parent_id INTEGER REFERENCES andon_issue(andon_issue_id),
	is_group BOOLEAN NOT NULL DEFAULT FALSE,
	assigned_to_team INTEGER REFERENCES team(team_id),
	severity TEXT,

	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	created_by INTEGER NOT NULL REFERENCES app_user(user_id),
	updated_at TIMESTAMPTZ,
	updated_by INTEGER REFERENCES app_user(user_id),

	UNIQUE (parent_id, issue_name)
);

CREATE TABLE andon_event (
	andon_event_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	issue_id INTEGER REFERENCES andon_issue(andon_issue_id),
	issue_description TEXT NOT NULL,
	source TEXT NOT NULL,
	location TEXT NOT NULL,
	raised_by INT NOT NULL REFERENCES app_user(user_id),
	raised_at TIMESTAMPTZ DEFAULT NOW(),
	acknowledged_by INT REFERENCES app_user(user_id),
	acknowledged_at TIMESTAMPTZ,
	resolved_by INT REFERENCES app_user(user_id),
	resolved_at TIMESTAMPTZ,
	cancelled_by INT REFERENCES app_user(user_id),
	cancelled_at TIMESTAMPTZ,
	status TEXT NOT NULL,
	last_updated TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE andon_change (
  andon_change_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  andon_event_id INT NOT NULL REFERENCES andon_event(andon_event_id) ON DELETE CASCADE,
  change_by INT NOT NULL REFERENCES app_user(user_id),
  change_at TIMESTAMPTZ DEFAULT NOW(),
  issue_id INTEGER REFERENCES andon_issue(andon_issue_id),
  issue_description TEXT,
  location TEXT,
  raised_by INT REFERENCES app_user(user_id),
  raised_at TIMESTAMPTZ,
  acknowledged_by INT REFERENCES app_user(user_id),
  acknowledged_at TIMESTAMPTZ,
  resolved_by INT REFERENCES app_user(user_id),
  resolved_at TIMESTAMPTZ,
  cancelled_by INT REFERENCES app_user(user_id),
  cancelled_at TIMESTAMPTZ,
  status TEXT
);

CREATE TABLE comment (
  comment_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  entity TEXT NOT NULL,
  entity_id INT NOT NULL,
  comment TEXT NOT NULL,
  commented_by INT NOT NULL REFERENCES app_user(user_id),
  commented_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_comment_entity ON comment(entity, entity_id);

