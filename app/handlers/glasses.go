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
)

func (h *Handler) renderSidebar() []models.SidebarItem {
	sidebar := []models.SidebarItem{
		{Path: "/", Label: "Home"},
		{Path: "/glasses/portal", Label: "Insert"},
		{Path: "/glasses/portal", Label: "Delete"},
		{Path: "/glasses/portal", Label: "Update"},
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
	println(pageSize)
	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	//brand := r.FormValue("brand")

	if err != nil {
		page = 1
	}

	g, err := h.service.GetGlasses(context.Background(), page, pageSize, orderBy, sortBy)
	fmt.Printf("Glasses: %v\n", g)

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
		{Title: "Created At", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Updated At", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Action", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
	}

	page, g, _ := h.getGlasses(w, r)
	fmt.Printf("Glasses: %v\n Page: %v", g, page)

	nextPage := page + 1
	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 1
	}

	lastPage, err := h.service.GetSum()
	fmt.Printf("Last Page: %v\n", lastPage)
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
