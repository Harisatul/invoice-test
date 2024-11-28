CREATE TABLE IF NOT EXISTS invoice (
                                       invoice_number TEXT     NOT NULL PRIMARY KEY,
                                       date DATE NOT NULL,
                                       customer_name TEXT      NOT NULL ,
                                       salesperson TEXT        NOT NULL,
                                       notes TEXT,
                                       created_at TIMESTAMP DEFAULT NOW() NOT NULL,
                                       updated_at TIMESTAMP DEFAULT NOW() NOT NULL,
                                       payment_type PAYMENT_STATUS     NOT NULL
);
