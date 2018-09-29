-- gorph DEPENDS 
users/0001_create_users 
-- end --
-- gorph UP
CREATE TABLE foo (
    invoice_id INTEGER,
    user_id INTEGER,
    FOREIGN KEY(invoice_id) REFERENCES invoices(id),
    FOREIGN KEY(user_id) REFERENCES users(id)
);
-- end --
-- gorph DOWN
DROP TABLE foo;
-- end --