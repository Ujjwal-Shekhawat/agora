CREATE KEYSPACE IF NOT EXISTS guild_messages WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 3 };

USE guild_messages;

CREATE TABLE IF NOT EXISTS messages (
    guild_name TEXT,
    user_name TEXT,
    channel_name TEXT,
    user_message TEXT,
    timestamp TIMESTAMP,
    PRIMARY KEY ((guild_name, channel_name), timestamp)
) WITH CLUSTERING ORDER BY (timestamp DESC);
