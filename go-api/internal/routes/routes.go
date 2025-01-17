package routes

import (
	"github.com/Energie-Burgenland/ausaestung-info/internal/handlers"
	"github.com/Energie-Burgenland/ausaestung-info/internal/middleware"
	"github.com/Energie-Burgenland/ausaestung-info/internal/repositories"
	"github.com/Energie-Burgenland/ausaestung-info/internal/validation"
	"github.com/Energie-Burgenland/ausaestung-info/utils/database"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, verifier *oidc.IDTokenVerifier, dbContext *database.DbContext, validator *validation.Validator) {
	router := app.Group("/", middleware.JWTMiddleware(verifier))

	registerUserRoutes(router, dbContext, validator)
}

func registerUserRoutes(router fiber.Router, dbContext *database.DbContext, validator *validation.Validator) {
	repo := repositories.NewUserRepository(dbContext)
	handler := handlers.NewUserHandler(&repo, validator)

	users := router.Group("/user")
	users.Get("/", handler.GetUsers)
	users.Get("/:id", handler.GetUser)
	users.Post("/", handler.CreateUser)
	users.Put("/:id", handler.UpdateUser)
	users.Delete("/:id", handler.DeleteUser)
	users.Post("/import", handler.ImportUsers)
}
