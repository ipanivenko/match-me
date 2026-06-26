-- ============================================
-- DATABASE SETUP FOR MATCHME APPLICATION
-- All tables, functions, triggers, and indexes
-- ============================================

-- Extensions (must be created first!)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS postgis;

-- ============================================
-- TABLES
-- ============================================

-- Users table
CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  email CITEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Parent profiles table (includes geography column)
CREATE TABLE IF NOT EXISTS parent_profiles (
  user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  name TEXT,
  gender TEXT,
  about TEXT,
  languages TEXT[] DEFAULT '{}',
  address_city TEXT,
  lat DOUBLE PRECISION,
  lon DOUBLE PRECISION,
  preferred_distance_km INTEGER CHECK (preferred_distance_km BETWEEN 0 AND 200),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  geog geography(Point, 4326)
);

-- Children table
CREATE TABLE IF NOT EXISTS children (
  user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  name TEXT,
  birthday DATE,
  gender TEXT,
  about_short TEXT,
  interests TEXT[] DEFAULT '{}',
  activity_level TEXT,
  limitations TEXT[] DEFAULT '{}',
  allergies TEXT[] DEFAULT '{}',
  play_styles TEXT[] DEFAULT '{}'
);

-- User photos table
CREATE TABLE IF NOT EXISTS user_photos (
  user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  photo_public_id TEXT NOT NULL,
  photo_version INTEGER NOT NULL
);

-- Matching preferences table
CREATE TABLE IF NOT EXISTS matching_preferences (
  user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  interests_weight INTEGER DEFAULT 1 CHECK (interests_weight BETWEEN 0 AND 5),
  activity_level_weight INTEGER DEFAULT 2 CHECK (activity_level_weight BETWEEN 0 AND 5),
  limitations_weight INTEGER DEFAULT 3 CHECK (limitations_weight BETWEEN 0 AND 5),
  allergies_weight INTEGER DEFAULT 3 CHECK (allergies_weight BETWEEN 0 AND 5),
  play_styles_weight INTEGER DEFAULT 1 CHECK (play_styles_weight BETWEEN 0 AND 5),
  max_age_difference INTEGER DEFAULT 2 CHECK (max_age_difference >= 0)
);

-- User reactions table (likes/dislikes and matches)
CREATE TABLE IF NOT EXISTS user_reactions (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  target_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  reaction TEXT NOT NULL CHECK (reaction IN ('like', 'dislike')),
  is_match BOOLEAN NOT NULL DEFAULT false,
  CONSTRAINT user_reactions_no_self CHECK (user_id <> target_user_id),
  CONSTRAINT user_reactions_user_target_unique UNIQUE (user_id, target_user_id)
);

-- Chat rooms table
CREATE TABLE IF NOT EXISTS chats (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user1_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  user2_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT chats_different_users CHECK (user1_id <> user2_id)
);

-- Messages table
CREATE TABLE IF NOT EXISTS messages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  chat_id UUID NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
  sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  content TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Connections table
CREATE TABLE IF NOT EXISTS connections (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  requester_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  target_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  status TEXT NOT NULL CHECK (status IN ('pending', 'accepted', 'rejected')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT connections_different_users CHECK (requester_user_id <> target_user_id),
  CONSTRAINT connections_unique_pair UNIQUE (requester_user_id, target_user_id)
);


-- ============================================
-- INDEXES
-- ============================================

CREATE INDEX IF NOT EXISTS parent_profiles_geog_idx ON parent_profiles USING GIST (geog);
CREATE UNIQUE INDEX IF NOT EXISTS idx_chats_unique_pair ON chats ((LEAST(user1_id, user2_id)), (GREATEST(user1_id, user2_id)));
CREATE INDEX IF NOT EXISTS idx_chats_user1 ON chats(user1_id);
CREATE INDEX IF NOT EXISTS idx_chats_user2 ON chats(user2_id);
CREATE INDEX IF NOT EXISTS idx_messages_chat_id ON messages(chat_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_messages_sender ON messages(sender_id);
CREATE INDEX IF NOT EXISTS idx_connections_requester ON connections(requester_user_id);
CREATE INDEX IF NOT EXISTS idx_connections_target ON connections(target_user_id);
CREATE INDEX IF NOT EXISTS idx_connections_status ON connections(status);


-- ============================================
-- FUNCTIONS AND TRIGGERS
-- ============================================

-- Function to update timestamps automatically
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS trigger AS $$
BEGIN
  NEW.updated_at := now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for parent_profiles timestamps
DROP TRIGGER IF EXISTS trg_parent_profiles_timestamps ON parent_profiles;
CREATE TRIGGER trg_parent_profiles_timestamps
BEFORE UPDATE ON parent_profiles
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- Function to keep geography column in sync with lat/lon
CREATE OR REPLACE FUNCTION set_parent_geog() RETURNS trigger AS $$
BEGIN
  IF NEW.lat IS NOT NULL AND NEW.lon IS NOT NULL THEN
    NEW.geog := ST_SetSRID(ST_MakePoint(NEW.lon, NEW.lat), 4326)::geography;
  ELSE
    NEW.geog := NULL;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for geography column
DROP TRIGGER IF EXISTS trg_parent_geog ON parent_profiles;
CREATE TRIGGER trg_parent_geog
BEFORE INSERT OR UPDATE ON parent_profiles
FOR EACH ROW EXECUTE FUNCTION set_parent_geog();

-- Backfill geography for existing rows
UPDATE parent_profiles
SET geog = CASE
             WHEN lat IS NOT NULL AND lon IS NOT NULL
             THEN ST_SetSRID(ST_MakePoint(lon, lat), 4326)::geography
             ELSE NULL
           END
WHERE geog IS NULL;

-- Function for updating is_match status when users like each other
CREATE OR REPLACE FUNCTION update_is_match()
RETURNS TRIGGER AS $$
BEGIN
  -- When the current user likes someone
  IF NEW.reaction = 'like' THEN
    -- Check if the target has already liked this user
    UPDATE user_reactions
    SET is_match = true
    WHERE user_id = NEW.target_user_id
      AND target_user_id = NEW.user_id
      AND reaction = 'like';

    -- If the reverse "like" exists, mark this one as a match too
    IF EXISTS (
      SELECT 1 FROM user_reactions
      WHERE user_id = NEW.target_user_id
        AND target_user_id = NEW.user_id
        AND reaction = 'like'
    ) THEN
      NEW.is_match := true;
    END IF;
  ELSE
    -- If it's a dislike, make sure match is false
    NEW.is_match := false;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for updating is_match automatically
DROP TRIGGER IF EXISTS trg_update_is_match ON user_reactions;
CREATE TRIGGER trg_update_is_match
BEFORE INSERT OR UPDATE ON user_reactions
FOR EACH ROW
EXECUTE FUNCTION update_is_match();


-- ============================================
-- MATCHING ALGORITHM FUNCTIONS
-- ============================================

-- Prefilter function: returns candidates within user's preferred radius
CREATE OR REPLACE FUNCTION prefilter_candidates_postgis(
  viewer uuid,
  p_limit int DEFAULT 500,
  p_offset int DEFAULT 0
)
RETURNS TABLE (
  candidate_user_id uuid,
  distance_km double precision,
  address_city text
)
LANGUAGE sql
AS $$
WITH me_parent AS (
  SELECT p.user_id, p.geog, p.preferred_distance_km
  FROM parent_profiles p
  WHERE p.user_id = viewer
    AND p.name IS NOT NULL AND btrim(p.name) <> ''
    AND p.gender IS NOT NULL AND btrim(p.gender) <> ''
    AND p.about IS NOT NULL AND btrim(p.about) <> ''
    AND array_length(p.languages,1) IS NOT NULL AND array_length(p.languages,1) > 0
    AND p.address_city IS NOT NULL AND btrim(p.address_city) <> ''
    AND p.geog IS NOT NULL
    AND p.preferred_distance_km IS NOT NULL
),
me_child AS (
  SELECT c.user_id
  FROM children c
  WHERE c.user_id = viewer
    AND c.name IS NOT NULL AND btrim(c.name) <> ''
    AND c.birthday IS NOT NULL
    AND c.gender IS NOT NULL AND btrim(c.gender) <> ''
    AND c.about_short IS NOT NULL AND btrim(c.about_short) <> ''
    AND array_length(c.interests,1) IS NOT NULL AND array_length(c.interests,1) > 0
    AND c.activity_level IS NOT NULL AND btrim(c.activity_level) <> ''
    AND array_length(c.allergies,1) IS NOT NULL AND array_length(c.allergies,1) > 0
    AND array_length(c.play_styles,1) IS NOT NULL AND array_length(c.play_styles,1) > 0
),
me AS (
  SELECT mp.user_id, mp.geog, mp.preferred_distance_km
  FROM me_parent mp JOIN me_child mc ON mc.user_id = mp.user_id
),
candidates_full AS (
  SELECT p.user_id, p.address_city, p.geog
  FROM parent_profiles p
  JOIN children c ON c.user_id = p.user_id
  JOIN me ON TRUE
  WHERE p.user_id <> me.user_id
    AND p.name IS NOT NULL AND btrim(p.name) <> ''
    AND p.gender IS NOT NULL AND btrim(p.gender) <> ''
    AND p.about IS NOT NULL AND btrim(p.about) <> ''
    AND array_length(p.languages,1) IS NOT NULL AND array_length(p.languages,1) > 0
    AND p.address_city IS NOT NULL AND btrim(p.address_city) <> ''
    AND p.geog IS NOT NULL
    AND c.name IS NOT NULL AND btrim(c.name) <> ''
    AND c.birthday IS NOT NULL
    AND c.gender IS NOT NULL AND btrim(c.gender) <> ''
    AND c.about_short IS NOT NULL AND btrim(c.about_short) <> ''
    AND array_length(c.interests,1) IS NOT NULL AND array_length(c.interests,1) > 0
    AND c.activity_level IS NOT NULL AND btrim(c.activity_level) <> ''
    AND array_length(c.allergies,1) IS NOT NULL AND array_length(c.allergies,1) > 0
    AND array_length(c.play_styles,1) IS NOT NULL AND array_length(c.play_styles,1) > 0
)
SELECT
  cf.user_id AS candidate_user_id,
  ST_Distance(me.geog, cf.geog) / 1000.0 AS distance_km,
  cf.address_city
FROM candidates_full cf
JOIN me ON TRUE
WHERE ST_DWithin(me.geog, cf.geog, me.preferred_distance_km * 1000.0)
ORDER BY distance_km ASC, cf.user_id
LIMIT p_limit OFFSET p_offset;
$$;

-- Profile completion percentage function
CREATE OR REPLACE FUNCTION profile_completion_percent(p_user_id uuid)
RETURNS numeric AS $$
DECLARE
  total_fields int := 14;
  parent_filled int := 0;
  child_filled int := 0;
BEGIN
  SELECT
    ((name IS NOT NULL AND name <> '')::int) +
    ((gender IS NOT NULL AND gender <> '')::int) +
    ((about IS NOT NULL AND about <> '')::int) +
    ((languages IS NOT NULL AND cardinality(languages) > 0)::int) +
    ((address_city IS NOT NULL AND address_city <> '')::int)
  INTO parent_filled
  FROM parent_profiles
  WHERE user_id = p_user_id;

  parent_filled := COALESCE(parent_filled, 0);

  SELECT
    ((name IS NOT NULL AND name <> '')::int) +
    ((birthday IS NOT NULL)::int) +
    ((gender IS NOT NULL AND gender <> '')::int) +
    ((about_short IS NOT NULL AND about_short <> '')::int) +
    ((interests IS NOT NULL AND cardinality(interests) > 0)::int) +
    ((activity_level IS NOT NULL AND activity_level <> '')::int) +
    ((limitations IS NOT NULL AND cardinality(limitations) > 0)::int) +
    ((allergies IS NOT NULL AND cardinality(allergies) > 0)::int) +
    ((play_styles IS NOT NULL AND cardinality(play_styles) > 0)::int)
  INTO child_filled
  FROM children
  WHERE user_id = p_user_id;

  child_filled := COALESCE(child_filled, 0);

  RETURN ROUND(((parent_filled + child_filled)::numeric / total_fields) * 100.0, 1);
END;

$$ LANGUAGE plpgsql;

-- Table to track unread messages
CREATE TABLE IF NOT EXISTS unread_messages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  message_id UUID NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT unread_messages_unique UNIQUE (user_id, message_id)
);

CREATE INDEX IF NOT EXISTS idx_unread_messages_user ON unread_messages(user_id);
CREATE INDEX IF NOT EXISTS idx_unread_messages_message ON unread_messages(message_id);

-- Function to mark message as unread for recipient
CREATE OR REPLACE FUNCTION mark_message_unread() RETURNS TRIGGER AS $$
BEGIN
  -- Add unread entry for the recipient (not the sender)
  INSERT INTO unread_messages (user_id, message_id)
  SELECT 
    CASE 
      WHEN c.user1_id = NEW.sender_id THEN c.user2_id
      ELSE c.user1_id
    END,
    NEW.id
  FROM chats c
  WHERE c.id = NEW.chat_id
  ON CONFLICT (user_id, message_id) DO NOTHING;
  
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to auto-create unread entries
DROP TRIGGER IF EXISTS trg_mark_message_unread ON messages;
CREATE TRIGGER trg_mark_message_unread
AFTER INSERT ON messages
FOR EACH ROW
EXECUTE FUNCTION mark_message_unread();

-- Add last_message_at column to chats for sorting
ALTER TABLE chats ADD COLUMN IF NOT EXISTS last_message_at TIMESTAMPTZ DEFAULT NOW();

-- Update last_message_at when new message is sent
CREATE OR REPLACE FUNCTION update_chat_last_message() RETURNS TRIGGER AS $$
BEGIN
  UPDATE chats 
  SET last_message_at = NEW.created_at 
  WHERE id = NEW.chat_id;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_update_chat_last_message ON messages;
CREATE TRIGGER trg_update_chat_last_message
AFTER INSERT ON messages
FOR EACH ROW
EXECUTE FUNCTION update_chat_last_message();

-- Backfill last_message_at for existing chats
UPDATE chats c
SET last_message_at = (
  SELECT MAX(m.created_at)
  FROM messages m
  WHERE m.chat_id = c.id
)
WHERE EXISTS (SELECT 1 FROM messages WHERE chat_id = c.id);


