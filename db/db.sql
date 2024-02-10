CREATE TABLE orders
(
    order_uid          varchar(255) NOT NULL PRIMARY KEY,
    track_number       varchar(50),
    entry              varchar(50),
    locale             varchar(10),
    internal_signature varchar(100),
    customer_id        varchar(50),
    delivery_service   varchar(50),
    shardkey           varchar(10),
    sm_id              int4,
    date_created       timestamp,
    oof_shard          varchar(10)
);

CREATE TABLE delivery
(
    id       serial4 PRIMARY KEY,
    order_id varchar(255) REFERENCES orders (order_uid) ON DELETE CASCADE NOT NULL,
    name     varchar(100),
    phone    varchar(20),
    zip      varchar(20),
    city     varchar(100),
    address  varchar(200),
    region   varchar(100),
    email    varchar(100)
);

CREATE TABLE payment
(
    id             serial4 PRIMARY KEY,
    transaction_id varchar(50),
    order_id       varchar(255) REFERENCES orders (order_uid) ON DELETE CASCADE,
    request_id     varchar(50),
    currency       varchar(5),
    provider       varchar(50),
    amount         int4,
    payment_dt     int4,
    bank           varchar(50),
    delivery_cost  int4,
    goods_total    int4,
    custom_fee     int4
);

CREATE TABLE items
(
    id           serial4 PRIMARY KEY,
    order_id     varchar(255) REFERENCES orders (order_uid) ON DELETE CASCADE NOT NULL,
    chrt_id      int4,
    track_number varchar(50),
    price        int4,
    rid          varchar(50),
    name         varchar(100),
    sale         int4,
    size         varchar(20),
    total_price  int4,
    nm_id        int4,
    brand        varchar(100),
    status       int4
);

