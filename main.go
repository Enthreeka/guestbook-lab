package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
	"lab-3/app/controller"
	"lab-3/app/controller/handler"
	"lab-3/repository"
	"lab-3/repository/connect"
)

func main() {
	psql, err := connect.New(context.Background(), 5, "postgres://postgres:postgres@db:5432/guestbook")
	if err != nil {
		log.Fatal("failed to connect PostgreSQL: %v", err)
	}

	defer psql.Close()

	store := session.New(session.Config{
		CookieSecure:   true,
		CookieHTTPOnly: true,
	})

	lRepo := repository.NewListRepository(psql)
	sRepo := repository.NewSessionRepository(psql)
	uRepo := repository.NewUserRepository(psql)
	gRepo := repository.NewGuestbookRepository(psql)

	list := handler.NewListHandler(lRepo, uRepo)
	user := handler.NewUserHandler(uRepo, sRepo, store)
	guestbook := handler.NewGuestbookHandler(gRepo, uRepo)

	auth := controller.AuthMiddleware(sRepo, store)

	engine := html.New("./app/view", ".html")
	engine.Debug(true)
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	v1 := app.Group("/")
	v1.Get("/", func(c *fiber.Ctx) error {
		return c.Render("auth", fiber.Map{})
	})
	v1.Post("/register", user.Register)
	v1.Post("/login", user.Login)
	v1.Post("/logout", user.Logout)

	v1.Use(auth).Get("/list", list.Get)
	v1.Use(auth).Post("/list/create", list.Create)
	v1.Use(auth).Post("list/delete", list.Delete)

	v2 := v1.Group("/list")
	v2.Use(auth).Get("/guestbook", guestbook.Get)
	v2.Use(auth).Post("/guestbook/add", guestbook.Create)

	log.Info("Starting http server: localhost:8080")
	if err := app.Listen(fmt.Sprintf(":%d", 8080)); err != nil {
		log.Fatal("Server listening failed:%s", err)
	}
}
