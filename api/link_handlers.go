package api

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type LinksHandlers struct {
	AppScheme string // e.g. "paysplit"
}

func NewLinksHandlers(appScheme string) *LinksHandlers {
	if appScheme == "" {
		appScheme = "paysplit"
	}
	return &LinksHandlers{AppScheme: appScheme}
}

type settleReq struct {
	ToVPA       string `json:"to_vpa"`       // e.g. apurva@okhdfcbank
	ToName      string `json:"to_name"`      // e.g. Apurva S
	AmountPaise int64  `json:"amount_paise"` // e.g. 12500 for ₹125.00
	Note        string `json:"note"`         // e.g. Goa Trip settle-up
}

func (h *LinksHandlers) HandleBuildSettleLink(c *fiber.Ctx) error {
	var req settleReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json"})
	}
	if req.ToVPA == "" || req.AmountPaise <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "to_vpa and amount_paise required"})
	}
	// super-basic VPA sanity (user@psp). Keep it loose to avoid false negatives.
	if !strings.Contains(req.ToVPA, "@") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid vpa"})
	}

	// UPI amount is in rupees with 2 decimals; we store paise
	rupees := float64(req.AmountPaise) / 100.0

	// Build UPI deep link: upi://pay?pa=<vpa>&pn=<name>&am=<rupees>&cu=INR&tn=<note>
	upi := fmt.Sprintf(
		"upi://pay?pa=%s&pn=%s&am=%.2f&cu=INR&tn=%s",
		url.QueryEscape(req.ToVPA),
		url.QueryEscape(req.ToName),
		rupees,
		url.QueryEscape(req.Note),
	)

	// Your app deep link to prefill a "Settle" screen
	app := fmt.Sprintf(
		"%s://settle?vpa=%s&name=%s&amount=%d&note=%s",
		h.AppScheme,
		url.QueryEscape(req.ToVPA),
		url.QueryEscape(req.ToName),
		req.AmountPaise,
		url.QueryEscape(req.Note),
	)

	return c.JSON(fiber.Map{
		"upi": upi, // open directly to payments apps (GPay/PhonePe/Paytm/BHIM)
		"app": app, // opens your app’s Settle screen (if you register the scheme)
	})
}
