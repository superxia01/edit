-- Edit Business Database Schema
-- PostgreSQL 15+

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    auth_center_user_id UUID NOT NULL UNIQUE,
    role VARCHAR(50) NOT NULL DEFAULT 'USER',
    profile JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Bloggers table
CREATE TABLE IF NOT EXISTS bloggers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    xhs_id VARCHAR(50) UNIQUE,
    blogger_name VARCHAR(100),
    avatar_url VARCHAR(500),
    description TEXT,
    followers_count INTEGER NOT NULL DEFAULT 0,
    blogger_url VARCHAR(500),
    capture_timestamp BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Notes table
CREATE TABLE IF NOT EXISTS notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    url VARCHAR(500) NOT NULL UNIQUE,
    title VARCHAR(500),
    author VARCHAR(100),
    content TEXT,
    tags TEXT[],
    image_urls TEXT[],
    video_url VARCHAR(500),
    note_type VARCHAR(20),
    cover_image_url VARCHAR(500),
    likes INTEGER NOT NULL DEFAULT 0,
    collects INTEGER NOT NULL DEFAULT 0,
    comments INTEGER NOT NULL DEFAULT 0,
    publish_date BIGINT,
    capture_timestamp BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_notes_author ON notes(author);
CREATE INDEX IF NOT EXISTS idx_notes_publish_date ON notes(publish_date DESC);
CREATE INDEX IF NOT EXISTS idx_notes_tags ON notes USING GIN(tags);
CREATE INDEX IF NOT EXISTS idx_notes_likes ON notes(likes DESC);
CREATE INDEX IF NOT EXISTS idx_bloggers_xhs_id ON bloggers(xhs_id);
CREATE INDEX IF NOT EXISTS idx_bloggers_followers ON bloggers(followers_count DESC);

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Triggers for automatic updated_at
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_bloggers_updated_at
    BEFORE UPDATE ON bloggers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_notes_updated_at
    BEFORE UPDATE ON notes
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
