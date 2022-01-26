-- gorph UP
CREATE TABLE embedded_migrations (
    id INTEGER PRIMARY KEY NOT NULL,
    name VARCHAR(20)
);
-- end --

-- gorph DOWN
DROP TABLE embedded_migrations;
-- end --