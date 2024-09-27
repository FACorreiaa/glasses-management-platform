CREATE TABLE user_sessions (
                             token TEXT PRIMARY KEY,
                             user_id UUID REFERENCES "user"(user_id),
                             created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
