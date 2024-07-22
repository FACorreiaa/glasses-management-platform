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
	Errors      []string
	FieldErrors map[string]string
}

type RegisterFormValues struct {
	Errors          []string
	FieldErrors     map[string]string
	Values          map[string]string
	Username        string `form:"username" validate:"required"`
	Email           string `form:"email" validate:"required,email"`
	Password        string `form:"password" validate:"required,min=8,max=72"`
	PasswordConfirm string `form:"password_confirm" validate:"required,eqfield=Password"`
}

type Columns struct {
	Title string
}

type Glasses struct {
	UserName  string    `json:"username"`
	UserEmail string    `json:"email"`
	UserID    uuid.UUID `json:"user_id"`
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

type CustomerShippingForm struct {
	CustomerID     uuid.UUID `json:"customer_id" schema:"customer_id"`
	UserID         uuid.UUID `json:"user_id"     schema:"user_id"`
	GlassesID      uuid.UUID `json:"glasses_id"  schema:"glasses_id"`
	Name           string    `json:"name"        schema:"name"`
	CardID         string    `json:"card_id_number" schema:"card_id_number"`
	Address        string    `json:"address" schema:"address"`
	AddressDetails string    `json:"address_details" schema:"address_details"`
	City           string    `json:"city" schema:"city"`
	Country        string    `json:"country" schema:"country"`
	Continent      string    `json:"continent" schema:"continent"`
	PostalCode     string    `json:"postal_code" schema:"postal_code"`
	PhoneNumber    string    `json:"phone_number" schema:"phone_number"`
	Email          string    `json:"email" schema:"email"`
	Updated        bool
	Values         map[string]string
	FieldErrors    map[string]string
}

type Customer struct {
	CustomerID     uuid.UUID `json:"customer_id"`
	Name           string    `json:"name"       `
	CardID         string    `json:"card_id_number" `
	Address        string    `json:"address" `
	AddressDetails string    `json:"address_details" `
	City           string    `json:"city" `
	Country        string    `json:"country" `
	Continent      string    `json:"continent" `
	PostalCode     string    `json:"postal_code" `
	PhoneNumber    string    `json:"phone_number" `
	Email          string    `json:"email"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Shipping struct {
	ShippingID   uuid.UUID `json:"shipping_id"`
	GlassesID    uuid.UUID `json:"glasses_id"`
	CustomerID   uuid.UUID `json:"customer_id"`
	ShippedBy    uuid.UUID `json:"shipped_by"`
	ShippingDate time.Time `json:"shipping_date"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
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

type ShippingDetails struct {
	Name             string    `json:"name"`
	CardID           string    `json:"card_id_number"`
	Email            string    `json:"email"`
	Reference        string    `json:"reference"`
	LeftEyeStrength  float64   `json:"left_eye_strength"`
	RightEyeStrength float64   `json:"right_eye_strength"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type ShippingDetailsTable struct {
	Column      []ColumnItems
	Shipping    []ShippingDetails
	PrevPage    int
	NextPage    int
	Page        int
	LastPage    int
	FilterBrand string
	OrderParam  string
	SortParam   string
}

type SettingsShippingDetails struct {
	CollaboratorName  string    `json:"username"`
	CollaboratorEmail string    `json:"collaborator_email"`
	Name              string    `json:"name"`
	CardID            string    `json:"card_id_number"`
	Email             string    `json:"email"`
	Reference         string    `json:"reference"`
	LeftEyeStrength   float64   `json:"left_eye_strength"`
	RightEyeStrength  float64   `json:"right_eye_strength"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type SettingsShippingDetailsTable struct {
	Column      []ColumnItems
	Shipping    []SettingsShippingDetails
	PrevPage    int
	NextPage    int
	Page        int
	LastPage    int
	FilterBrand string
	OrderParam  string
	SortParam   string
}
