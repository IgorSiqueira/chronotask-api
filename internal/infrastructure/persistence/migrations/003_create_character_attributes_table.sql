-- Create character_attributes table
CREATE TABLE IF NOT EXISTS character_attributes (
    id VARCHAR(255) PRIMARY KEY,
    attribute_name VARCHAR(50) NOT NULL,
    value INTEGER NOT NULL DEFAULT 0,
    character_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    -- Foreign key constraint
    CONSTRAINT fk_attribute_character
        FOREIGN KEY (character_id)
        REFERENCES characters(id)
        ON DELETE CASCADE,

    -- Unique constraint to prevent duplicate attribute names per character
    CONSTRAINT uq_character_attribute_name
        UNIQUE (character_id, attribute_name)
);

-- Create index on character_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_character_attributes_character_id ON character_attributes(character_id);

-- Create index on attribute_name for filtering by specific attributes
CREATE INDEX IF NOT EXISTS idx_character_attributes_name ON character_attributes(attribute_name);

-- Create index on value for queries that filter/sort by attribute values
CREATE INDEX IF NOT EXISTS idx_character_attributes_value ON character_attributes(value DESC);

-- Create index on created_at for sorting
CREATE INDEX IF NOT EXISTS idx_character_attributes_created_at ON character_attributes(created_at);
