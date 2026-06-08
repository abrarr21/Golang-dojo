package database

import (
	"context"
	"log"
	"time"

	"github.com/abrarr21/auth-practice-3/internal/config"
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
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := c.Ping(ctx, nil); err != nil {
		log.Printf("failed to Ping MongoDB")
	}

	log.Println("Connected to MongoDB ✅")

	db := c.Database(cfg.DBName)
	users := db.Collection("users")

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
		log.Printf("failed to Disconnect from MongoDB: %v", err)
	}
	log.Println("Disconnected from MongoDB")
}
