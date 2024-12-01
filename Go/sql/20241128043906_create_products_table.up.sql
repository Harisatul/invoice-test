CREATE TABLE IF NOT EXISTS product (
                                       id UUID PRIMARY KEY,
                                       item_name TEXT          NOT NULL,
                                       quantity INT            NOT NULL,
                                       total_cogs BIGINT  DEFAULT 0    NOT NULL,
                                       total_price_sold BIGINT DEFAULT 0 NOT NULL,
                                       created_at TIMESTAMP DEFAULT NOW() NOT NULL,
                                       updated_at TIMESTAMP DEFAULT NOW() NOT NULL,
                                       invoice_number TEXT NOT NULL,
                                       CONSTRAINT fk_invoice FOREIGN KEY (invoice_number) REFERENCES invoice ON DELETE CASCADE
);
