package models

import (
<<<<<<< HEAD
=======
	"fmt"
	"strings"
>>>>>>> master
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
	Updated         bool
	Values          map[string]string
	FieldErrors     map[string]string
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

// EyePrescription holds the optical details for a single eye
type EyePrescription struct {
	Sph   *float64 // Sphere
	Cyl   *float64 // Cylinder
	Axis  *int     // Axis (degrees)
	Add   *float64 // Addition (for progressives/bifocals)
	Prism *float64 // Prism power
	Base  *string  // Prism base direction (UP, DOWN, IN, OUT)
}

// Helper method to format the prescription for display
func (ep EyePrescription) String() string {
	var parts []string
	if ep.Sph != nil {
		parts = append(parts, fmt.Sprintf("Sph: %+0.2f", *ep.Sph))
	}
	if ep.Cyl != nil && ep.Axis != nil { // Cyl and Axis go together
		parts = append(parts, fmt.Sprintf("Cyl: %+0.2f", *ep.Cyl))
		parts = append(parts, fmt.Sprintf("Ax: %d°", *ep.Axis))
	}
	if ep.Add != nil {
		parts = append(parts, fmt.Sprintf("Add: %+0.2f", *ep.Add))
	}
	if ep.Prism != nil && ep.Base != nil { // Prism and Base go together
		parts = append(parts, fmt.Sprintf("Prism: %0.2f", *ep.Prism))
		parts = append(parts, fmt.Sprintf("Base: %s", *ep.Base))
	}
	if len(parts) == 0 {
		return "N/A" // Or empty string
	}
	return strings.Join(parts, " ")
}

type Glasses struct {
	UserName          string          `json:"username"`
	UserEmail         string          `json:"email"`
	UserID            uuid.UUID       `json:"user_id"`
	GlassesID         uuid.UUID       `json:"glasses_id"`
	Reference         string          `json:"reference"`
	Brand             string          `json:"brand"`
	Color             string          `json:"color"`
	LeftPrescription  EyePrescription `json:"left_prescription"`
	RightPrescription EyePrescription `json:"right_prescription"`
	Type              string          `json:"type"`
	IsInStock         bool            `json:"is_in_stock"`
	Feature           string          `json:"features"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	FieldErrors       map[string]string
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
	FieldErrors map[string]string
}

type UsersTable struct {
	Column      []ColumnItems
	Users       []UserSession
	PrevPage    int
	NextPage    int
	Page        int
	LastPage    int
	OrderParam  string
	SortParam   string
	FieldErrors map[string]string
}

type GlassesForm struct {
	UserID    uuid.UUID `json:"user_id" schema:"user_id"`
	GlassesID uuid.UUID `json:"glasses_id" schema:"glasses_id"`
	Reference string    `json:"reference" schema:"reference"`
	Brand     string    `json:"brand" schema:"brand"`

	LeftSph   float64 `json:"left_sph" schema:"left_sph"`
	LeftCyl   float64 `json:"left_cyl" schema:"left_cyl"`
	LeftAxis  float64 `json:"left_axis" schema:"left_axis"`
	LeftAdd   float64 `json:"left_add" schema:"left_add"`
	LeftPrism float64 `json:"left_prism" schema:"left_prism"`
	LeftBase  float64 `json:"left_base" schema:"left_base"`

	RightSph   float64 `json:"right_sph" schema:"right_sph"`
	RightCyl   float64 `json:"right_cyl" schema:"right_cyl"`
	RightAxis  float64 `json:"right_axis" schema:"right_axis"`
	RightAdd   float64 `json:"right_add" schema:"right_add"`
	RightPrism float64 `json:"right_prism" schema:"right_prism"`
	RightBase  float64 `json:"right_base" schema:"right_base"`

	Color       string `json:"color" schema:"color"`
	Type        string `json:"type" schema:"type"`
	IsInStock   bool   `json:"is_in_stock" schema:"is_in_stock"`
	Feature     string `json:"features" schema:"features"`
	Updated     bool
	Values      map[string]string
	FieldErrors map[string]string
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
	FieldErrors     map[string]string
}

type ShippingDetails struct {
	CustomerID uuid.UUID `json:"customer_id"`
	Name       string    `json:"name"`
	CardID     string    `json:"card_id_number"`
	Email      string    `json:"email"`
	Reference  string    `json:"reference"`
	LeftEye    float64   `json:"left_eye_strength"`
	RightEye   float64   `json:"right_eye_strength"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
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
	CustomerID        uuid.UUID `json:"customer_id"`
	CollaboratorName  string    `json:"username"`
	CollaboratorEmail string    `json:"collaborator_email"`
	Name              string    `json:"name"`
	CardID            string    `json:"card_id_number"`
	Email             string    `json:"email"`
	Reference         string    `json:"reference"`
	LeftEye           float64   `json:"left_eye_strength"`
	RightEye          float64   `json:"right_eye_strength"`
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

type ShippingDetailsForm struct {
	CustomerID uuid.UUID `json:"customer_id"`
	Name       string    `json:"name" schema:"name"`
	CardID     string    `json:"card_id_number" schema:"card_id_number"`
	Email      string    `json:"email" schema:"email"`
	Reference  string    `json:"reference" schema:"reference"`
	LeftEye    float64   `json:"left_eye_strength" schema:"left_eye_strength"`
	RightEye   float64   `json:"right_eye_strength" schema:"right_eye_strength"`
	// CreatedAt        time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Updated     bool
	Values      map[string]string
	FieldErrors map[string]string
}
