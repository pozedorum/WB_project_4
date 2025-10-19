CREATE TABLE IF NOT EXISTS events (
		id SERIAL PRIMARY KEY,
		usertoken VARCHAR(255) NOT NULL,
		title VARCHAR(500) NOT NULL,
		text TEXT,
		datetime TIMESTAMP WITH TIME ZONE NOT NULL,
		remind_before INTEGER DEFAULT 0,
		is_archived BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

CREATE INDEX IF NOT EXISTS idx_events_usertoken ON events(usertoken);
CREATE INDEX IF NOT EXISTS idx_events_datetime ON events(datetime);
CREATE INDEX IF NOT EXISTS idx_events_user_datetime ON events(usertoken, datetime);