package models

import (
	"time"

	"github.com/a-h/templ"
	"github.com/google/uuid"
)

type ctxKey int

const (
	CtxKeyAuthUser ctxKey = iota
)

type LoginForm struct {
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required"`
}

type UserSession struct {
	ID           uuid.UUID
	Username     string
	Email        string
	PasswordHash []byte
	Role         string
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
}

type UserSessionToken struct {
	Token     string
	CreatedAt *time.Time
	User      *UserSession
}

type RegisterForm struct {
	Username        string `form:"username" validate:"required"`
	Email           string `form:"email" validate:"required,email"`
	Password        string `form:"password" validate:"required,min=8,max=72"`
	PasswordConfirm string `form:"password_confirm" validate:"required,eqfield=Password"`
}

type NavItem struct {
	Path     string
	Icon     templ.Component
	Label    string
	IsLogout bool
}

type TabItem struct {
	Path  string
	Icon  string
	Label string
}

type SidebarItem struct {
	Path       string
	Label      string
	ActivePath string
	SubItems   []SidebarItem
}

type LayoutTempl struct {
	Title     string
	Nav       []NavItem
	ActiveNav string
	User      *UserSession
	Content   templ.Component
}

type SettingsPage struct {
	Updated bool
	Errors  []string
	User    UserSession
}

type LoginPage struct {
	Errors []string
}

type RegisterPage struct {
	Errors []string
	Values map[string]string
}

type Columns struct {
	Title string
}

type Glasses struct {
	GlassesID uuid.UUID `json:"glasses_id"`
	Reference string    `json:"reference"`
	Brand     string    `json:"brand"`
	Color     string    `json:"color"`
	LeftEye   float64   `json:"left_eye_strength"`
	RightEye  float64   `json:"right_eye_strength"`
	Type      string    `json:"type"`
	IsInStock bool      `json:"is_in_stock"`
	Feature   string    `json:"features"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ColumnItems struct {
	Title     string
	Icon      templ.Component
	SortParam string
}

type GlassesTable struct {
	Column      []ColumnItems
	Glasses     []Glasses
	PrevPage    int
	NextPage    int
	Page        int
	LastPage    int
	FilterBrand string
	OrderParam  string
	SortParam   string
}

type UsersTable struct {
	Column     []ColumnItems
	Users      []UserSession
	PrevPage   int
	NextPage   int
	Page       int
	LastPage   int
	OrderParam string
	SortParam  string
}

type GlassesForm struct {
	Reference string  `json:"reference" schema:"reference"`
	Brand     string  `json:"brand" schema:"brand"`
	Color     string  `json:"color" schema:"color"`
	LeftEye   float64 `json:"left_eye_strength" schema:"left_eye_strength"`
	RightEye  float64 `json:"right_eye_strength" schema:"right_eye_strength"`
	Type      string  `json:"type" schema:"type"`
	IsInStock bool    `json:"is_in_stock" schema:"is_in_stock"`
	Feature   string  `json:"features" schema:"features"`
	Updated   bool
	Values    map[string]string
	Errors    map[string]string
}

type UpdateUserForm struct {
	UserID          uuid.UUID `json:"user_id"`
	Username        string    `form:"username"`
	Email           string    `form:"email"`
	Role            string    `form:"role"`
	Password        string    `form:"password"`
	PasswordConfirm string    `form:"password_confirm"`
	Updated         bool
	Values          map[string]string
	Errors          map[string]string
}
