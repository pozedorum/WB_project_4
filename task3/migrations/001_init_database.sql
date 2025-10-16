CREATE TABLE IF NOT EXISTS events (
		id SERIAL PRIMARY KEY,
		user_id VARCHAR(255) NOT NULL,
		title VARCHAR(500),
		text TEXT NOT NULL,
		datetime TIMESTAMP WITH TIME ZONE NOT NULL,
		remind_before INTEGER DEFAULT 0,
		is_archived BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

CREATE INDEX IF NOT EXISTS idx_events_user_id ON events(user_id);
CREATE INDEX IF NOT EXISTS idx_events_datetime ON events(datetime);
CREATE INDEX IF NOT EXISTS idx_events_user_datetime ON events(user_id, datetime);