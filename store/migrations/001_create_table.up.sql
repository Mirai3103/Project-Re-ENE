
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
    content blob ,
    created_at TIMESTAMP  DEFAULT CURRENT_TIMESTAMP
);


create table if not exists characters (
    id text primary key,
    name text ,
    base_prompt text ,
    description text ,
    updated_at timestamp  default current_timestamp,
    created_at timestamp  default current_timestamp
);

create table if not exists character_facts (
    id text primary key,
    character_id text ,
    name text ,
    value text ,
    type text ,
    created_at timestamp  default current_timestamp,
    updated_at timestamp  default current_timestamp
);


create table if not exists users (
    id text primary key,
    name text ,
    bio text ,
    created_at timestamp  default current_timestamp,
    updated_at timestamp  default current_timestamp
);

create table if not exists user_facts (
    id text primary key,
    name text ,
    value text ,
    type text ,
    user_id text ,
    created_at timestamp  default current_timestamp,
    updated_at timestamp  default current_timestamp
);

create table if not exists memories (
    id text primary key,
    user_id text ,
    character_id text ,
    content text ,
    embedding blob ,
    importance REAL  default 0.0,
    confidence REAL  default 0.0,
    source text ,
    tags text ,
    access_count INTEGER  default 0,
    decay_score REAL  default 0.0,
    last_accessed_at timestamp  default current_timestamp,
    created_at timestamp  default current_timestamp,
    updated_at timestamp  default current_timestamp
);
