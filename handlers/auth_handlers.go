package handlers

import (
	"encoding/gob"
	"encoding/json"
	"io"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/oauth2"
)

func init() {
	// Register the map[string]interface{} type with gob
	gob.Register(map[string]interface{}{})
}

func HandleGoogleLogin(config *oauth2.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		url := config.AuthCodeURL("state")
		return c.Redirect(url)
	}
}

func HandleGoogleCallback(config *oauth2.Config, store *session.Store) fiber.Handler {
	log.Println("Entering GoogleCallbackHandler")
	return func(c *fiber.Ctx) error {
		code := c.Query("code")
		if code == "" {
			log.Println("Code is empty")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Code is empty"})
		}
		token, err := config.Exchange(c.Context(), code)
		if err != nil {
			log.Println("Failed to exchange token", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Failed to exchange token"})
		}

		client := config.Client(c.Context(), token)
		resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")

		if err != nil {
			log.Println("Failed to get user info", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get user info"})
		}
		defer resp.Body.Close()

		userInfo, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Failed to read user info", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read user info"})
		}

		var userInfoMap map[string]interface{}
		if err := json.Unmarshal(userInfo, &userInfoMap); err != nil {
			log.Println("Failed to parse user info", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse user info"})
		}

		sess, err := store.Get(c)
		if err != nil {
			log.Println("Failed to get session", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get session"})
		}

		sess.Set("user", userInfoMap)
		if err := sess.Save(); err != nil {
			log.Println("Failed to save session", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save session"})
		}

		return c.Redirect("/") // Redirect to home page after successful login
	}
}

func AuthMiddleware(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			log.Println("Failed to get session", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		user := sess.Get("user")
		if user == nil {
			log.Println("User not found in session")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		return c.Next()
	}
}

func HandleAuthStatus(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"isLoggedIn": false,
				"user":       nil,
			})
		}

		user := sess.Get("user")
		if user == nil {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"isLoggedIn": false,
				"user":       nil,
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"isLoggedIn": true,
			"user":       user,
		})
	}
}

func HandleLogout(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get session"})
		}

		sess.Delete("user")
		if err := sess.Save(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save session"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logged out successfully"})
	}
}
