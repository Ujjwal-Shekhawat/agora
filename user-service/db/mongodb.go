package db

import (
	"context"
	"log"
	"time"
	"user_service/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoSession struct {
	client *mongo.Client
}

var mongodb *mongoSession = &mongoSession{
	client: nil,
}

func InitMongoSess() error {
	cfg := config.LoadConfig()
	log.Println("Init Mongo Session Now")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.MongoURI).SetMinPoolSize(10).SetMaxPoolSize(20).SetMaxConnIdleTime(10 * time.Second)

	var err error
	mongodb.client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Println("Mongo url", cfg.MongoURI)
		log.Fatal("Error connecting to mongo client", err)
	}

	err = mongodb.client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Could not connect to MongoDB: ", err)
	}

	log.Println("Connected to MongoDB")

	database := mongodb.client.Database("users")

	collection := database.Collection("user")

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal("Something went wrong while reading colection")
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var doc map[string]interface{}
		if err := cursor.Decode(&doc); err != nil {
			log.Fatal("Error decoding cursor")
		}

		log.Println(doc)
	}

	return nil
}

func MongoCreateUser(name, email, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user := bson.D{
		{Key: "username", Value: name},
		{Key: "email", Value: email},
		{Key: "password", Value: password},
	}

	_, err := mongodb.client.Database("main_dab").Collection("users").InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func MongoGetUser(name string) map[string]interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res := mongodb.client.Database("main_dab").Collection("users").FindOne(ctx, bson.D{{Key: "username", Value: name}})

	var dRes map[string]interface{}
	res.Decode(&dRes)

	return dRes
}

func MongoCreateGuild(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	guild := &bson.D{
		{Key: "name", Value: name},
		{Key: "channels", Value: bson.A{"general"}},
	}

	_, err := mongodb.client.Database("main_dab").Collection("guilds").InsertOne(ctx, guild)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func MongoGetGuild(name string) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res := mongodb.client.Database("main_dab").Collection("guilds").FindOne(ctx, bson.D{{Key: "name", Value: name}})

	var rDec map[string]interface{}

	if err := res.Decode(rDec); err != nil {
		log.Println(err)
		return nil, err
	}

	return rDec, nil
}
