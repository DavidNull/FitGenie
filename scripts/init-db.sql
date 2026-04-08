-- Initialize pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- Add vector column to clothing_items for image embeddings (future use)
ALTER TABLE clothing_items ADD COLUMN IF NOT EXISTS image_embedding vector(512);

-- Create index for similarity search
CREATE INDEX IF NOT EXISTS idx_clothing_image_embedding 
ON clothing_items USING ivfflat (image_embedding vector_cosine_ops)
WITH (lists = 100);
