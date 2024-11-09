-- +goose Up
CREATE TABLE IF NOT EXISTS orders
(
    id                       BIGINT PRIMARY KEY,
    user_id                  BIGINT                                             NOT NULL,
    created_at               TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expiry_date              TIMESTAMP WITH TIME ZONE,
    accept_return_order_date TIMESTAMP                                          NULL,
    returned_from_client_at  TIMESTAMP                                          NULL,
    returned_to_courier_at   TIMESTAMP                                          NULL,
    packaging                TEXT                                               NOT NULL,
    weigh                    float                                              NOT NULL,
    cost                     float                                              NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS orders;
