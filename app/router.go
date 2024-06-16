package app

import (
	"embed"
	"log/slog"
	"net/http"

	"github.com/FACorreiaa/Aviation-tracker/app/handlers"
	"github.com/FACorreiaa/Aviation-tracker/app/repository"
	"github.com/FACorreiaa/Aviation-tracker/app/services"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

//go:embed static
var staticFS embed.FS

func setupBusinessComponents(pool *pgxpool.Pool, redisClient *redis.Client, validate *validator.Validate,
	sessionSecret []byte) (*handlers.Handler, *repository.MiddlewareRepository) {
	// Business components

	authRepo := repository.NewAccountRepository(pool, redisClient, validate, sessions.NewCookieStore(sessionSecret))

	// Middleware
	middleware := &repository.MiddlewareRepository{
		Pgpool:      pool,
		RedisClient: redisClient,
		Validator:   validate,
		Sessions:    sessions.NewCookieStore(sessionSecret),
	}

	// Service
	service := services.NewService(authRepo)

	// Handler
	handler := handlers.NewHandler(service, sessions.NewCookieStore(sessionSecret), pool, redisClient)

	return handler, middleware
}

func Router(pool *pgxpool.Pool, sessionSecret []byte, redisClient *redis.Client) http.Handler {
	validate := validator.New()
	translator, _ := ut.New(en.New(), en.New()).GetTranslator("en")
	if err := enTranslations.RegisterDefaultTranslations(validate, translator); err != nil {
		slog.Error("Error registering translations", "error", err)
	}

	r := mux.NewRouter()

	// Static files
	r.PathPrefix("/static/").Handler(http.FileServer(http.FS(staticFS)))
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, _ *http.Request) {
		file, _ := staticFS.ReadFile("static/favicon.ico")
		// etag := setEtag("v1", file)

		w.Header().Set("Cache-Control", "max-age=3600")
		// w.Header().Set("Content-Type", http.DetectContentType(file))
		// w.Header().Set("Etag", etag)

		_, err := w.Write(file)
		if err != nil {
			return
		}
	})

	h, middleware := setupBusinessComponents(pool, redisClient, validate, sessionSecret)

	// Public routes, authentication is optional
	optAuth := r.NewRoute().Subrouter()
	optAuth.Use(middleware.AuthMiddleware)
	optAuth.HandleFunc("/", handler(h.Homepage)).Methods(http.MethodGet)

	// Routes that shouldn't be available to authenticated users
	noAuth := r.NewRoute().Subrouter()
	noAuth.Use(middleware.AuthMiddleware)
	noAuth.Use(middleware.RedirectIfAuth)

	noAuth.HandleFunc("/login", handler(h.LoginPage)).Methods(http.MethodGet)
	noAuth.HandleFunc("/login", handler(h.LoginPost)).Methods(http.MethodPost)
	noAuth.HandleFunc("/register", handler(h.RegisterPage)).Methods(http.MethodGet)
	noAuth.HandleFunc("/register", handler(h.RegisterPost)).Methods(http.MethodPost)

	// Authenticated routes
	auth := r.NewRoute().Subrouter()
	auth.Use(middleware.AuthMiddleware)
	auth.Use(middleware.RequireAuth)

	auth.HandleFunc("/logout", handler(h.Logout)).Methods(http.MethodPost)
	auth.HandleFunc("/settings", handler(h.SettingsPage)).Methods(http.MethodGet)
	return r
}

func handler(fn func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			slog.Error("Error handling request", "error", err)
		}
	}
}
