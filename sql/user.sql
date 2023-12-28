create table if not exists user(
    id integer primary key,
    username text not null unique,
    password text not null,
    last_room_id integer not null
);