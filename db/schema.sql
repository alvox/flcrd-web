create schema authorization flcrd;

create extension "uuid-ossp";

create table flcrd.deck (
    id varchar(40) not null default uuid_generate_v4(),
    name varchar(255) not null,
    description varchar(255) not null default '',
    created timestamp not null default now(),

    primary key (id)
);
create unique index deck_name_idx on flcrd.deck (name);

create table flcrd.flashcard (
    id varchar(40) not null default uuid_generate_v4(),
    deck_id varchar(40) not null references flcrd.deck on delete cascade,
    front varchar(255) not null,
    rear varchar(255) not null,
    created timestamp not null default now()
);

--- TEST DB ---
create user test with password 'pass';
create database test_flcrd owner test;
grant flcrd to test;
alter user test with superuser;