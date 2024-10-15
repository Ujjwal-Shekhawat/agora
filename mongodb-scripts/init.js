
db = db.getSiblingDB('main_dab');


const pwd = process.env.MONGO_INITDB_ROOT_PASSWORD

db.createUser({
    user: "user",
    pwd: pwd,
    roles: [{ role: "readWrite", db: "user_service" }]
});

db.createCollection("users", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: ["username", "email", "password"],
            additionalProperties: false,
            properties: {
                username: {
                    bsonType: "string",
                    description: "must be a string and is required"
                },
                email: {
                    bsonType: "string",
                    description: "must be a string and is required"
                },
                password: {
                    bsonType: "string",
                    description: "must be a string and is required"
                },
                guildNames: {
                    bsonType: "array",
                    items: {
                        bsonType: "string",
                        description: "must be an array of GuildNames representing guilds the user belongs to"
                    }
                }
            }
        }
    }
});

db.users.createIndex(
    { "username": 1 },
    { unique: true, name: "username_unique_index" }
);

db.users.createIndex(
    { "email": 1 },
    { unique: true, name: "email_unique_index" }
);

db.users.createIndex(
    { "guildIds": 1 },
    { name: "guild_pointer_index" }
);

db.createCollection("guilds", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: ["name"],
            additionalProperties: false,
            properties: {
                name: {
                    bsonType: "string",
                    description: "must be a string and is required"
                },
                channels: {
                    bsonType: "array",
                    items: {
                        bsonType: "string",
                        description: "must be an array of strings representing channel names"
                    },
                    description: "must be an array"
                },
                members: {
                    bsonType: "array",
                    items: {
                        bsonType: "string",
                        description: "must be an array of strings representing user names"
                    },
                    description: "must be an array"
                }
            }
        }
    }
});

db.guilds.createIndex(
    { "name": 1 },
    { unique: true, name: "guild_name_unique_index" }
);
