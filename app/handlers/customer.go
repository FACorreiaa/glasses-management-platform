package handlers

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/FACorreiaa/glasses-management-platform/app/models"
	"github.com/FACorreiaa/glasses-management-platform/app/view/customer"
	"github.com/FACorreiaa/glasses-management-platform/app/view/pages"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (h *Handler) InsertShippingFormPage(w http.ResponseWriter, r *http.Request) error {
	sidebar := h.renderSidebar()

	form := customer.CustomerShipingDetailsForm(models.CustomerShippingForm{})
	home := pages.MainLayoutPage("Insert Shipping Form", "Insert Shipping Form", sidebar, form)
	return h.CreateLayout(w, r, "Insert Shipping Form", home).Render(context.Background(), w)
}

func (h *Handler) InsertShippingForm(w http.ResponseWriter, r *http.Request) error {
	var user *models.UserSession
	fieldErrors := make(map[string]string)

	userCtx := r.Context().Value(models.CtxKeyAuthUser)
	if userCtx != nil {
		switch u := userCtx.(type) {
		case *models.UserSession:
			user = u
		default:
			log.Printf("Unexpected type in userCtx: %T", userCtx)
		}
	}

	vars := mux.Vars(r)
	glassesIDStr := vars["glasses_id"]
	glassesID, err := uuid.Parse(glassesIDStr)
	if err != nil {
		http.Error(w, "Invalid glasses ID", http.StatusBadRequest)
		return err
	}

	if err := r.ParseForm(); err != nil {
		HandleError(err, "parsing form")
		return err
	}

	customerForm := models.CustomerShippingForm{
		UserID:         user.ID,
		GlassesID:      glassesID,
		Name:           r.FormValue("name"),
		CardID:         r.FormValue("card_id_number"),
		Address:        r.FormValue("address"),
		AddressDetails: r.FormValue("address_details"),
		City:           r.FormValue("city"),
		Country:        r.FormValue("country"),
		Continent:      r.FormValue("continent"),
		PostalCode:     r.FormValue("postal_code"),
		PhoneNumber:    r.FormValue("phone_number"),
		Email:          r.FormValue("email"),
	}

	shipping := models.Shipping{
		GlassesID:    glassesID,
		ShippedBy:    user.ID,
		ShippingDate: time.Now(),
	}

	cardIDNumber, err := h.service.GetCardIDNumber(r.Context(), customerForm.CustomerID)
	if err != nil {
		slog.Error("Error fetching card_id_number", "err", err)
		http.Error(w, "Error fetching card_id_number", http.StatusInternalServerError)
		return nil
	}

	if cardIDNumber == customerForm.CardID {
		fieldErrors["card_id_number"] = "Card ID number already exists"
	}

	if len(fieldErrors) > 0 {
		form := models.CustomerShippingForm{
			FieldErrors: fieldErrors,
		}
		sidebar := h.renderSidebar()
		f := customer.CustomerShipingDetailsForm(form)
		register := pages.MainLayoutPage("Insert Shipping Form", "Insert Shipping Form", sidebar, f)
		return h.CreateLayout(w, r, "Insert Shipping Form", register).Render(context.Background(), w)
	}

	err = h.service.InsertShippingDetails(r.Context(), glassesID, user.ID, customerForm, shipping)
	if err != nil {
		slog.Error("inserting shipping details", "err", err)
		http.Error(w, "Error processing shipping", http.StatusInternalServerError)
		return nil
	}

	http.Redirect(w, r, "/glasses", http.StatusSeeOther)

	return nil
}
