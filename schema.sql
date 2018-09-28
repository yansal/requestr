create table pools (
    id serial primary key
);

create table jobs (
    id serial primary key,
    name text,
    pool_id int references pools(id)
)