package handlers

import (
	"context"
	"net/http"

	"github.com/FACorreiaa/glasses-management-platform/app/models"
	"github.com/FACorreiaa/glasses-management-platform/app/view/admin"
)

func (h *Handler) renderAdminSidebar() []models.SidebarItem {
	sidebar := []models.SidebarItem{
		{Path: "/", Label: "Home"},
		{Path: "/collaborators", Label: "List Collaborators"},
		{Path: "/collaborators/register", Label: "Insert collaborators"},
		{Path: "/log-out", Label: "Log out"},
	}
	return sidebar
}

// UsersPage users page for admin to manage views
func (h *Handler) UsersPage(w http.ResponseWriter, r *http.Request) error {
	sidebar := h.renderAdminSidebar()
	test := admin.TestComponent()
	users := admin.UserLayoutPage("List of collaborators", "List of collaborators", sidebar, test)
	data := h.CreateLayout(w, r, "Users", users).Render(context.Background(), w)
	return data
}
