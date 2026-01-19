-- Initial schema for Idea Hamster Phase 1
-- Run this in your Supabase SQL editor

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Ideas table
CREATE TABLE ideas (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(80) NOT NULL,
    description VARCHAR(500) NOT NULL,
    category VARCHAR(20) NOT NULL CHECK (category IN ('Frontend', 'Backend', 'Full Stack')),
    tags TEXT[] DEFAULT '{}',
    submitter VARCHAR(100),
    vote_count INTEGER DEFAULT 0,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'eligible', 'building', 'built', 'archived')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Votes table
CREATE TABLE votes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    idea_id UUID NOT NULL REFERENCES ideas(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(idea_id, email) -- One vote per email per idea
);

-- Indexes for performance
CREATE INDEX idx_ideas_vote_count ON ideas(vote_count DESC);
CREATE INDEX idx_ideas_created_at ON ideas(created_at DESC);
CREATE INDEX idx_ideas_status ON ideas(status);
CREATE INDEX idx_votes_idea_id ON votes(idea_id);
CREATE INDEX idx_votes_email ON votes(email);

-- Function to update vote count
CREATE OR REPLACE FUNCTION update_idea_vote_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE ideas
        SET vote_count = vote_count + 1,
            status = CASE
                WHEN vote_count + 1 >= 50 THEN 'eligible'
                ELSE status
            END,
            updated_at = NOW()
        WHERE id = NEW.idea_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE ideas
        SET vote_count = GREATEST(vote_count - 1, 0),
            updated_at = NOW()
        WHERE id = OLD.idea_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Trigger to automatically update vote counts
CREATE TRIGGER trigger_update_vote_count
AFTER INSERT OR DELETE ON votes
FOR EACH ROW
EXECUTE FUNCTION update_idea_vote_count();

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to automatically update updated_at
CREATE TRIGGER trigger_ideas_updated_at
BEFORE UPDATE ON ideas
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Row Level Security (RLS) Policies
-- Enable RLS
ALTER TABLE ideas ENABLE ROW LEVEL SECURITY;
ALTER TABLE votes ENABLE ROW LEVEL SECURITY;

-- Ideas: Everyone can read, authenticated users can insert
CREATE POLICY "Anyone can view ideas"
    ON ideas FOR SELECT
    TO anon, authenticated
    USING (true);

CREATE POLICY "Anyone can insert ideas"
    ON ideas FOR INSERT
    TO anon, authenticated
    WITH CHECK (true);

-- Votes: Everyone can read, authenticated users can insert their own votes
CREATE POLICY "Anyone can view votes"
    ON votes FOR SELECT
    TO anon, authenticated
    USING (true);

CREATE POLICY "Anyone can insert votes"
    ON votes FOR INSERT
    TO anon, authenticated
    WITH CHECK (true);

-- Prevent vote deletion by users (only admins can delete)
CREATE POLICY "Prevent vote deletion"
    ON votes FOR DELETE
    TO authenticated
    USING (false);

-- Insert some seed data for development
INSERT INTO ideas (title, description, category, tags, submitter) VALUES
(
    'AI-Powered Recipe Generator',
    'Generate custom recipes based on ingredients you have at home. Uses AI to suggest creative combinations and cooking instructions.',
    'Full Stack',
    ARRAY['AI', 'Food', 'Mobile'],
    'Chef Mike'
),
(
    'Retro Pomodoro Timer',
    'A productivity timer with 90s aesthetics. Track your work sessions with pixel art animations and chiptune sounds.',
    'Frontend',
    ARRAY['Productivity', 'Retro', 'Timer'],
    'ProductivityGuru'
),
(
    'Local Event Discovery API',
    'API that aggregates events from multiple sources and provides a unified interface for discovering local happenings.',
    'Backend',
    ARRAY['API', 'Events', 'Location'],
    'DevCommunity'
);

-- Add some initial votes for testing
INSERT INTO votes (idea_id, email, verified)
SELECT id, 'test' || generate_series(1, 15) || '@example.com', true
FROM ideas
LIMIT 1;
