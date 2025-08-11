CREATE TABLE conversations (
  id           uuid PRIMARY KEY NOT NULL,
  type         text NOT NULL CHECK (type IN ('direct','group')),
  direct_key   text UNIQUE, -- only for direct, and the key should be canonical
  group_name   text, -- only for group, optional group name
  created_at   timestamp WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE conversation_participants (
  conversation_id uuid NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
  uid             uuid NOT NULL,
  role            text NOT NULL DEFAULT 'member',
  last_read_seq   bigint NOT NULL DEFAULT 0,
  joined_at       timestamp WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (conversation_id, uid)
);


CREATE TABLE conversation_counters (
  conversation_id uuid PRIMARY KEY REFERENCES conversations(id) ON DELETE CASCADE,
  last_seq        bigint NOT NULL DEFAULT 0, -- last message sequence number
);

-- create counter row with next_seq=1 whenever a conversation is created
-- (do it in app tx or via trigger)

CREATE TABLE messages (
  id              uuid PRIMARY KEY,
  conversation_id uuid NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
  seq             bigint NOT NULL,
  from_uid        uuid NOT NULL,
  body            jsonb NOT NULL,
  created_at      timestamp WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (conversation_id, seq)
);


-- Delivery status 
CREATE TABLE message_receipts (
  message_id      uuid REFERENCES messages(id),
  uid             uuid NOT NULL,
  delivered_at    timestamp WITH TIME ZONE,
  read_at         timestamp WITH TIME ZONE,
  PRIMARY KEY (message_id, uid)
);