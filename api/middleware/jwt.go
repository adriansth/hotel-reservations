package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/adriansth/go-hotel-reservations/db"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()["X-Api-Token"]
		if !ok {
			return fmt.Errorf("Unauthorized")
		}
		claims, err := validateToken(token[0])
		if err != nil {
			return err
		}
		// check token expiration
		expiresFloat, ok := claims["expires"].(float64)
		if !ok {
			return fmt.Errorf("Token invalid: 'expires' claim is missing or not a float64")
		}
		expires := int64(expiresFloat)
		if time.Now().Unix() > expires {
			return fmt.Errorf("Token expired")
		}
		fmt.Println(claims)
		userID, ok := claims["userID"].(string)
		if !ok {
			return fmt.Errorf("Token invalid: 'id' claim is missing or not a string")
		}
		user, err := userStore.GetUserByID(c.Context(), userID)
		if err != nil {
			return fmt.Errorf("Unauthorized")
		}
		// set the current authenticated user to the context
		c.Context().SetUserValue("user", user)
		return c.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("Invalid signing method", token.Header["alg"])
			return nil, fmt.Errorf("Unauthorized")
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("Failed to parse JWT token: ", err)
		return nil, fmt.Errorf("Unauthorized")
	}
	if !token.Valid {
		fmt.Println("Invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Unauthorized")
	}
	return claims, nil
}
