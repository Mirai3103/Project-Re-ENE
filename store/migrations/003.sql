-- +migrate Up

create table if not exists characters (
    id text primary key,
    name text ,
    base_prompt text ,
    description text ,
    updated_at timestamp not null default current_timestamp,
    created_at timestamp not null default current_timestamp
);

create table if not exists character_facts (
    id text primary key,
    character_id text not null,
    name text not null,
    value text not null,
    type text not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);


create table if not exists users (
    id text primary key,
    name text not null,
    bio text not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

create table if not exists user_facts (
    id text primary key,
    name text not null,
    value text not null,
    type text not null,
    user_id text not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

create table if not exists memories (
    id text primary key,
    user_id text not null,
    character_id text not null,
    content text not null,
    embedding F32_BLOB(1024) not null,
    importance REAL not null default 0.0,
    last_accessed_at timestamp not null default current_timestamp,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);



-- +migrate Down
drop table if exists characters;
drop table if exists character_facts;
drop table if exists users;
drop table if exists user_facts;
drop table if exists memories;