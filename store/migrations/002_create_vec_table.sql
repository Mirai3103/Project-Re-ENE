-- +migrate Up

CREATE TABLE IF NOT EXISTS conversations (
    id TEXT PRIMARY KEY,
    title TEXT ,
    max_window_size INTEGER ,
    character_id TEXT ,
    user_id TEXT ,
    current_summary TEXT ,
    created_at TIMESTAMP  DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP  DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS conversation_messages (
    id TEXT PRIMARY KEY,
    conversation_id TEXT ,
    role TEXT ,
    content TEXT ,
    created_at TIMESTAMP  DEFAULT CURRENT_TIMESTAMP
);




-- +migrate Down

DROP TABLE IF EXISTS user_fact_vecs;
DROP TABLE IF EXISTS user_facts;
