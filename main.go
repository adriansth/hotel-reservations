package main

import (
	"context"
	"flag"
	"log"

	"github.com/adriansth/go-hotel-reservations/api"
	"github.com/adriansth/go-hotel-reservations/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	// database
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(db.DBURI).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}
	// server
	listenAddr := flag.String("listenAddr", ":8080", "The listen address of the API server.")
	flag.Parse()
	// handlers initialization
	var (
		app          = fiber.New(config)
		apiv1        = app.Group("/api/v1")
		userHandler  = api.NewUserHandler(db.NewMongoUserStore(client, db.DBNAME))
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		hotelHandler = api.NewHotelHandler(hotelStore, roomStore)
	)
	// user handlers
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	// hotel handlers
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	// initialization
	app.Listen(*listenAddr)
}
