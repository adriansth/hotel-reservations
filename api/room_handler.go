package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/adriansth/go-hotel-reservations/db"
	"github.com/adriansth/go-hotel-reservations/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookRoomParams struct {
	FromDate  time.Time `json:"fromDate"`
	TillDate  time.Time `json:"tillDate"`
	NumPeople int       `json:"numPeople"`
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
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
	booking := types.Booking{
		UserID:    user.ID,
		RoomID:    roomID,
		FromDate:  params.FromDate,
		TillDate:  params.TillDate,
		NumPeople: params.NumPeople,
	}
	fmt.Printf("%+v\n", booking)
	return nil
}
