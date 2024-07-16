package handlers

import (
	"log"
	"net/http"

	"context"

	"github.com/FACorreiaa/glasses-management-platform/app/models"
	"github.com/FACorreiaa/glasses-management-platform/app/view/admin"
	"github.com/FACorreiaa/glasses-management-platform/app/view/settings"
)

func (h *Handler) renderSettingsSidebar() []models.SidebarItem {
	sidebar := []models.SidebarItem{
		{Path: "/settings", Label: "Home"},
		{Path: "/settings/admin", Label: "Change details"},
		{Path: "/settings/collaborators", Label: "View collaborators"},
		{Path: "/settings/glasses", Label: "View glasses"},
		{Path: "/settings/operations", Label: "View transactions"},
		{Path: "/log-out", Label: "Log out"},
	}
	return sidebar
}

func (h *Handler) UpdateAdmin(w http.ResponseWriter, r *http.Request) error {
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

	if user == nil || user.Role != "admin" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return nil
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return err
	}

	g := models.UpdateUserForm{
		UserID:          user.ID,
		Email:           r.FormValue("email"),
		Username:        r.FormValue("username"),
		Password:        r.FormValue("password"),
		PasswordConfirm: r.FormValue("password_confirm"),
	}

	err = h.service.UpdateUser(context.Background(), g)
	if err != nil {
		if err.Error() == "username already exists" || err.Error() == "email already exists" {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
		}
		return err
	}

	w.Header().Set("HX-Redirect", "/settings/admin")
	return nil
}

func (h *Handler) UpdateAdminPage(w http.ResponseWriter, r *http.Request) error {
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

	if user == nil {
		log.Printf("User not authenticated")
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return nil
	}

	g, err := h.service.GetAdminID(context.Background(), user.ID)
	if err != nil {
		log.Printf("Failed to retrieve user information: %v", err)
		http.Error(w, "Failed to retrieve glasses", http.StatusInternalServerError)
		return err
	}

	if g == nil {
		log.Printf("No user found for given ID")
		http.Error(w, "No user found for given ID", http.StatusNotFound)
		return nil
	}

	form := models.UpdateUserForm{
		Values: map[string]string{
			"Username": g.Username,
			"Email":    g.Email,
		},
	}

	f := settings.AdminUpdateForm(form, user.ID)
	updatePage := admin.UserLayoutPage("Update users", "form to update users", f)
	return h.CreateLayout(w, r, "Update Glasses", updatePage).Render(context.Background(), w)
}

func (h *Handler) SettingsPage(w http.ResponseWriter, r *http.Request) error {
	s := settings.AdminSettingsPage(models.SettingsPage{})
	page := settings.AdminSettingsLayoutPage("Settings Page", "Admin settings main page", h.renderSettingsSidebar(), s)
	data := h.CreateLayout(w, r, "Settings", page).Render(context.Background(), w)
	return data
}
