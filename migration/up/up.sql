
create table "user"(
  id int generated always as identity,
  login varchar(100) unique,
  password varchar(100) unique,
  primary key (id)
);

create table list(
    id int generated always as identity,
    name text,
    primary key (id)
);

create table guestbook(
    id int generated always as identity,
    message text,
    created_at timestamp default current_timestamp not null,
    user_id int,
    list_id int,
    primary key (id),
    foreign key (user_id)
            references "user" (id) on delete cascade ,
    foreign key (list_id)
        references list (id) on delete cascade
);

CREATE TABLE IF NOT EXISTS session(
    id int generated always as identity,
    token uuid not null,
    user_id int,
    primary key (id),
    foreign key (user_id)
        references "user" (id) on delete cascade
);
