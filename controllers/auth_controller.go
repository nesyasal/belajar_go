package controllers

import (
	"context"
	"todo-app/database"
	"todo-app/models"

	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Data tidak valid"})
	}

	collection := database.DB.Collection("users")
	_, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan pengguna"})
	}

	return c.JSON(fiber.Map{"message": "Registrasi berhasil"})
}

func Login(c *fiber.Ctx) error {
	var input models.User
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Data tidak valid"})
	}

	var user models.User
	collection := database.DB.Collection("users")
	err := collection.FindOne(context.Background(), fiber.Map{"email": input.Email, "password": input.Password}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Email atau password salah"})
	}

	return c.JSON(fiber.Map{"message": "Login berhasil", "user": user})
}
