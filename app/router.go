package app

import (
	"embed"
	"log/slog"
	"net/http"

	"github.com/FACorreiaa/glasses-management-platform/app/handlers"
	"github.com/FACorreiaa/glasses-management-platform/app/repository"
	"github.com/FACorreiaa/glasses-management-platform/app/services"
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
	glassesRepo := repository.NewGlassesRepository(pool)
	adminRepo := repository.NewAdminRepository(pool, redisClient, validate, sessions.NewCookieStore(sessionSecret))
	customerRepo := repository.NewCustomerRepository(pool)
	// Middleware
	middleware := &repository.MiddlewareRepository{
		Pgpool:      pool,
		RedisClient: redisClient,
		Validator:   validate,
		Sessions:    sessions.NewCookieStore(sessionSecret),
	}

	// Service
	service := services.NewService(authRepo, glassesRepo, adminRepo, customerRepo)

	// Handler
	handler := handlers.NewHandler(service, sessions.NewCookieStore(sessionSecret), pool, redisClient)

	return handler, middleware
}

func Router(pool *pgxpool.Pool, sessionSecret []byte, redisClient *redis.Client) http.Handler {
	validate := validator.New()
	translator, _ := ut.New(en.New(), en.New()).GetTranslator("en")
	if err := enTranslations.RegisterDefaultTranslations(validate, translator); err != nil {
		slog.Error(" registering translations", "error", err)
	}

	r := mux.NewRouter()

	// Static files
	r.PathPrefix("/static/").Handler(http.FileServer(http.FS(staticFS)))
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, _ *http.Request) {
		file, _ := staticFS.ReadFile("static/favicon.ico")

		w.Header().Set("Cache-Control", "max-age=3600")

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
	// noAuth.HandleFunc("/register", handler(h.RegisterPage)).Methods(http.MethodGet)
	// noAuth.HandleFunc("/register", handler(h.RegisterPost)).Methods(http.MethodPost)

	// Authenticated routes
	auth := r.NewRoute().Subrouter()
	auth.Use(middleware.AuthMiddleware)
	auth.Use(middleware.RequireAuth)

	auth.HandleFunc("/logout", handler(h.Logout)).Methods(http.MethodPost)
	auth.HandleFunc("/settings", handler(h.SettingsPage)).Methods(http.MethodGet)

	// Glasses
	auth.HandleFunc("/glasses", handler(h.GlassesPage)).Methods(http.MethodGet)
	auth.HandleFunc("/glasses/register", handler(h.GlassesRegisterPage)).Methods(http.MethodGet)
	auth.HandleFunc("/glasses/register", handler(h.InsertGlasses)).Methods(http.MethodPost)
	auth.HandleFunc("/glasses/{glasses_id}", handler(h.DeleteGlasses)).Methods(http.MethodDelete)
	auth.HandleFunc("/glasses/type/{type}", handler(h.GlassesTypePage)).Methods(http.MethodGet)
	auth.HandleFunc("/glasses/{glasses_id}/edit", handler(h.UpdateGlassesPage)).Methods(http.MethodGet)
	auth.HandleFunc("/glasses/{glasses_id}/update", handler(h.UpdateGlasses)).Methods(http.MethodPut)
	auth.HandleFunc("/glasses/{stock}/inventory", handler(h.GlassesStockPage)).Methods(http.MethodGet)

	// Collaborators
	auth.HandleFunc("/collaborators/register", handler(h.UserInsertPage)).Methods(http.MethodGet)
	auth.HandleFunc("/collaborators/register", handler(h.UserRegisterPost)).Methods(http.MethodPost)
	auth.HandleFunc("/collaborators/{user_id}", handler(h.DeleteUser)).Methods(http.MethodDelete)
	auth.HandleFunc("/collaborators/{user_id}/edit", handler(h.UpdateUserPage)).Methods(http.MethodGet)
	auth.HandleFunc("/collaborators/{user_id}/update", handler(h.UpdateUser)).Methods(http.MethodPut)

	auth.HandleFunc("/settings/collaborators", handler(h.UsersPage)).Methods(http.MethodGet)
	auth.HandleFunc("/settings/admin", handler(h.UpdateAdminPage)).Methods(http.MethodGet)
	auth.HandleFunc("/settings/admin/update", handler(h.UpdateAdmin)).Methods(http.MethodPut)
	auth.HandleFunc("/settings/glasses", handler(h.SettingsGlassesPage)).Methods(http.MethodGet)

	auth.HandleFunc("/customer/glasses/{glasses_id}/send", handler(h.InsertShippingFormPage)).Methods(http.MethodGet)
	// Settings

	return r
}

func handler(fn func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			slog.Error(" handling request", "error", err)
		}
	}
}
