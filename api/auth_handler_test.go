package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/adriansth/go-hotel-reservations/db"
	"github.com/adriansth/go-hotel-reservations/types"
	"github.com/gofiber/fiber/v2"
)

func insertTestUser(t *testing.T, userStore db.UserStore) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     "james@foo.com",
		FirstName: "James",
		LastName:  "Foo",
		Password:  "supersecurepassword",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = userStore.InsertUser(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}
	return user
}

func TestAuthenticateWithWrongPassword(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/", authHandler.HandleAuthenticate)
	params := AuthParams{
		Email:    "james@foo.com",
		Password: "wrongsupersecurepassword",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected http status of 400 but got %d", resp.StatusCode)
	}
	var genResp genericResponse
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		t.Fatal(err)
	}
	if genResp.Type != "error" {
		t.Fatalf("Expected generic response type to be <error> but got %s", genResp.Type)
	}
	if genResp.Msg != "Invalid credentials" {
		t.Fatalf("Expected generic response message to be <Invalid credentials> but got %s", genResp.Msg)
	}
}

func TestAuthenticateSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	insertedUser := insertTestUser(t, tdb.UserStore)
	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/", authHandler.HandleAuthenticate)
	params := AuthParams{
		Email:    "james@foo.com",
		Password: "supersecurepassword",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected http status of 200 but got %d", resp.StatusCode)
	}
	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatal(err)
	}
	if authResp.Token == "" {
		t.Fatalf("Expected the jwt token to be present in the auth response.")
	}
	// set the encrypted password to an empty string, because we do not return that in any json response
	insertedUser.EncryptedPassword = ""
	if reflect.DeepEqual(insertedUser, authResp.User) {
		t.Fatalf("Expected the user to be the inserted user.")
	}
}
