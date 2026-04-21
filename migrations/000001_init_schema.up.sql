-- FitGenie Initial Schema
-- Created with golang-migrate

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Users
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Style Profiles
CREATE TABLE IF NOT EXISTS style_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    preferred_styles TEXT[] DEFAULT '{}',
    body_type VARCHAR(100) DEFAULT '',
    lifestyle VARCHAR(100) DEFAULT '',
    occasion TEXT[] DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Color Profiles
CREATE TABLE IF NOT EXISTS color_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    color_season VARCHAR(100) DEFAULT '',
    skin_tone VARCHAR(100) DEFAULT '',
    undertone VARCHAR(100) DEFAULT '',
    favorite_colors TEXT[] DEFAULT '{}',
    avoid_colors TEXT[] DEFAULT '{}',
    recommended_colors TEXT[] DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Clothing Items
CREATE TABLE IF NOT EXISTS clothing_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL,
    brand VARCHAR(255) DEFAULT '',
    size VARCHAR(50) DEFAULT '',
    primary_color VARCHAR(100) DEFAULT '',
    secondary_color VARCHAR(100) DEFAULT '',
    material VARCHAR(100) DEFAULT '',
    style VARCHAR(100) DEFAULT '',
    season TEXT[] DEFAULT '{}',
    occasion TEXT[] DEFAULT '{}',
    image_url TEXT DEFAULT '',
    notes TEXT DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Outfits
CREATE TABLE IF NOT EXISTS outfits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT DEFAULT '',
    style VARCHAR(100) DEFAULT '',
    occasion TEXT[] DEFAULT '{}',
    season TEXT[] DEFAULT '{}',
    weather VARCHAR(100) DEFAULT '',
    color_harmony_score DECIMAL DEFAULT 0,
    style_coherence_score DECIMAL DEFAULT 0,
    overall_score DECIMAL DEFAULT 0,
    rating INTEGER DEFAULT 0,
    worn BOOLEAN DEFAULT FALSE,
    favorite BOOLEAN DEFAULT FALSE,
    notes TEXT DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Outfit-Clothing Items (many-to-many)
CREATE TABLE IF NOT EXISTS outfit_clothing_items (
    outfit_id UUID NOT NULL REFERENCES outfits(id) ON DELETE CASCADE,
    clothing_item_id UUID NOT NULL REFERENCES clothing_items(id) ON DELETE CASCADE,
    PRIMARY KEY (outfit_id, clothing_item_id)
);

-- Outfit Recommendations
CREATE TABLE IF NOT EXISTS outfit_recommendations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    outfit_id UUID NOT NULL REFERENCES outfits(id) ON DELETE CASCADE,
    confidence DECIMAL DEFAULT 0,
    reasoning TEXT DEFAULT '',
    requested_occasion VARCHAR(100) DEFAULT '',
    requested_season VARCHAR(100) DEFAULT '',
    requested_weather VARCHAR(100) DEFAULT '',
    requested_style VARCHAR(100) DEFAULT '',
    viewed BOOLEAN DEFAULT FALSE,
    accepted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Favorite Outfits
CREATE TABLE IF NOT EXISTS favorite_outfits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    outfit_id UUID NOT NULL REFERENCES outfits(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, outfit_id)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_clothing_items_user_id ON clothing_items(user_id);
CREATE INDEX IF NOT EXISTS idx_clothing_items_category ON clothing_items(category);
CREATE INDEX IF NOT EXISTS idx_outfits_user_id ON outfits(user_id);
CREATE INDEX IF NOT EXISTS idx_outfit_recommendations_outfit_id ON outfit_recommendations(outfit_id);
