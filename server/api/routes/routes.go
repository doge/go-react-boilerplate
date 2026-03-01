package routes

import (
	"server/api/controllers"
	"server/api/repository"
	"server/api/service"
	"server/internal/middleware"
	"server/internal/security"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type routes struct {
	router *chi.Mux
}

func SetupRoutes(db *mongo.Database, crypto *security.AESGCM) *chi.Mux {
	api := routes{
		router: chi.NewMux(),
	}

	api.router.Use(middleware.LoggingMiddleware)
	api.router.Use(middleware.RateLimitMiddleware)

	api.AuthRoutes(api.router, db, crypto)

	return api.router
}

func (r routes) AuthRoutes(router *chi.Mux, db *mongo.Database, crypto *security.AESGCM) {
	userRepository := repository.NewUserRepository(db.Collection("users"))
	userService := service.NewUserService(userRepository, crypto)
	userController := controllers.NewUserController(userService)

	router.Route("/auth", func(authRouter chi.Router) {
		authRouter.Post("/register", userController.Create)
		authRouter.Post("/login", userController.Login)
	})
}
