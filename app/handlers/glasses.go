package handlers

import (
	"context"
	"net/http"

	"github.com/FACorreiaa/glasses-management-platform/app/view/glasses"
)

func (h *Handler) GlassesPage(w http.ResponseWriter, r *http.Request) error {
	home := glasses.GlassesMainPage()
	return h.CreateLayout(w, r, "Glasses Management Page", home).Render(context.Background(), w)
}
