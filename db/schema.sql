create schema authorization flcrd;

create extension "uuid-ossp";

create table flcrd.deck (
    id            varchar(40)  not null default uuid_generate_v4(),
    name          varchar(255) not null,
    description   varchar(255) not null default '',
    created       timestamp    not null default now(),
    created_by    varchar(40)  not null default 'anonymous',
    public        boolean      not null default false,
    search_tokens tsvector,

    primary key (id)
);
create unique index deck_name_user_idx on flcrd.deck(name, created_by);

create table flcrd.flashcard (
    id      varchar(40)  not null default uuid_generate_v4(),
    deck_id varchar(40)  not null references flcrd.deck on delete cascade,
    front   varchar(255) not null,
    rear    varchar(255) not null,
    created timestamp    not null default now(),

    primary key (id)
);

create table flcrd.user (
    id                varchar(40)  not null default uuid_generate_v4(),
    name              varchar(128) not null,
    email             varchar(128) not null,
    password          varchar(255) not null,
    status            varchar(50)  not null,
    refresh_token     varchar(40)  not null default '',
    refresh_token_exp timestamp    not null default now(),
    created           timestamp    not null default now(),

    primary key (id)
);
create unique index user_email_idx on flcrd.user(email);

--- TEST DB ---
create user test with password 'pass';
create database test_flcrd owner test;
grant flcrd to test;
alter user test with superuser;