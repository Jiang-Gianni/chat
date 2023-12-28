create table if not exists message(
    id integer primary key,
    room_id integer not null,
    username text not null,
    message text not null,
    sent_at timestamp not null
);