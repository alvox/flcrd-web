-- SCHEMA --
create schema authorization flcrd;
create extension "uuid-ossp";

create table flcrd.deck (
    id          uuid      not null default uuid_generate_v4(),
    name        text      not null,
    description text      not null default '',
    created     timestamp not null default now(),
    created_by  uuid      not null,
    public      boolean   not null default false,
    search_tokens tsvector,

    primary key (id)
);
create unique index deck_name_user_idx on flcrd.deck(name, created_by);

create table flcrd.flashcard (
    id         uuid      not null default uuid_generate_v4(),
    deck_id    uuid      not null references flcrd.deck on delete cascade,
    front      text      not null,
    front_type text      not null,
    rear       text      not null,
    rear_type  text      not null,
    created    timestamp not null default now(),

    primary key (id)
);

create table flcrd.user (
    id      uuid      not null default uuid_generate_v4(),
    name    text      not null,
    email   text      not null,
    status  text      not null,
    created timestamp not null default now(),

    primary key (id)
);
create unique index user_email_idx on flcrd.user(email);

create table flcrd.verification_code(
    user_id  uuid      not null,
    code     text      not null,
    code_exp timestamp not null,

    primary key (user_id)
);
create unique index auth_code_idx on flcrd.verification_code(code);

create table flcrd.credentials(
    user_id           uuid not null,
    password          text not null,
    refresh_token     text not null default '',
    refresh_token_exp timestamp not null default now(),

    primary key (user_id)
);


-- DATA --
insert into flcrd.user (id, name, email, status, created)
values ('40AFBC9A-27E3-4B38-97F9-2930B8790A9F', 'Testuser1', 'testuser1@example.com', 'ACTIVE', '2019-01-01 09:00:00'),
       ('DD4A5E3A-4D95-44C1-8AA3-E29FA9A29570', 'Testuser2', 'testuser2@example.com', 'PENDING', '2019-01-01 12:00:00');

insert into flcrd.credentials (user_id, password, refresh_token, refresh_token_exp)
values ('40AFBC9A-27E3-4B38-97F9-2930B8790A9F', '12345', 'refreshtoken1', '2019-02-02 10:00:00'),
       ('DD4A5E3A-4D95-44C1-8AA3-E29FA9A29570', '54321', 'refreshtoken2', '2019-03-03 10:00:00');

insert into flcrd.deck (id, name, description, created, created_by, public, search_tokens)
values ('9F2556FB-0B84-4B8D-AB0A-B5ACB0C89F6E', 'Test Name 1', 'Test Description 1', '2019-01-01 10:00:00', '40AFBC9A-27E3-4B38-97F9-2930B8790A9F', false, to_tsvector('Test Name 1 Test Description 1')),
       ('2601DA50-56A6-41A1-A92E-5624598A7D19', 'Test Name 2', 'Test Description 2', '2019-02-02 12:22:00', 'DD4A5E3A-4D95-44C1-8AA3-E29FA9A29570', true, to_tsvector('Test Name 2 Test Description 2')),
       ('4735BDB2-45A5-42D9-A6D4-6DB29787B5F1', 'Test Name 3', 'Test Description 3', '2019-03-03 12:22:00', '40AFBC9A-27E3-4B38-97F9-2930B8790A9F', true, to_tsvector('Test Name 3 Test Description 3'));

insert into flcrd.flashcard (id, deck_id, front, front_type, rear, rear_type, created)
values ('9F814806-E2DF-4598-A323-1380D47B9C35', '9F2556FB-0B84-4B8D-AB0A-B5ACB0C89F6E', 'Test Front 1 1', 'TEXT', 'Test Rear 1 1', 'TEXT', '2019-01-01 10:00:00'),
       ('1CFE647E-3D5F-4C6E-B935-281EFC61E5C9', '9F2556FB-0B84-4B8D-AB0A-B5ACB0C89F6E', 'Test Front 1 2', 'TEXT', 'Test Rear 1 2', 'TEXT', now()),
       ('4686D1B1-6776-4C25-947A-E1285E1B538F', '9F2556FB-0B84-4B8D-AB0A-B5ACB0C89F6E', 'Test Front 1 3', 'TEXT', 'Test Rear 1 3', 'TEXT', now()),
       ('CE233F42-4D4C-4E67-B70B-2B5735916687', '2601DA50-56A6-41A1-A92E-5624598A7D19', 'Test Front 2 1', 'TEXT', 'Test Rear 2 1', 'TEXT', now()),
       ('A5F65E8D-CA03-4D3E-978E-F5A612881231', '2601DA50-56A6-41A1-A92E-5624598A7D19', 'Test Front 2 2', 'TEXT', 'https://s3/testuser/testdeck/testimg.jpeg', 'IMAGE_URL', '2019-05-05 15:55:00');

insert into flcrd.verification_code (user_id, code, code_exp)
values ('DD4A5E3A-4D95-44C1-8AA3-E29FA9A29570', 'code_for_user_2', '2019-01-01 09:00:00');