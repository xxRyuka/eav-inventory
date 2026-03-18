create table if not exists categories(
    id serial primary key,
    name varchar(255) not null,
    parent_id integer references categories(id) null
);





create table if not exists category_attributes(
    id serial primary key,
    name varchar(255) not null,
    data_type VARCHAR(50) NOT NULL,
    category_id integer references categories(id) on delete cascade not null,
    is_required boolean default false
);

