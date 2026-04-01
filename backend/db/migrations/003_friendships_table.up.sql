CREATE TABLE IF NOT EXISTS friendships (
    requester_uid uuid NOT NULL REFERENCES users(uid) ON DELETE CASCADE,
    receiver_uid  uuid NOT NULL REFERENCES users(uid) ON DELETE CASCADE,
    status        text NOT NULL CHECK (status IN ('pending', 'accepted', 'blocked')),
    created_at    timestamp WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    timestamp WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_friendships_requester_receiver ON friendships(requester_uid, receiver_uid);
CREATE INDEX idx_friendships_receiver ON friendships(receiver_uid);
