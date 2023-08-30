create table slugs
(
    id   serial primary key,
    name varchar(255) not null unique
);

create table operations
(
    id serial primary key,
    user_id int not null,
    slug_id int not null,
    created_at timestamp not null default now(),
    removed_at timestamp not null
);
