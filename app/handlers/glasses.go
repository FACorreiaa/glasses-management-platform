package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/FACorreiaa/glasses-management-platform/app/models"
	"github.com/FACorreiaa/glasses-management-platform/app/static/svg"
	"github.com/FACorreiaa/glasses-management-platform/app/view/glasses"
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (h *Handler) renderSidebar() []models.SidebarItem {
	sidebar := []models.SidebarItem{
		{Path: "/", Label: "Home"},
		{Path: "/glasses", Label: "Glasses stock"},
		{Path: "/glasses/register", Label: "Insert"},
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

func (h *Handler) getGlasses(_ http.ResponseWriter, r *http.Request) (int, []models.Glasses, error) {
	pageSize := 20
	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))

	if err != nil {
		page = 1
	}

	g, err := h.service.GetGlasses(context.Background(), page, pageSize, orderBy, sortBy)

	if err != nil {
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
		{Title: "Created At", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Updated At", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Action", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
	}

	page, g, _ := h.getGlasses(w, r)

	nextPage := page + 1
	prevPage := page - 1
	if prevPage <= 1 {
		prevPage = 1
	}

	lastPage, err := h.service.GetSum()
	if err != nil {
		HandleError(err, "Error fetching tax")
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
		HandleError(err, " rendering glasses table")
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
	}

	err = h.service.InsertGlasses(context.Background(), g)
	if err != nil {
		return fmt.Errorf("error inserting glasses: %v", err)
	}

	actionType := r.FormValue("action")

	if actionType == "insert_more" {
		w.Header().Set("HX-Trigger", "glassesAdded")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `
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
	w.WriteHeader(http.StatusNoContent)
	return nil
}
