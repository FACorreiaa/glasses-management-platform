package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	httperror "github.com/FACorreiaa/glasses-management-platform/app/errors"
	"github.com/FACorreiaa/glasses-management-platform/app/models"
	"github.com/FACorreiaa/glasses-management-platform/app/static/svg"
	"github.com/FACorreiaa/glasses-management-platform/app/view/components"
	"github.com/FACorreiaa/glasses-management-platform/app/view/pages"
	"github.com/FACorreiaa/glasses-management-platform/app/view/shipping"
	"github.com/a-h/templ"
	"github.com/gorilla/mux"
)

func (h *Handler) renderShippingSidebar() []models.SidebarItem {
	sidebar := []models.SidebarItem{
		{Path: "/", Label: "Home"},
		{Path: "/log-out", Label: "Log out"},
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
	id := vars["card_id_number"]
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return err
	}

	form := models.ShippingDetailsForm{
		Name:             r.FormValue("name"),
		CardID:           r.FormValue("card_id_number"),
		Email:            r.FormValue("email"),
		Reference:        r.FormValue("reference"),
		LeftEyeStrength:  parseFloat(r.FormValue("left_eye_strength")),
		RightEyeStrength: parseFloat(r.FormValue("right_eye_strength")),
	}

	f := shipping.ShippingUpdateForm(form, id)
	home := pages.MainLayoutPage("Insert Shipping Form", "Insert Shipping Form", sidebar, f)
	return h.CreateLayout(w, r, "Insert Shipping Form", home).Render(context.Background(), w)
}
