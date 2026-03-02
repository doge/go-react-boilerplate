package routes

import (
	"context"
	"net/http"
	"server/api/controllers"
	"server/api/repository"
	"server/api/service"
	"server/internal/middleware"
	"server/internal/models"
	"server/internal/security"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type routes struct {
	router *chi.Mux
}

func SetupRoutes(db *mongo.Database, crypto *security.AESGCM, config *models.Config) (*chi.Mux, error) {
	api := routes{
		router: chi.NewMux(),
	}

	api.router.Use(middleware.LoggingMiddleware)
	api.router.Use(middleware.RateLimitMiddleware)

	if err := api.AuthRoutes(api.router, db, crypto, config); err != nil {
		return nil, err
	}

	return api.router, nil
}

func (r routes) AuthRoutes(router *chi.Mux, db *mongo.Database, crypto *security.AESGCM, config *models.Config) error {
	accessTTL := time.Duration(config.Auth.AccessTTLMinutes) * time.Minute
	if accessTTL <= 0 {
		accessTTL = 15 * time.Minute
	}

	refreshTTL := time.Duration(config.Auth.RefreshTTLDays) * 24 * time.Hour
	if refreshTTL <= 0 {
		refreshTTL = 30 * 24 * time.Hour
	}

	tokenManager, err := security.NewTokenManager(config.Auth.JWTSecret, config.Auth.Issuer, accessTTL)
	if err != nil {
		return err
	}

	userRepository := repository.NewUserRepository(db.Collection("users"))
	if err = userRepository.EnsureIndexes(context.Background()); err != nil {
		return err
	}

	refreshRepository := repository.NewRefreshSessionRepository(db.Collection("refresh_sessions"))
	if err = refreshRepository.EnsureIndexes(context.Background()); err != nil {
		return err
	}

	userService := service.NewUserService(
		userRepository,
		refreshRepository,
		crypto,
		tokenManager,
		service.AuthSettings{
			RefreshTTL:        refreshTTL,
			RefreshCookieName: config.Auth.RefreshCookieName,
		},
	)
	userController := controllers.NewUserController(
		userService,
		config.Auth.RefreshCookieName,
		config.Auth.SecureCookies,
		refreshTTL,
		parseSameSite(config.Auth.RefreshSameSite),
	)

	router.Route("/auth", func(authRouter chi.Router) {
		authRouter.Post("/register", userController.Create)
		authRouter.Post("/login", userController.Login)
		authRouter.With(middleware.TrustedOriginMiddleware(config)).Post("/refresh", userController.Refresh)
		authRouter.With(middleware.TrustedOriginMiddleware(config)).Post("/logout", userController.Logout)
		authRouter.With(middleware.AuthMiddleware(tokenManager)).Get("/session", userController.Session)
	})

	return nil
}

func parseSameSite(value string) http.SameSite {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "none":
		return http.SameSiteNoneMode
	case "lax":
		return http.SameSiteLaxMode
	default:
		return http.SameSiteStrictMode
	}
}
