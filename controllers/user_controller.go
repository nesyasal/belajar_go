package controllers

import (
	"context"
	"todo-app/database"
	"todo-app/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllUsers(c *fiber.Ctx) error {
	collection := database.DB.Collection("users")

	// Ambil semua dokumen
	cursor, err := collection.Find(context.Background(), bson.M{}, options.Find())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data pengguna"})
	}
	defer cursor.Close(context.Background())

	var users []models.UserResponse
	for cursor.Next(context.Background()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mendekode data pengguna"})
		}

		users = append(users, models.UserResponse{
			ID:    user.ID.Hex(),
			Name:  user.Name,
			Email: user.Email,
		})
	}

	return c.JSON(users)
}
