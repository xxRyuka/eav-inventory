drop table if exists product_attribute_values;
drop table if exists category_attributes;

create table if not exists attributes
(
    id        serial primary key,
    code      varchar(255) not null unique,
    name      varchar(255) not null,
    data_type varchar(255) not null
);

create table if not exists category_attributes
(
    attribute_id integer not null references attributes (id) on delete cascade,
    category_id  integer not null references categories (id) on DELETE cascade,
    is_required  boolean not null default false,
    PRIMARY KEY (category_id, attribute_id)
);

create table if not exists product_attribute_values
(
    attribute_id integer not null references attributes (id) on delete cascade,
    product_id   integer not null references products (id),
    value        text    not null,
    primary key (attribute_id, product_id)
);