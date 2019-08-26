-- SCHEMA --
create schema authorization flcrd;
create extension "uuid-ossp";

create table flcrd.deck (
    id varchar(40) not null default uuid_generate_v4(),
    name varchar(255) not null,
    description varchar(255) not null default '',
    created timestamp not null default now(),
    created_by varchar(40) not null default 'anonymous',
    private boolean not null default true,

    primary key (id)
);
create unique index deck_name_user_idx on flcrd.deck (name, created_by);

create table flcrd.flashcard (
    id varchar(40) not null default uuid_generate_v4(),
    deck_id varchar(40) not null references flcrd.deck on delete cascade,
    front varchar(255) not null,
    rear varchar(255) not null,
    created timestamp not null default now(),

    primary key (id)
);

create table flcrd.user (
    id varchar(40) not null default uuid_generate_v4(),
    name varchar(128) not null,
    email varchar(128) not null,
    password varchar(255) not null,
    created timestamp not null default now(),

    primary key (id)
);
create unique index user_email_idx on flcrd.user (email);

-- DATA --
insert into flcrd.deck (id, name, description, created, private) values
('test_deck_id_1', 'Test Name 1', 'Test Description 1', '2019-01-01 10:00:00', true),
('test_deck_id_2', 'Test Name 2', 'Test Description 2', '2019-02-02 12:22:00', false);

insert into flcrd.flashcard (id, deck_id, front, rear, created) values
('test_flashcard_id_1', 'test_deck_id_1', 'Test Front 1 1', 'Test Rear 1 1', '2019-01-01 10:00:00'),
('test_flashcard_id_2', 'test_deck_id_1', 'Test Front 1 2', 'Test Rear 1 2', now()),
('test_flashcard_id_3', 'test_deck_id_1', 'Test Front 1 3', 'Test Rear 1 3', now()),
('test_flashcard_id_4', 'test_deck_id_2', 'Test Front 2 1', 'Test Rear 2 1', now()),
('test_flashcard_id_5', 'test_deck_id_2', 'Test Front 2 2', 'Test Rear 2 2', '2019-05-05 15:55:00');

insert into flcrd.user (id, name, email, password) values
('testuser_id_1', 'Testuser', 'testuser@example.com', '12345');
