-- +migrate Up

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
    embedding F32_BLOB(1024) ,
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



-- +migrate Down
drop table if exists characters;
drop table if exists character_facts;
drop table if exists users;
drop table if exists user_facts;
drop table if exists memories;