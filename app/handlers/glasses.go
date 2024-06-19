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

//func (h *Handler) GlassesPost(w http.ResponseWriter, r *http.Request) error {
//	values := r.Form
//	leftVal, err := strconv.ParseFloat(values.Get("left_eye_strength"), 32)
//	if err != nil {
//		HandleError(err, "parse left eye strength")
//	}
//
//	rightVal, err := strconv.ParseFloat(values.Get("right_eye_strength"), 32)
//	if err != nil {
//		HandleError(err, "parse right eye strength")
//	}
//	g := models.Glasses{
//		Brand:     values.Get("brand"),
//		Color:     values.Get("color"),
//		Reference: values.Get("reference"),
//		LeftEye:   leftVal,
//		RightEye:  rightVal,
//		Type:      values.Get("type"),
//	}
//	err = h.service.InsertGlasses(context.Background(), g)
//	if err != nil {
//		HandleError(err, "Inserting glasses")
//	}
//	return h.GlassesPage(w, r)
//}

func (h *Handler) GlassesPost(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
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

	http.Redirect(w, r, "/glasses", http.StatusSeeOther)
	return nil
}
