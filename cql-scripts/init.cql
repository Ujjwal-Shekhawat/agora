CREATE KEYSPACE IF NOT EXISTS users WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 3 };

USE users;

CREATE TABLE IF NOT EXISTS users (
    user_id UUID,
    user_name TEXT,
    user_email TEXT,
    PRIMARY KEY ((user_name), user_email)
);
