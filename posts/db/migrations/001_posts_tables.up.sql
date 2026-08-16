CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS posts (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    seq bigserial NOT NULL,
    author_uid uuid NOT NULL,          -- no FK: separate database, intentional
    body text NOT NULL,
    likes_count integer NOT NULL DEFAULT 0,
    comments_count integer NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX idx_posts_seq ON posts(seq);
CREATE INDEX idx_posts_author ON posts(author_uid);

CREATE TABLE IF NOT EXISTS post_likes (
    post_id uuid NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id uuid NOT NULL,             -- no FK: separate database, intentional
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (post_id, user_id)     -- makes like/unlike idempotent
);

CREATE TABLE IF NOT EXISTS comments (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id uuid NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    author_uid uuid NOT NULL,          -- no FK: separate database, intentional
    body text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_comments_post ON comments(post_id, created_at);
