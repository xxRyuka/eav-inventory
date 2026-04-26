create table if not exists warehouses
(
    id       int generated always as identity,
    name     varchar(255) not null,
    code     varchar(255) not null unique,
    location varchar(255) not null,


    constraint pk_warehouses primary key (id) -- acık olarak belirtiyorum id oldugunu
);


create table if not exists stocks
(
    id                 bigint generated always as identity,-- kendi keyi olmayacak mi ? bence olmalı ama
    warehouse_id       int not null,
    product_id         int not null,
    available_quantity int default 0,
    reserved_quantity  int default 0,


    constraint pk_stocks primary key (id),
    constraint uq_stocks_warehouse_products unique (warehouse_id, product_id), -- aynı depoda aynı üründen iki ayrı satır olmasını engeller
    constraint fk_stocks_warehouse foreign key (warehouse_id) references warehouses (id) on DELETE cascade
);


create table if not exists stock_movements
(
    -- ayni sekilde id'si olmalı mı olmamalı mı ?
    id            bigint generated always as identity,
    warehouse_id  int  not null,
    product_id    int  not null,
    quantity      int  not null,
    movement_type text not null,

    constraint chk_stock_movements_movement_type check (
        movement_type in ('PURCHASE_IN', 'TRANSFER_IN', 'ORDER_OUT', 'TRANSFER_OUT', 'ADJUSTMENT')
        ),
    constraint pk_stock_movements primary key (id),
    constraint fk_stocks_movements_warehouse foreign key (warehouse_id) references warehouses (id)
);

create index if not exists idx_stock_movements_watehouse_product on stock_movements (warehouse_id, product_id);


