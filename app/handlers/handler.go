package handlers

import (
	"log"
	"net/http"

	"context"

	"github.com/FACorreiaa/glasses-management-platform/app/models"
	"github.com/FACorreiaa/glasses-management-platform/app/services"
	"github.com/FACorreiaa/glasses-management-platform/app/view/pages"
	"github.com/a-h/templ"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

const ASC = "ASC"
const DESC = "DESC"

type contextKey string

var themeContextKey contextKey = "theme"

type Handler struct {
	service     *services.Service
	formDecoder *form.Decoder
	validator   *validator.Validate
	translator  ut.Translator
	sessions    *sessions.CookieStore
	pool        *pgxpool.Pool
	redisClient *redis.Client
}

func NewHandler(s *services.Service, sessions *sessions.CookieStore,
	pool *pgxpool.Pool, redisClient *redis.Client) *Handler {
	decoder := form.NewDecoder()
	validate := validator.New()
	translator, _ := ut.New(en.New(), en.New()).GetTranslator("en")
	return &Handler{
		service:     s,
		formDecoder: decoder,
		validator:   validate,
		translator:  translator,
		sessions:    sessions,
		pool:        pool,
		redisClient: redisClient,
	}
}

func HandleError(err error, message string) {
	if err != nil {
		log.Printf("%s: %v", message, err)
	}
}

func (h *Handler) CreateLayout(_ http.ResponseWriter, r *http.Request, title string, data templ.Component) templ.Component {
	var user *models.UserSession
	userCtx := r.Context().Value(models.CtxKeyAuthUser)
	if userCtx != nil {
		switch u := userCtx.(type) {
		case *models.UserSession:
			user = u
		default:
			log.Printf("Unexpected type in userCtx: %T", userCtx)
		}
	}

	var nav []models.NavItem

	switch {
	case user == nil:
		nav = []models.NavItem{
			{Path: "/", Label: "Home"},
			{Path: "/login", Label: "Sign in"},
			// {Path: "/register", Label: "Sign up"},
		}
	case user.Role == "admin":
		nav = []models.NavItem{
			{Path: "/", Label: "Home"},
			{Path: "/glasses", Label: "Glasses Inventory"},
			{Path: "/shipping", Label: "Shipping Orders"},
			{Path: "/settings", Label: "Settings"},
			{Path: "/logout", Label: "Sign out", IsLogout: true},
		}
	default:
		nav = []models.NavItem{
			{Path: "/", Label: "Home"},
			{Path: "/glasses", Label: "Glasses Inventory"},
			{Path: "/shipping", Label: "Shipping Orders"},
			{Path: "/logout", Label: "Sign out", IsLogout: true},
		}
	}

	l := models.LayoutTempl{
		Title:     title,
		Nav:       nav,
		User:      user,
		ActiveNav: r.URL.Path,
		Content:   data,
	}

	return pages.LayoutPage(l)
}

func (h *Handler) Homepage(w http.ResponseWriter, r *http.Request) error {
	home := pages.HomePage()
	return h.CreateLayout(w, r, "Home Page", home).Render(context.Background(), w)
}
