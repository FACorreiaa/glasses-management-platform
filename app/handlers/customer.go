package handlers

import (
	"context"
	"net/http"

	"github.com/FACorreiaa/glasses-management-platform/app/models"
	"github.com/FACorreiaa/glasses-management-platform/app/view/customer"
	"github.com/FACorreiaa/glasses-management-platform/app/view/glasses"
)

func (h *Handler) InsertShippingFormPage(w http.ResponseWriter, r *http.Request) error {
	sidebar := h.renderSidebar()

	form := customer.CustomerShipingDetailsForm(models.CustomerShippingForm{})
	home := glasses.GlassesLayoutPage("Insert Shipping Form", "Insert Shipping Form", sidebar, form)
	return h.CreateLayout(w, r, "Insert Shipping Form", home).Render(context.Background(), w)
}
