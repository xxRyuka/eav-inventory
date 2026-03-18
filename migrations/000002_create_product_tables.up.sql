create table if not exists products
(
    id          serial primary key,
    category_id integer      not null references categories (id),
    name        varchar(255) not null,
    sku         VARCHAR(100) NOT NULL UNIQUE
);


create table if not exists product_attribute_values
(
    id         serial primary key,
    product_id integer not null references products (id) on delete cascade,
    category_attribute_id integer not null references category_attributes(id) on delete cascade,
    value text not null 

);