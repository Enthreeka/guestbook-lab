package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"lab-3/app/model"
	"lab-3/repository"
	"net/http"
	"strconv"
)

type guestbookHandler struct {
	guestbookRepo repository.GuestbookRepository
	userRepo      repository.UserRepository
}

func NewGuestbookHandler(guestbookRepo repository.GuestbookRepository, userRepo repository.UserRepository) *guestbookHandler {
	return &guestbookHandler{
		guestbookRepo: guestbookRepo,
		userRepo:      userRepo,
	}
}

func (g *guestbookHandler) Get(c *fiber.Ctx) error {
	id := c.Query("list_id")
	if id == "" {
		c.Status(http.StatusBadRequest)
	}
	listID, _ := strconv.Atoi(id)

	token := getUserID(c)
	if token == nil {
		c.Status(http.StatusInternalServerError)
	}

	user, err := g.userRepo.GetByID(context.Background(), token.UserID)
	if err != nil {
		log.Errorf("%v\n", err)
		if errors.Is(err, repository.ErrNoRows) {
			return c.SendStatus(http.StatusNotFound)
		}
		c.Status(http.StatusInternalServerError)
	}

	guestbook, err := g.guestbookRepo.GetAllByListID(context.Background(), listID)
	if err != nil {
		log.Errorf("%v\n", err)
		c.Status(http.StatusInternalServerError)
	}

	return c.Render("guestbook", fiber.Map{
		"guestbook": guestbook,
		"user":      user,
		"list":      listID,
	})
}

func (g *guestbookHandler) Create(c *fiber.Ctx) error {
	message := c.FormValue("message")
	listIDStr := c.FormValue("list_id")
	listID, _ := strconv.Atoi(listIDStr)

	token := getUserID(c)
	if token == nil {
		c.Status(http.StatusInternalServerError)
	}

	guestbook := &model.Guestbook{
		UserID:  token.UserID,
		ListID:  listID,
		Message: message,
	}

	_, err := g.guestbookRepo.Create(context.Background(), guestbook)
	if err != nil {
		log.Errorf("%v\n", err)
		c.Status(http.StatusInternalServerError)
	}

	redirectURL := fmt.Sprintf("/list/guestbook?list_id=%d", listID)
	return c.Redirect(redirectURL, fiber.StatusSeeOther)
}
