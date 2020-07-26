create schema authorization flcrd;

create extension "uuid-ossp";

create table flcrd.deck (
    id            uuid      not null default uuid_generate_v4(),
    name          text      not null,
    description   text      not null default '',
    created       timestamp not null default now(),
    created_by    uuid      not null,
    public        boolean   not null default false,
    search_tokens tsvector,

    primary key (id)
);
create unique index deck_name_user_idx on flcrd.deck(name, created_by);

create table flcrd.flashcard (
    id      uuid      not null default uuid_generate_v4(),
    deck_id uuid      not null references flcrd.deck on delete cascade,
    front   text      not null,
    rear    text      not null,
    created timestamp not null default now(),

    primary key (id)
);

create table flcrd.user (
    id                uuid      not null default uuid_generate_v4(),
    external_id       text      not null,
    email             text,
    created           timestamp not null default now(),

    primary key (id)
);
create unique index user_email_idx on flcrd.user(email);
create unique index user_email_idx on flcrd.user(external_id);

create table flcrd.verification_code(
    user_id  uuid      not null,
    code     text      not null,
    code_exp timestamp not null,

    primary key (user_id)
);
create unique index auth_code_idx on flcrd.verification_code(code);

--- TEST DB ---
-- create user test with password 'pass';
-- create database test_flcrd owner test;
-- grant flcrd to test;
-- alter user test with superuser;