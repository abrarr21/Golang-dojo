package database

import (
	"context"
	"log"
	"time"

	"github.com/abrarr21/golang-auth/internal/config"
	"github.com/abrarr21/golang-auth/internal/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Database struct {
	client *mongo.Client
	DB     *mongo.Database
	Users  *mongo.Collection
}

func ConnectDB(cfg *config.DatabaseConfig) *Database {
	c, err := mongo.Connect(options.Client().ApplyURI(cfg.MongoDB_URI))
	if err != nil {
		log.Fatal("failed to connect to MongoDB: ", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := c.Ping(ctx, nil); err != nil {
		log.Fatal("failed to ping MongoDB: ", err)
	}

	log.Println("Connected to MongoDB ✅")

	db := c.Database(cfg.DBName)
	users := db.Collection("users")

	if err := models.EnsureIndexes(users); err != nil {
		log.Fatal("failed to create indexes: ", err)
	}

	return &Database{
		client: c,
		DB:     db,
		Users:  users,
	}
}

func (d *Database) Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := d.client.Disconnect(ctx); err != nil {
		log.Printf("failed to connect from MongoDB: %v", err)
	}

	log.Println("Disconnected from MongoDB")
}
