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
	"github.com/FACorreiaa/glasses-management-platform/app/view/glasses"
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (h *Handler) renderSidebar() []models.SidebarItem {
	sidebar := []models.SidebarItem{
		{Path: "/", Label: "Home"},
		{Path: "/glasses", Label: "Glasses stock"},
		{Path: "/glasses/register", Label: "Insert glasses"},
		{
			Label: "Type",
			SubItems: []models.SidebarItem{
				{Path: "/glasses/type/adult", Label: "Adult"},
				{Path: "/glasses/type/children", Label: "Children"},
			},
		},
		{
			Label: "Inventory",
			SubItems: []models.SidebarItem{
				{Path: "/glasses/current/inventory", Label: "Check Current Inventory"},
				{Path: "/glasses/shipped/inventory", Label: "Check Shipped Inventory"},
			},
		},
		{Path: "/log-out", Label: "Log out"},
	}
	return sidebar
}

func (h *Handler) getGlasses(w http.ResponseWriter, r *http.Request) (int, []models.Glasses, error) {
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

	g, err := h.service.GetGlasses(context.Background(), page, pageSize, orderBy, sortBy, reference, leftEye, rightEye)
	if err != nil {
		httperror.ErrNotFound.WriteError(w)
		return 0, nil, err
	}

	return page, g, nil
}

func (h *Handler) renderGlassesTable(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
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
		{Title: "Brand", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Color", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Reference", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Left Eye", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Right Eye", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Type", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Has Stock", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Features", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Created At", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Updated At", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
	}

	page, g, _ := h.getGlasses(w, r)

	if len(g) == 0 {
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
	data := models.GlassesTable{
		Column:      columnNames,
		Glasses:     g,
		PrevPage:    prevPage,
		NextPage:    nextPage,
		Page:        page,
		LastPage:    lastPage,
		FilterBrand: brand,
		OrderParam:  orderBy,
		SortParam:   sortAux,
	}
	taxTable := glasses.GlassesTable(data)

	return taxTable, nil
}

func (h *Handler) GlassesPage(w http.ResponseWriter, r *http.Request) error {
	sidebar := h.renderSidebar()
	renderTable, err := h.renderGlassesTable(w, r)
	if err != nil {
		HandleError(err, "rendering glasses table")
	}
	home := glasses.GlassesLayoutPage("Glasses Management Page", "Glasses Management Page", sidebar, renderTable)
	return h.CreateLayout(w, r, "Glasses Management Page", home).Render(context.Background(), w)
}

func (h *Handler) GlassesRegisterPage(w http.ResponseWriter, r *http.Request) error {
	form := glasses.GlassesRegisterForm(models.GlassesForm{})
	sidebar := h.renderSidebar()
	insertPagePage := glasses.GlassesLayoutPage("Insert glasses", "form to insert new glasses", sidebar, form)
	return h.CreateLayout(w, r, "Insert glasses", insertPagePage).Render(context.Background(), w)
}

func (h *Handler) InsertGlasses(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		HandleError(err, "parsing form")
		return err
	}

	leftVal, err := strconv.ParseFloat(r.FormValue("left_eye_strength"), 64)
	if err != nil {
		return fmt.Errorf("invalid left eye strength: %v", err)
	}

	rightVal, err := strconv.ParseFloat(r.FormValue("right_eye_strength"), 64)
	if err != nil {
		return fmt.Errorf("invalid right eye strength: %v", err)
	}

	g := models.Glasses{
		Reference: r.FormValue("reference"),
		Brand:     r.FormValue("brand"),
		Color:     r.FormValue("color"),
		LeftEye:   leftVal,
		RightEye:  rightVal,
		Type:      r.FormValue("type"),
		Feature:   r.FormValue("features"),
	}

	err = h.service.InsertGlasses(context.Background(), g)
	if err != nil {
		return fmt.Errorf("error inserting glasses: %v", err)
	}

	actionType := r.FormValue("action")

	if actionType == "insert_more" {
		w.Header().Set("HX-Trigger", "glassesAdded")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = fmt.Fprintf(w, `
			<div class="flex items-center">
                <div class="success success-message">Glasses successfully added! You can add another.</div>
                <svg class="w-[18px] h-[18px] ml-2" viewBox="0 0 24 24" fill="green" xmlns="http://www.w3.org/2000/svg">
                    <path fill-rule="evenodd" clip-rule="evenodd" d="M12 2C6.47715 2 2 6.47715 2 12C2 17.5228 6.47715 22 12 22C17.5228 22 22 17.5228 22 12C22 6.47715 17.5228 2 12 2ZM16.7744 9.63269C17.1238 9.20501 17.0604 8.57503 16.6327 8.22559C16.2051 7.87615 15.5751 7.93957 15.2256 8.36725L10.6321 13.9892L8.65936 12.2524C8.24484 11.8874 7.61295 11.9276 7.248 12.3421C6.88304 12.7566 6.92322 13.3885 7.33774 13.7535L9.31046 15.4903C10.1612 16.2393 11.4637 16.1324 12.1808 15.2547L16.7744 9.63269Z" fill="currentColor"></path>
                </svg>
			</div>`)
	} else if actionType == "insert_and_redirect" {
		w.Header().Set("HX-Redirect", "/glasses")
	}

	return nil
}

func (h *Handler) DeleteGlasses(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	glassesIDStr := vars["glasses_id"]
	glassesID, err := uuid.Parse(glassesIDStr)
	if err != nil {
		http.Error(w, "Invalid glasses ID", http.StatusBadRequest)
		return err
	}

	// Delete the glasses
	err = h.service.DeleteGlasses(context.Background(), glassesID)
	if err != nil {
		http.Error(w, "Failed to delete glasses", http.StatusInternalServerError)
		return err
	}

	// Return a success response
	w.Header().Set("HX-Redirect", "/glasses")

	return nil
}

func (h *Handler) UpdateGlassesPage(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	glassesIDStr := vars["glasses_id"]
	glassesID, err := uuid.Parse(glassesIDStr)
	if err != nil {
		http.Error(w, "Invalid glasses ID", http.StatusBadRequest)
		return err
	}

	g, err := h.service.GetGlassesByID(context.Background(), glassesID)
	if err != nil {
		http.Error(w, "Failed to retrieve glasses", http.StatusInternalServerError)
		return err
	}

	le := strconv.FormatFloat(g.LeftEye, 'f', 2, 64)
	re := strconv.FormatFloat(g.RightEye, 'f', 2, 64)
	form := models.GlassesForm{
		Values: map[string]string{
			"Reference": g.Reference,
			"Brand":     g.Brand,
			"LeftEye":   le,
			"RightEye":  re,
			"Color":     g.Color,
			"Type":      g.Type,
			"Features":  g.Feature,
		},
	}

	f := glasses.GlassesUpdateForm(form, glassesIDStr)
	sidebar := h.renderSidebar()
	updatePage := glasses.GlassesLayoutPage("Update Glasses", "form to update glasses", sidebar, f)
	return h.CreateLayout(w, r, "Update Glasses", updatePage).Render(context.Background(), w)
}

func (h *Handler) UpdateGlasses(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	glassesIDStr := vars["glasses_id"]
	glassesID, err := uuid.Parse(glassesIDStr)
	if err != nil {
		http.Error(w, "Invalid glasses ID", http.StatusBadRequest)
		return err
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return err
	}

	g := models.Glasses{
		GlassesID: glassesID,
		Reference: r.FormValue("reference"),
		Brand:     r.FormValue("brand"),
		LeftEye:   parseFloat(r.FormValue("left_eye_strength")),
		RightEye:  parseFloat(r.FormValue("right_eye_strength")),
		Color:     r.FormValue("color"),
		Type:      r.FormValue("type"),
		Feature:   r.FormValue("features"),
	}

	err = h.service.UpdateGlasses(context.Background(), g)
	if err != nil {
		http.Error(w, "Failed to update glasses", http.StatusInternalServerError)
		return err
	}

	w.Header().Set("HX-Redirect", "/glasses")

	return nil
}

func parseFloat(value string) float64 {
	f, _ := strconv.ParseFloat(value, 64)
	return f
}

// Filtered Side bar views

func (h *Handler) getGlassesByType(w http.ResponseWriter, r *http.Request) (int, []models.Glasses, error) {
	pageSize := 10
	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")
	vars := mux.Vars(r)
	filter := vars["type"]
	page, err := strconv.Atoi(r.URL.Query().Get("page"))

	if err != nil {
		page = 1
	}

	g, err := h.service.GetGlassesByType(context.Background(), page, pageSize, orderBy, sortBy, filter)

	if err != nil {
		httperror.ErrNotFound.WriteError(w)
		return 0, nil, err
	}

	return page, g, nil
}

func (h *Handler) renderTypeTable(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	var page int
	var sortAux string
	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")
	brand := r.FormValue("brand")
	vars := mux.Vars(r)
	filter := vars["type"]
	if sortBy == ASC {
		sortAux = DESC
	} else {
		sortAux = ASC
	}

	columnNames := []models.ColumnItems{
		{Title: "Brand", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Color", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Reference", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Left Eye", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Right Eye", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Type", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Has Stock", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Features", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Created At", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Updated At", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
	}

	page, g, _ := h.getGlassesByType(w, r)

	if len(g) == 0 {
		message := components.EmptyPageComponent()
		return message, nil
	}

	nextPage := page + 1
	prevPage := page - 1
	if prevPage <= 1 {
		prevPage = 1
	}

	lastPage, err := h.service.GetSumByType(filter)
	if err != nil {
		HandleError(err, " fetching tax")
		return nil, err
	}
	data := models.GlassesTable{
		Column:      columnNames,
		Glasses:     g,
		PrevPage:    prevPage,
		NextPage:    nextPage,
		Page:        page,
		LastPage:    lastPage,
		FilterBrand: brand,
		OrderParam:  orderBy,
		SortParam:   sortAux,
	}
	taxTable := glasses.GlassesByFilter(data)

	return taxTable, nil
}

func (h *Handler) GlassesTypePage(w http.ResponseWriter, r *http.Request) error {
	sidebar := h.renderSidebar()
	renderTable, err := h.renderTypeTable(w, r)
	if err != nil {
		HandleError(err, " rendering glasses table")
	}
	home := glasses.GlassesLayoutPage("Glasses Management Page", "Glasses Management Page", sidebar, renderTable)
	return h.CreateLayout(w, r, "Glasses Management Page", home).Render(context.Background(), w)
}

// Filter by stock

func (h *Handler) renderInventoryTable(w http.ResponseWriter, r *http.Request, hasStock bool) (templ.Component, error) {
	pageSize := 10
	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	var sortAux string

	if sortBy == ASC {
		sortAux = DESC
	} else {
		sortAux = ASC
	}
	if err != nil {
		page = 1
	}

	g, err := h.service.GetGlassesByStock(context.Background(), page, pageSize, orderBy, sortBy, hasStock)
	if err != nil {
		httperror.ErrNotFound.WriteError(w)
		return nil, err
	}

	columnNames := []models.ColumnItems{
		{Title: "Brand", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Color", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Reference", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Left Eye", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Right Eye", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Type", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Has Stock", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Features", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Created At", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Updated At", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
	}

	if len(g) == 0 {
		message := components.EmptyPageComponent()
		return message, nil
	}

	nextPage := page + 1
	prevPage := page - 1
	if prevPage <= 1 {
		prevPage = 1
	}

	lastPage, err := h.service.GetSumByStock(hasStock)
	if err != nil {
		HandleError(err, " fetching glasses count")
		return nil, err
	}

	data := models.GlassesTable{
		Column:     columnNames,
		Glasses:    g,
		PrevPage:   prevPage,
		NextPage:   nextPage,
		Page:       page,
		LastPage:   lastPage,
		OrderParam: orderBy,
		SortParam:  sortBy,
	}
	glassesTable := glasses.GlassesByFilter(data)

	return glassesTable, nil
}

func (h *Handler) GlassesStockPage(w http.ResponseWriter, r *http.Request) error {
	sidebar := h.renderSidebar()
	var hasStock bool
	vars := mux.Vars(r)
	stock := vars["stock"]

	if stock == "current" {
		hasStock = true
	} else {
		hasStock = false
	}

	renderTable, err := h.renderInventoryTable(w, r, hasStock)

	if err != nil {
		HandleError(err, "rendering glasses table")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	home := glasses.GlassesLayoutPage("Glasses Inventory Management", "Glasses Inventory Management", sidebar, renderTable)
	return h.CreateLayout(w, r, "Glasses Inventory Management", home).Render(context.Background(), w)
}
