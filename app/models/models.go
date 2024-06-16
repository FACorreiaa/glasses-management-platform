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
	Bio          string
	Image        *string
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

type FlightStatus string

// const (
//	Scheduled FlightStatus = "scheduled"
//	Active    FlightStatus = "active"
//	Landed    FlightStatus = "landed"
//	Canceled FlightStatus = "canceled"
//	Incident  FlightStatus = "incident"
//	Diverted  FlightStatus = "diverted"
// )

type NavItem struct {
	Path  string
	Icon  templ.Component
	Label string
}

type TabItem struct {
	Path  string
	Icon  string
	Label string
}

type SidebarItem struct {
	Path       string
	Icon       templ.Component
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
