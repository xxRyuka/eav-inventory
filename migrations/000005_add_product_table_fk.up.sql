alter table stocks
    add constraint fk_stocks_products
        foreign key (product_id)
            references products (id)
            on delete restrict;

alter table stock_movements
    add constraint fk_stocks_movements_products
        foreign key (product_id)
            references products (id)
            on delete restrict;