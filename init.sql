CREATE TABLE orders (
                        order_uid VARCHAR(255) PRIMARY KEY,
                        track_number VARCHAR(255) NOT NULL,
                        entry VARCHAR(255) NOT NULL,
                        locale VARCHAR(50) NOT NULL,
                        internal_signature VARCHAR(255),  --no NOT NULL
                        customer_id VARCHAR(255) NOT NULL,
                        delivery_service VARCHAR(255) NOT NULL,
                        shardkey VARCHAR(100) NOT NULL,
                        sm_id INT NOT NULL,
                        date_created TIMESTAMPTZ NOT NULL,
                        oof_shard VARCHAR(50) NOT NULL
);

CREATE TABLE delivery(
                         id BIGSERIAL PRIMARY KEY,
                         order_uid VARCHAR(255) NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
                         name VARCHAR(255) NOT NULL,
                         phone VARCHAR(50) NOT NULL,
                         zip VARCHAR(50) NOT NULL,
                         city VARCHAR(255) NOT NULL,
                         address VARCHAR(255) NOT NULL,
                         region VARCHAR(255) NOT NULL,
                         email VARCHAR(255) NOT NULL,
                         UNIQUE(order_uid)
);

CREATE TABLE payment(
                        order_id VARCHAR(255) NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
                        transaction VARCHAR(255) PRIMARY KEY,
                        request_id VARCHAR(255) NOT NULL,
                        currency VARCHAR(50) NOT NULL,
                        provider VARCHAR(255) NOT NULL,
                        amount INT NOT NULL,
                        payment_dt BIGINT NOT NULL,
                        bank VARCHAR(255) NOT NULL,
                        delivery_cost INT NOT NULL,
                        goods_total INT NOT NULL,
                        custom_fee INT DEFAULT 0 -- maybe set to zero ?
);
CREATE TABLE items(
                      id BIGSERIAL PRIMARY KEY,
                      order_uid VARCHAR(255) NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
                      chrt_id BIGINT NOT NULL,
                      track_number VARCHAR(255) NOT NULL,
                      price INT NOT NULL,
                      rid VARCHAR(255) NOT NULL,
                      name VARCHAR(255) NOT NULL,
                      sale INT DEFAULT 0,
                      size VARCHAR(50) DEFAULT '0',
                      total_price INT NOT NULL,
                      nm_id BIGINT NOT NULL,
                      brand VARCHAR(255), --does not matter?
                      status INT NOT NULL
);