package handler

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"lab-3/app/model"
	"lab-3/repository"
	"net/http"
	"strconv"
)

type listHandler struct {
	listRepo repository.ListRepository
	userRepo repository.UserRepository
}

func NewListHandler(listRepo repository.ListRepository, userRepo repository.UserRepository) *listHandler {
	return &listHandler{
		listRepo: listRepo,
		userRepo: userRepo,
	}
}

func (l *listHandler) Get(c *fiber.Ctx) error {
	token := getUserID(c)
	if token == nil {
		c.Status(http.StatusInternalServerError)
	}

	user, err := l.userRepo.GetByID(context.Background(), token.UserID)
	if err != nil {
		log.Errorf("%v\n", err)
		c.Status(http.StatusInternalServerError)
	}

	list, err := l.listRepo.GetAll(context.Background())
	if err != nil {
		log.Errorf("%v\n", err)
		c.Status(http.StatusInternalServerError)
	}

	return c.Render("index", fiber.Map{
		"list": list,
		"user": user,
	})
}

func (l *listHandler) Create(c *fiber.Ctx) error {
	name := c.FormValue("name")

	list := &model.List{
		Name: name,
	}
	_, err := l.listRepo.Create(context.Background(), list)
	if err != nil {
		log.Errorf("%v\n", err)
		c.Status(http.StatusInternalServerError)
	}

	return c.Redirect("/list", http.StatusSeeOther)
}

func (l *listHandler) Delete(c *fiber.Ctx) error {
	listIDStr := c.FormValue("list_id")
	listID, _ := strconv.Atoi(listIDStr)

	err := l.listRepo.DeleteByID(context.Background(), listID)
	if err != nil {
		log.Errorf("%v\n", err)
		c.Status(http.StatusInternalServerError)
	}

	return c.Redirect("/list", http.StatusSeeOther)
}

func getUserID(c *fiber.Ctx) *model.Session {
	t := c.Locals("token")

	token := t.(*model.Session)

	return token
}
