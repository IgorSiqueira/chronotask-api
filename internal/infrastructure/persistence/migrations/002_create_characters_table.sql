-- Create characters table
CREATE TABLE IF NOT EXISTS characters (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    level INTEGER NOT NULL DEFAULT 1,
    current_xp INTEGER NOT NULL DEFAULT 0,
    total_xp INTEGER NOT NULL DEFAULT 0,
    user_id VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    -- Foreign key constraint
    CONSTRAINT fk_character_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- Create index on user_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_characters_user_id ON characters(user_id);

-- Create index on level for leaderboard queries
CREATE INDEX IF NOT EXISTS idx_characters_level ON characters(level DESC);

-- Create index on total_xp for leaderboard queries
CREATE INDEX IF NOT EXISTS idx_characters_total_xp ON characters(total_xp DESC);

-- Create index on created_at for sorting
CREATE INDEX IF NOT EXISTS idx_characters_created_at ON characters(created_at);
