-- gorph UP
CREATE TABLE invoices (
    id INTEGER PRIMARY KEY NOT NULL,
    name VARCHAR(20)
);
-- end --

-- gorph DOWN
DROP TABLE invoices;
-- end --