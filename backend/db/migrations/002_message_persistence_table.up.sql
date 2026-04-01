CREATE TABLE IF NOT EXISTS conversations (
  id           uuid PRIMARY KEY NOT NULL,
  type         text NOT NULL CHECK (type IN ('direct','group')),
  group_name   text, -- only for group, optional group name
  group_avatar_url text, -- only for group, optional group avatar URL
  created_at   timestamp WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS conversation_participants (
  conversation_id uuid NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
  uid             uuid NOT NULL REFERENCES users(uid) ON DELETE CASCADE,
  role            text NOT NULL DEFAULT 'member', -- Role: 'creator', 'admin', 'member'
  last_read_seq   bigint NOT NULL DEFAULT 0,
  joined_at       timestamp WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (conversation_id, uid)
);

CREATE INDEX idx_cp_uid ON conversation_participants(uid);


CREATE TABLE IF NOT EXISTS conversation_counters (
  conversation_id uuid PRIMARY KEY REFERENCES conversations(id) ON DELETE CASCADE,
  last_seq        bigint NOT NULL DEFAULT 0 -- last message sequence number
);

CREATE TABLE IF NOT EXISTS messages (
  id              uuid PRIMARY KEY,
  conversation_id uuid NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
  seq             bigint NOT NULL,
  from_uid        uuid NOT NULL REFERENCES users(uid) ON DELETE CASCADE,
  body            jsonb NOT NULL,
  created_at      timestamp WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (conversation_id, seq)
);
