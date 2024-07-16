package handlers

import (
	"net/http"

	"context"

	"github.com/FACorreiaa/glasses-management-platform/app/models"
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

func (h *Handler) SettingsPage(w http.ResponseWriter, r *http.Request) error {
	s := settings.UserSettingsPage(models.SettingsPage{})
	page := settings.UserSettingsLayoutPage("Settings Page", "Admin settings main page", h.renderSettingsSidebar(), s)
	data := h.CreateLayout(w, r, "Settings", page).Render(context.Background(), w)
	return data
}
