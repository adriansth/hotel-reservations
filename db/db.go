package db

import "os"

const (
	DBNAME     = "hotel-reservation"
	TestDBNAME = "hotel-reservation-test"
)

var DBURI = os.Getenv("MONGO_URI")

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}
