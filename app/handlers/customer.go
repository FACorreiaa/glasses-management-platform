package handlers

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	httperror "github.com/FACorreiaa/glasses-management-platform/app/errors"
	"github.com/FACorreiaa/glasses-management-platform/app/models"
	"github.com/FACorreiaa/glasses-management-platform/app/static/svg"
	"github.com/FACorreiaa/glasses-management-platform/app/view/components"
	"github.com/FACorreiaa/glasses-management-platform/app/view/customer"
	"github.com/FACorreiaa/glasses-management-platform/app/view/pages"
	"github.com/FACorreiaa/glasses-management-platform/app/view/shipping"
	"github.com/a-h/templ"
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

	if err = h.service.InsertShippingDetails(r.Context(), glassesID, user.ID, customerForm, shipping); err != nil {
		slog.Error("inserting shipping details", "err", err)
		http.Error(w, "Error processing shipping", http.StatusInternalServerError)
		return nil
	}

	http.Redirect(w, r, "/glasses", http.StatusSeeOther)

	return nil
}

// shipping details

func (h *Handler) renderShippingSidebar() []models.SidebarItem {
	sidebar := []models.SidebarItem{
		{Path: "/", Label: "Home"},
		{Path: "/logout", Label: "Log out"},
	}
	return sidebar
}

func (h *Handler) getShipping(w http.ResponseWriter, r *http.Request) (int, []models.ShippingDetails, error) {
	pageSize := 10
	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	reference := r.FormValue("reference")

	leftEyeStr := r.FormValue("left_eye_strength")
	rightEyeStr := r.FormValue("right_eye_strength")

	var leftEye, rightEye *float64

	if leftEyeStr != "" {
		parsedLeftEye, err := strconv.ParseFloat(leftEyeStr, 64)
		if err != nil {
			HandleError(err, "parse left eye")
			return 0, nil, err
		}
		leftEye = &parsedLeftEye
	}

	if rightEyeStr != "" {
		parsedRightEye, err := strconv.ParseFloat(rightEyeStr, 64)
		if err != nil {
			HandleError(err, "parse right eye")
			return 0, nil, err
		}
		rightEye = &parsedRightEye
	}

	// Debug print statements to check values
	if leftEye != nil {
		fmt.Printf("leftEye: %f\n", *leftEye)
	} else {
		fmt.Println("leftEye is nil")
	}

	if rightEye != nil {
		fmt.Printf("rightEye: %f\n", *rightEye)
	} else {
		fmt.Println("rightEye is nil")
	}

	s, err := h.service.GetShippingDetails(context.Background(), page, pageSize, orderBy, sortBy, reference, leftEye, rightEye)
	if err != nil {
		httperror.ErrNotFound.WriteError(w)
		return 0, nil, err
	}

	return page, s, nil
}

func (h *Handler) renderShippingDetailsTable(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	var page int
	var sortAux string
	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")
	brand := r.FormValue("brand")

	if sortBy == ASC {
		sortAux = DESC
	} else {
		sortAux = ASC
	}

	columnNames := []models.ColumnItems{
		{Title: "Name", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Card ID", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Email", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Glasses Reference", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Right Eye Strength", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Left Eye Strength", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Created at:", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Updated at:", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
	}

	page, s, _ := h.getShipping(w, r)

	if len(s) == 0 {
		message := components.EmptyPageComponent()
		return message, nil
	}

	nextPage := page + 1
	prevPage := page - 1
	if prevPage <= 1 {
		prevPage = 1
	}

	lastPage, err := h.service.GetSum()
	if err != nil {
		HandleError(err, " fetching tax")
		return nil, err
	}
	data := models.ShippingDetailsTable{
		Column:      columnNames,
		Shipping:    s,
		PrevPage:    prevPage,
		NextPage:    nextPage,
		Page:        page,
		LastPage:    lastPage,
		FilterBrand: brand,
		OrderParam:  orderBy,
		SortParam:   sortAux,
	}
	t := shipping.ShippingDetailsSimple(data, models.ShippingDetails{})

	return t, nil
}

func (h *Handler) GetShippingDetailsPage(w http.ResponseWriter, r *http.Request) error {
	sidebar := h.renderShippingSidebar()
	renderTable, err := h.renderShippingDetailsTable(w, r)
	if err != nil {
		HandleError(err, "rendering glasses table")
	}
	home := pages.MainLayoutPage("Insert Shipping Form", "Insert Shipping Form", sidebar, renderTable)
	return h.CreateLayout(w, r, "Insert Shipping Form", home).Render(context.Background(), w)
}

func (h *Handler) UpdateCustomerPage(w http.ResponseWriter, r *http.Request) error {
	sidebar := h.renderSettingsSidebar()

	vars := mux.Vars(r)
	id := vars["customer_id"]
	println("Received customer_id from URL UpdateCustomerPage:", id)
	customerID, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid glasses ID", http.StatusBadRequest)
		return err
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return err
	}

	g, err := h.service.GetCustomerGlassesID(context.Background(), customerID)
	if err != nil {
		http.Error(w, "Failed to retrieve glasses", http.StatusInternalServerError)
		return err
	}

	le := strconv.FormatFloat(g.LeftEye, 'f', 2, 64)
	re := strconv.FormatFloat(g.RightEye, 'f', 2, 64)

	form := models.ShippingDetailsForm{
		CustomerID: customerID,
		Values: map[string]string{
			"Name":      r.FormValue("name"),
			"CardID":    r.FormValue("card_id_number"),
			"Email":     r.FormValue("email"),
			"Reference": r.FormValue("reference"),
			"LeftEye":   le,
			"RightEye":  re,
		},
	}

	f := shipping.ShippingUpdateForm(form)
	home := pages.MainLayoutPage("Insert Shipping Form", "Insert Shipping Form", sidebar, f)
	return h.CreateLayout(w, r, "Insert Shipping Form", home).Render(context.Background(), w)
}

func (h *Handler) UpdateCustomer(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["customer_id"]
	println("Received customer_id from URL UpdateCustomer:", id)
	customerID, err := uuid.Parse(id)

	println("Received customer_id from URL:", id)
	if err != nil {
		http.Error(w, "Invalid glasses ID", http.StatusBadRequest)
		return err
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return err
	}

	form := models.ShippingDetailsForm{
		CustomerID:  customerID,
		Name:        r.FormValue("name"),
		CardID:      r.FormValue("card_id_number"),
		Email:       r.FormValue("email"),
		Reference:   r.FormValue("reference"),
		LeftEye:     parseFloat(r.FormValue("left_eye_strength")),
		RightEye:    parseFloat(r.FormValue("right_eye_strength")),
		FieldErrors: make(map[string]string),
	}

	cardIDNumber, err := h.service.GetCardIDFromShipping(r.Context(), customerID)
	if err != nil {
		slog.Error("Error fetching card_id_number", "err", err)
		http.Error(w, "Error fetching card_id_number", http.StatusInternalServerError)
		return nil
	}

	if cardIDNumber == form.CardID {
		form.FieldErrors["card_id_number"] = "Card ID number already exists"
	}

	referenceNumber, err := h.service.GetReferenceNumberFromShipping(r.Context(), customerID)
	if err != nil {
		slog.Error("Error fetching reference_number", "err", err)
		http.Error(w, "Error fetching reference_number", http.StatusInternalServerError)
		return nil
	}

	if referenceNumber == form.Reference {
		form.FieldErrors["reference"] = "Reference number already exists"
	}

	if len(form.FieldErrors) > 0 {
		sidebar := h.renderSidebar()
		f := shipping.ShippingUpdateForm(form)
		register := pages.MainLayoutPage("Insert Shipping Form", "Insert Shipping Form", sidebar, f)
		return h.CreateLayout(w, r, "Insert Shipping Form", register).Render(context.Background(), w)
	}

	if err := h.service.UpdateShippingDetails(context.Background(), form, customerID); err != nil {
		http.Error(w, "Failed to update shipping details", http.StatusInternalServerError)
		return err
	}

	w.Header().Set("HX-Redirect", "/settings/shipping")

	return nil
}
