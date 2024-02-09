package main

import (
	"context"
	"flag"
	"log"

	"github.com/adriansth/go-hotel-reservations/api"
	"github.com/adriansth/go-hotel-reservations/api/middleware"
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
		app  = fiber.New(config)
		auth = app.Group("/api")
		// create handlers
		hotelStore   = db.NewMongoHotelStore(client)
		userStore    = db.NewMongoUserStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		bookingStore = db.NewMongoBookingStore(client)
		authHandler  = api.NewAuthHandler(userStore)
		store        = &db.Store{
			Hotel:   hotelStore,
			Room:    roomStore,
			User:    userStore,
			Booking: bookingStore,
		}
		apiv1          = app.Group("/api/v1", middleware.JWTAuthentication(userStore))
		admin          = apiv1.Group("/admin", middleware.AdminAuth)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
		userHandler    = api.NewUserHandler(userStore)
		hotelHandler   = api.NewHotelHandler(store)
	)
	// auth
	auth.Post("/auth", authHandler.HandleAuthenticate)
	// user handlers
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	// hotel handlers
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)
	// room handlers
	apiv1.Post("/room", roomHandler.HandleGetRooms)
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)
	// bookings handlers
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	// admin handlers
	admin.Get("/booking", bookingHandler.HandleGetBookings)

	// initialization
	app.Listen(*listenAddr)
}
