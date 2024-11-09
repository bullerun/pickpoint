-- +goose NO TRANSACTION
-- +goose Up
create index concurrently if not exists idx_orders_user_id  on orders using HASH (user_id);



-- +goose Down
drop index concurrently if exists idx_orders_user_id ;

