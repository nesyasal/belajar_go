package controllers

import (
	"context"
	"todo-app/database"
	"todo-app/models"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Data tidak valid"})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengenkripsi password"})
	}
	user.Password = string(hashedPassword)

	collection := database.DB.Collection("users")
	_, err = collection.InsertOne(context.Background(), user)
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

	// Cari user berdasarkan email saja
	err := collection.FindOne(context.Background(), bson.M{"email": input.Email}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Email atau password salah"})
	}

	// Cek password menggunakan bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Email atau password salah"})
	}

	return c.JSON(fiber.Map{
		"message": "Login berhasil",
		"name":    user.Name,
		"email":   user.Email,
	})
}

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
