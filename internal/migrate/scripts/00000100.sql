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



CREATE TABLE comment (
  comment_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  entity TEXT NOT NULL,
  entity_id INT NOT NULL,
  comment TEXT NOT NULL,
  commented_by INT NOT NULL REFERENCES app_user(user_id),
  commented_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_comment_entity ON comment(entity, entity_id);


CREATE TABLE stock_item (
	stock_item_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	stock_code TEXT NOT NULL,
	description TEXT NOT NULL,
	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
-- CREATE TABLE stock_item (
-- 	stock_code TEXT PRIMARY KEY,
-- 	description TEXT NOT NULL,
-- 	created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
-- );


CREATE TABLE stock_item_change (
	stock_item_change_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	stock_item_id INT NOT NULL REFERENCES stock_item(stock_item_id),
	stock_code TEXT,
	description TEXT,
	change_by INT REFERENCES app_user(user_id),
	changed_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
-- CREATE TABLE stock_item_change (
-- 	stock_item_change_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
-- 	stock_code TEXT NOT NULL REFERENCES stock_item(stock_code) ON UPDATE CASCADE,
-- 	stock_code_history TEXT,
-- 	description TEXT,
-- 	change_by INT REFERENCES app_user(user_id),
-- 	changed_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
-- );


CREATE TABLE stock_transaction (
	stock_transaction_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	transaction_type TEXT NOT NULL,
	-- stock_code TEXT NOT NULL REFERENCES stock_item(stock_code),
	stock_item_id INT NOT NULL REFERENCES stock_item(stock_item_id),
	transaction_by INT NOT NULL REFERENCES app_user(user_id),
	transaction_note TEXT NOT NULL,
	timestamp TIMESTAMPTZ NOT NULL
);
-- CREATE TABLE stock_transaction (
-- 	stock_transaction_id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
-- 	transaction_type TEXT NOT NULL,
-- 	stock_code TEXT NOT NULL REFERENCES stock_item(stock_code),
-- 	transaction_by INT NOT NULL REFERENCES app_user(user_id),
-- 	transaction_note TEXT NOT NULL,
-- 	timestamp TIMESTAMPTZ NOT NULL
-- );


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
	assigned_team INTEGER REFERENCES team(team_id),
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


CREATE VIEW andon_issue_tree_view AS
WITH RECURSIVE andon_issue_tree AS (
    SELECT
        ai.andon_issue_id,
        ai.issue_name,
        ai.parent_id,
        ARRAY[ai.issue_name] AS name_path,
        1 AS depth,
        ai.is_group,
        ai.is_archived,
        (
            SELECT COUNT(*)
            FROM andon_issue c
            WHERE c.parent_id = ai.andon_issue_id
        ) AS children_count
    FROM andon_issue ai
    WHERE ai.parent_id IS NULL

    UNION ALL

    SELECT
        child.andon_issue_id,
        child.issue_name,
        child.parent_id,
        parent.name_path || child.issue_name,
        parent.depth + 1,
        child.is_group,
        child.is_archived,
        (
            SELECT COUNT(*)
            FROM andon_issue c
            WHERE c.parent_id = child.andon_issue_id
        ) AS children_count
    FROM andon_issue child
    JOIN andon_issue_tree parent ON child.parent_id = parent.andon_issue_id
    WHERE parent.is_group = TRUE
),
down_depths AS (
    SELECT
        g.andon_issue_id,

        -- Downward depth
        (
            SELECT COALESCE(MAX(depth) - 1, 0)
            FROM (
                WITH RECURSIVE downward AS (
                    SELECT andon_issue_id, parent_id, 1 AS depth
                    FROM andon_issue
                    WHERE andon_issue_id = g.andon_issue_id

                    UNION ALL

                    SELECT ai.andon_issue_id, ai.parent_id, d.depth + 1
                    FROM andon_issue ai
                    JOIN downward d ON ai.parent_id = d.andon_issue_id
                )
                SELECT * FROM downward
            ) AS down_sub
        ) AS down_depth
    FROM andon_issue g
)
SELECT
    ait.andon_issue_id,
    ait.issue_name,
    ait.parent_id,
    ait.name_path,
    ait.depth,
    ait.is_group,
    ait.is_archived,
    ait.children_count,
    ai.severity,
    ai.assigned_team,
    t.team_name AS assigned_team_name,
    ai.created_at,
    ai.created_by,
    cu.username AS created_by_username,
    ai.updated_at,
    ai.updated_by,
    uu.username AS updated_by_username,
    COALESCE(dd.down_depth, 0) + 1 AS down_depth
FROM
    andon_issue_tree ait
    INNER JOIN andon_issue ai USING(andon_issue_id)
    LEFT JOIN team t ON t.team_id = ai.assigned_team
    INNER JOIN app_user cu ON cu.user_id = ai.created_by
    LEFT JOIN app_user uu ON uu.user_id = ai.updated_by
    LEFT JOIN down_depths dd ON dd.andon_issue_id = ait.andon_issue_id;


CREATE VIEW andon_issue_view AS
SELECT
    andon_issue_id,
    issue_name,
    parent_id,
    name_path,
    depth,
    is_archived,
    severity,
    assigned_team,
    assigned_team_name,
    -- COALESCE(assigned_team, 0) AS assigned_team,
    -- COALESCE(assigned_team_name, '') AS assigned_team_name,
    created_at,
    created_by,
    created_by_username,
    updated_at,
    updated_by,
    updated_by_username
FROM andon_issue_tree_view
WHERE is_group = false;


CREATE VIEW andon_issue_group_view AS
SELECT
    andon_issue_id,
    issue_name,
    parent_id,
    name_path,
    depth,
    children_count,
    is_archived,
    is_group,
    down_depth
FROM andon_issue_tree_view
WHERE is_group = true;


CREATE VIEW andon_view AS
SELECT
    ae.andon_event_id,
    ae.issue_description,
    ae.issue_id,
    ae.source,
    ae.location,
    ae.status,
    ae.raised_at,
    ae.raised_by,
    ae.acknowledged_at,
    ae.resolved_at,
    ae.cancelled_at,
    ae.last_updated,
    aiv.issue_name,
    aiv.assigned_team,
    aiv.assigned_team_name,
    aiv.name_path,
    u.username AS raised_by_username,
    acku.username AS acknowledged_by_username,
    ru.username AS resolved_by_username,
    cu.username AS cancelled_by_username,
    aiv.severity
FROM
    andon_event ae
    INNER JOIN app_user u ON ae.raised_by = u.user_id
    LEFT JOIN app_user acku ON ae.acknowledged_by = acku.user_id
    LEFT JOIN app_user ru ON ae.resolved_by = ru.user_id
    LEFT JOIN app_user cu ON ae.cancelled_by = cu.user_id
    INNER JOIN andon_issue_view aiv ON ae.issue_id = aiv.andon_issue_id;


CREATE VIEW user_view AS
SELECT
    u.user_id,
    u.is_api_user,
    u.username,
    u.email,
    u.first_name,
    u.last_name,
    u.created,
    u.last_login,
    u.session_duration_minutes,
    u.permissions,
    COALESCE(json_agg(
        json_build_object(
            'team_id', t.team_id,
            'team_name', t.team_name,
            'role', ut.role
        )
    ) FILTER (WHERE t.team_id IS NOT NULL), '[]') AS teams
FROM
    app_user u
LEFT JOIN user_team ut ON u.user_id = ut.user_id
LEFT JOIN team t ON ut.team_id = t.team_id
GROUP BY u.user_id;