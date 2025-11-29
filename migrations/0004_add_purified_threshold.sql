-- Add purified_threshold column to grumbles table
-- This allows each grumble to have a custom threshold for purification

-- Add purified_threshold column with default value of 10
-- purified_thresholdは余裕を持たせて1-10000に
ALTER TABLE grumbles ADD COLUMN purified_threshold INTEGER NOT NULL DEFAULT 10 CHECK (purified_threshold BETWEEN 1 AND 10000);

-- Update existing purified grumbles to set their threshold to their current vibe_count
-- This preserves the historical accuracy of when they were purified
UPDATE grumbles SET purified_threshold = vibe_count WHERE is_purified = TRUE;

-- Add same column to grumbles_archive table for consistency
ALTER TABLE grumbles_archive ADD COLUMN purified_threshold INTEGER NOT NULL DEFAULT 10 CHECK (purified_threshold BETWEEN 1 AND 10000);

-- Update existing archived purified grumbles
UPDATE grumbles_archive SET purified_threshold = vibe_count WHERE is_purified = TRUE;