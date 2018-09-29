-- gorph UP
CREATE TABLE users (
    ID INTEGER PRIMARY KEY NOT NULL,
    username VARCHAR(32)
);
-- end --
-- gorph DOWN
DROP TABLE users;
-- end --