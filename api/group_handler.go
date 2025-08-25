package api

import (
	"net/http"

	"github.com/akarshgo/paysplit/db"
	"github.com/akarshgo/paysplit/types"
	"github.com/gofiber/fiber/v2"
)

type GroupHandlers struct {
	groups db.GroupStore
}

func NewGroupHanlders(groups db.GroupStore) *GroupHandlers {
	return &GroupHandlers{
		groups: groups,
	}
}

type createGroupReq struct {
	Name      string `json:"name"`
	CreatedBy string `json:"created_by"`
}

func (h *GroupHandlers) HandleCreateGroup(c *fiber.Ctx) error {
	var req createGroupReq
	if err := c.BodyParser(&req); err != nil || req.Name == "" || req.CreatedBy == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "bad request"})
	}
	id, err := h.groups.Create(c.Context(), &types.Group{Name: req.Name, CreatedBy: req.CreatedBy})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create"})
	}
	_ = h.groups.AddMember(c.Context(), id, req.CreatedBy) // creator joins
	return c.Status(201).JSON(fiber.Map{"id": id})
}

type addMemberReq struct {
	UserID string `json:"user_id"`
}

func (h *GroupHandlers) HandleAddMember(c *fiber.Ctx) error {
	gid := c.Params("id")
	var req addMemberReq
	if err := c.BodyParser(&req); err != nil || gid == "" || req.UserID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "bad request"})
	}
	if err := h.groups.AddMember(c.Context(), gid, req.UserID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to add"})
	}
	return c.SendStatus(204)
}

func (h *GroupHandlers) HandleListGroups(c *fiber.Ctx) error {
	// If your interface supports ListByUser(userID), you can pass a query param later.
	out, err := h.groups.ListByUser(c.Context(), "")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to list"})
	}
	return c.JSON(out)
}
