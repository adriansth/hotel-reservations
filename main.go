package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/adriansth/go-hotel-reservations/api"
	"github.com/adriansth/go-hotel-reservations/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dburi = "mongodb://localhost:27017"
const dbname = "hotel-reservation"
const userColl = "users"

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	// database
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	uri := os.Getenv("MONGO_URI")
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}
	// server
	listenAddr := flag.String("listenAddr", ":8080", "The listen address of the API server.")
	flag.Parse()
	app := fiber.New(config)
	apiv1 := app.Group("/api/v1")
	// handlers initialization
	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))
	// routes
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	// initialization
	app.Listen(*listenAddr)
}
