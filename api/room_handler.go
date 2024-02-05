package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/adriansth/go-hotel-reservations/db"
	"github.com/adriansth/go-hotel-reservations/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookRoomParams struct {
	FromDate  time.Time `json:"fromDate"`
	TillDate  time.Time `json:"tillDate"`
	NumPeople int       `json:"numPeople"`
}

func (p BookRoomParams) validate() error {
	now := time.Now()
	if now.After(p.FromDate) || now.After(p.TillDate) || p.FromDate.After(p.TillDate) {
		return fmt.Errorf("Cannot book a room in the past")
	}
	return nil
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Room.GetRooms(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if err := params.validate(); err != nil {
		return err
	}
	roomID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(genericResponse{
			Type: "error",
			Msg:  "Internal server error",
		})
	}
	ok, err = h.isRoomAvailable(c.Context(), roomID, params)
	if err != nil {
		return err
	}
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(genericResponse{
			Type: "error",
			Msg:  fmt.Sprintf("Room %s has already been booked", c.Params("id")),
		})
	}
	booking := types.Booking{
		UserID:    user.ID,
		RoomID:    roomID,
		FromDate:  params.FromDate,
		TillDate:  params.TillDate,
		NumPeople: params.NumPeople,
	}
	inserted, err := h.store.Booking.InsertBooking(c.Context(), &booking)
	if err != nil {
		return err
	}
	return c.JSON(inserted)
}

func (h *RoomHandler) isRoomAvailable(ctx context.Context, roomID primitive.ObjectID, params BookRoomParams) (bool, error) {
	where := bson.M{
		"roomID": roomID,
		"fromDate": bson.M{
			"$gte": params.FromDate,
		},
		"tillDate": bson.M{
			"$lte": params.TillDate,
		},
	}
	bookings, err := h.store.Booking.GetBookings(ctx, where)
	if err != nil {
		return false, err
	}
	ok := len(bookings) == 0
	return ok, nil
}
