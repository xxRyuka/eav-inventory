alter table stocks
drop
constraint if exists fk_stocks_products;

alter table stock_movements
drop
constraint if exists  fk_stocks_movements_products;