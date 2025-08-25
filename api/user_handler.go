package api

import (
	"errors"
	"strconv"
	"strings"

	"github.com/akarshgo/paysplit/db"
	"github.com/akarshgo/paysplit/types"
	"github.com/gofiber/fiber/v2"
)

type UserHandlers struct {
	userStore db.UserStore
}

func NewUserHandlers(userStore db.UserStore) *UserHandlers {
	return &UserHandlers{userStore: userStore}
}

// POST /users
func (h *UserHandlers) HandleCreateUser(c *fiber.Ctx) error {
	var user types.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON body"})
	}
	user.Name = strings.TrimSpace(user.Name)
	if user.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name is required"})
	}

	// store.Create should set user.ID (and CreatedAt) if your PG impl uses RETURNING id
	if err := h.userStore.Create(c.Context(), &user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create user"})
	}
	return c.Status(fiber.StatusCreated).JSON(user)
}

// GET /users/:id
func (h *UserHandlers) HandleGetUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id is required"})
	}
	user, err := h.userStore.GetByID(c.Context(), userID)
	if err != nil {
		// Treat not found generically (your store can return sql.ErrNoRows; map to 404)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

// GET /users?email=&phone=&q=&limit=&offset=
func (h *UserHandlers) HandleGetUsers(c *fiber.Ctx) error {
	var f db.UserFilter

	if v := strings.TrimSpace(c.Query("email")); v != "" {
		f.Email = &v
	}
	if v := strings.TrimSpace(c.Query("phone")); v != "" {
		f.Phone = &v
	}
	if v := strings.TrimSpace(c.Query("q")); v != "" {
		f.Query = &v
	}

	limit, offset := parseLimitOffset(c.Query("limit"), c.Query("offset"))

	users, err := h.userStore.Find(c.Context(), f, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list users"})
	}
	return c.Status(fiber.StatusOK).JSON(users)
}

// PATCH /users/:id
func (h *UserHandlers) HandleUpdateUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id is required"})
	}

	var body types.User
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON body"})
	}

	// Optional: ensure the user exists first (nice DX)
	if _, err := h.userStore.GetByID(c.Context(), userID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	// Build the update model
	up := &types.User{
		ID:    userID,
		Name:  strings.TrimSpace(body.Name),
		Email: body.Email,
		Phone: body.Phone,
		UPI:   body.UPI,
	}
	if up.Name == "" {
		// allow partial? If name is empty and you want to keep old, you could re-fetch; here we enforce non-empty
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name is required"})
	}

	if err := h.userStore.Update(c.Context(), up); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update user"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "user updated", "id": userID})
}

// DELETE /users/:id
func (h *UserHandlers) HandleDeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id is required"})
	}
	if err := h.userStore.Delete(c.Context(), userID); err != nil {
		// If your store returns sql.ErrNoRows on missing, map to 404
		if errors.Is(err, fiber.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete user"})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// ---- helpers ----

func parseLimitOffset(limStr, offStr string) (int, int) {
	limit := 50
	offset := 0
	if limStr != "" {
		if v, err := strconv.Atoi(limStr); err == nil && v > 0 && v <= 200 {
			limit = v
		}
	}
	if offStr != "" {
		if v, err := strconv.Atoi(offStr); err == nil && v >= 0 {
			offset = v
		}
	}
	return limit, offset
}
