package handler

import (
	"context"
	"crypto/sha256"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"lab-3/app/model"
	"lab-3/repository"
	"net/http"
	"time"
)

type userHandler struct {
	userRepo repository.UserRepository
	sessRepo repository.SessionRepository
	store    *session.Store
}

func NewUserHandler(userRepo repository.UserRepository, sessRepo repository.SessionRepository, store *session.Store) *userHandler {
	return &userHandler{
		userRepo: userRepo,
		sessRepo: sessRepo,
		store:    store,
	}
}

func (u *userHandler) Register(c *fiber.Ctx) error {
	login := c.FormValue("login")
	password := c.FormValue("password")

	log.Info(login)
	log.Info(password)

	hashPassword := string(hashPassword(login, password))

	user := &model.User{
		Login:    login,
		Password: hashPassword,
	}

	createdUser, err := u.userRepo.Create(context.Background(), user)
	if err != nil {
		log.Errorf("%v\n", err)
		if err == repository.ErrUniqueViolation {
			return c.SendStatus(http.StatusConflict)
		}
		return c.SendStatus(http.StatusInternalServerError)
	}

	token := uuid.New().String()
	status := updateToken(c, u.store, token)
	if status != 0 {
		return c.SendStatus(status)
	}

	err = u.sessRepo.CreateToken(context.Background(), &model.Session{UserID: createdUser.ID, Token: token})
	if err != nil {
		log.Errorf("%v\n", err)
		return c.Status(http.StatusInternalServerError).JSON(err)
	}

	return c.Redirect("/list", http.StatusSeeOther)
}

func (u *userHandler) Login(c *fiber.Ctx) error {
	login := c.FormValue("login")
	password := c.FormValue("password")

	user, err := u.userRepo.GetByLogin(context.Background(), login)
	if err != nil {
		log.Errorf("%v\n", err)
		if errors.Is(err, repository.ErrNoRows) {
			return c.SendStatus(http.StatusNotFound)
		}
		return c.SendStatus(http.StatusInternalServerError)
	}

	hashedSha256Password := sha256Hash(login, password)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), hashedSha256Password)
	if err != nil {
		log.Errorf("%v\n", err)
		return c.SendStatus(http.StatusUnauthorized)
	}

	token := uuid.New().String()
	u.sessRepo.UpdateToken(context.Background(), user.ID, token)

	status := updateToken(c, u.store, token)
	if status != 0 {
		return c.SendStatus(status)
	}

	return c.Redirect("/list", http.StatusSeeOther)
}

func (u *userHandler) Logout(c *fiber.Ctx) error {
	sess, err := u.store.Get(c)
	if err != nil {
		log.Errorf("%v\n", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	err = sess.Destroy()
	if err != nil {
		log.Errorf("%v\n", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Redirect("/", fiber.StatusSeeOther)
}

func updateToken(c *fiber.Ctx, store *session.Store, token string) int {
	sess, err := store.Get(c)
	if err != nil {
		return fiber.StatusForbidden
	}

	sess.Set("session_token", token)
	sess.SetExpiry(1 * time.Hour)

	err = sess.Save()
	if err != nil {
		return http.StatusInternalServerError
	}

	return 0
}

func hashPassword(login, password string) []byte {
	hashedSha256Password := sha256Hash(login, password)

	bcryptPassword, err := bcrypt.GenerateFromPassword(hashedSha256Password, bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to generate bcrypt %v\n", err)
		return nil
	}

	return bcryptPassword
}

func sha256Hash(login, password string) []byte {
	passwordBytes := []byte(password + login)

	hash := sha256.New()
	_, err := hash.Write(passwordBytes)
	if err != nil {
		log.Fatalf("failed to write bytes password in hash %v\n", err)
		return nil
	}
	hashedSha256Password := hash.Sum(nil)

	return hashedSha256Password
}
