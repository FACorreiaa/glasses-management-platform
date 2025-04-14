package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	httperror "github.com/FACorreiaa/glasses-management-platform/app/errors"
	"github.com/FACorreiaa/glasses-management-platform/app/models"
	"github.com/FACorreiaa/glasses-management-platform/app/static/svg"
	"github.com/FACorreiaa/glasses-management-platform/app/view/glasses"
	"github.com/FACorreiaa/glasses-management-platform/app/view/pages"

	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var tracer = otel.Tracer("github.com/FACorreiaa/glasses-management-platform/services")

func (h *Handler) renderSidebar() []models.SidebarItem {
	sidebar := []models.SidebarItem{
		{Path: "/", Label: "Home"},
		{Path: "/glasses", Label: "Glasses stock"},
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
		{Path: "/logout", Label: "Log out"},
	}
	return sidebar
}

func (h *Handler) getGlasses(w http.ResponseWriter, r *http.Request) (int, []models.Glasses, error) {
	ctx, span := tracer.Start(r.Context(), "getGlassesHandler") // Use request context!
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	pageSize := 10
	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	reference := r.FormValue("reference")

	leftEyeStr := r.FormValue("left_sph")
	rightEyeStr := r.FormValue("right_sph")

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

	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("orderBy", orderBy),
		attribute.String("sortBy", sortBy),
	)

	g, err := h.service.GetGlasses(ctx, page, pageSize, orderBy, sortBy, reference, leftEye, rightEye)
	if err != nil {
		span.RecordError(err) // Record errors on the span
		span.SetStatus(codes.Error, err.Error())
		httperror.ErrNotFound.WriteError(w)
		return 0, nil, err
	}

	span.SetAttributes(attribute.Int("glasses.count", len(g)))

	return page, g, nil
}

func (h *Handler) renderGlassesTable(ctx context.Context, w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	ctx, span := tracer.Start(r.Context(), "renderGlassesTable")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

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
		// General Info
		{Title: "Reference", Icon: svg.ArrowOrderIcon(), SortParam: "reference"},
		{Title: "Brand", Icon: svg.ArrowOrderIcon(), SortParam: "brand"},
		{Title: "Type", Icon: svg.ArrowOrderIcon(), SortParam: "type"},
		{Title: "Color", Icon: svg.ArrowOrderIcon(), SortParam: "color"},

		// Left Eye Prescription
		{Title: "L Sph", Icon: svg.ArrowOrderIcon(), SortParam: "left_sph"},
		{Title: "L Cyl", Icon: svg.ArrowOrderIcon(), SortParam: "left_cyl"},
		{Title: "L Axis", Icon: svg.ArrowOrderIcon(), SortParam: "left_axis"},
		{Title: "L Add", Icon: svg.ArrowOrderIcon(), SortParam: "left_add"},

		// Right Eye Prescription
		{Title: "R Sph", Icon: svg.ArrowOrderIcon(), SortParam: "right_sph"},
		{Title: "R Cyl", Icon: svg.ArrowOrderIcon(), SortParam: "right_cyl"},
		{Title: "R Axis", Icon: svg.ArrowOrderIcon(), SortParam: "right_axis"},
		{Title: "R Add", Icon: svg.ArrowOrderIcon(), SortParam: "right_add"},

		// Status & Details
		{Title: "Stock", Icon: svg.ArrowOrderIcon(), SortParam: "is_in_stock"}, // Renamed from "Has Stock" for brevity
		{Title: "Features", Icon: svg.ArrowOrderIcon(), SortParam: "feature"},  // Assuming struct field is 'Feature' or db col is 'feature'

		// Timestamps
		{Title: "Created", Icon: svg.ArrowOrderIcon(), SortParam: "created_at"}, // Shortened title
		// {Title: "Updated At", Icon: svg.ArrowOrderIcon(), SortParam: "updated_at"}, // Often Updated At isn't shown by default unless needed
	}

	page, g, _ := h.getGlasses(w, r)

	if len(g) == 0 {
		message := glasses.GlassesEmptyPage()
		return message, nil
	}

	nextPage := page + 1
	prevPage := page - 1
	if prevPage <= 1 {
		prevPage = 1
	}

	lastPage, err := h.service.GetSum()
	if err != nil {
		HandleError(err, " fetching glasses")
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
		FieldErrors: map[string]string{},
	}

	t := glasses.GlassesTable(data, models.GlassesForm{})

	return t, nil
}

func (h *Handler) GlassesPage(w http.ResponseWriter, r *http.Request) error {
	ctx, span := tracer.Start(r.Context(), "GlassesPageHandler")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	sidebar := h.renderSidebar()
	renderTable, err := h.renderGlassesTable(ctx, w, r)
	//errorPage := error.ErrorPage() // Example

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		HandleError(err, "rendering glasses table")
		// Decide how to handle the error response (e.g., render an error page)
		// You might want to render an error component instead of just logging
		//return h.CreateLayout(ctx, w, r, "Error", errorPage).Render(ctx, w) // Render error page
	}
	home := pages.MainLayoutPage("Glasses Management Page", "Glasses Management Page", sidebar, renderTable)
	return h.CreateLayout(ctx, w, r, "Glasses Management Page", home).Render(context.Background(), w)
}

func (h *Handler) GlassesRegisterPage(w http.ResponseWriter, r *http.Request) error {
	ctx, span := tracer.Start(r.Context(), "GlassesRegisterPageHandler")
	defer span.End()

	form := glasses.GlassesRegisterForm(models.GlassesForm{})
	sidebar := h.renderSidebar()
	insertPagePage := pages.MainLayoutPage("Insert glasses", "form to insert new glasses", sidebar, form)
	return h.CreateLayout(ctx, w, r, "Insert glasses", insertPagePage).Render(context.Background(), w)
}

func (h *Handler) InsertGlasses(w http.ResponseWriter, r *http.Request) error {
	var user *models.UserSession
	userCtx := r.Context().Value(models.CtxKeyAuthUser)

	fieldError := make(map[string]string)
	if userCtx != nil {
		switch u := userCtx.(type) {
		case *models.UserSession:
			user = u
		default:
			log.Printf("Unexpected type in userCtx: %T", userCtx)
		}
	}

	if err := r.ParseForm(); err != nil {
		HandleError(err, "parsing form")
		return err
	}

	leftSph, err := strconv.ParseFloat(r.FormValue("left_sph"), 64)
	if err != nil {
		fieldError["left_sph"] = "invalid left eye strength"
	}

	rightSph, err := strconv.ParseFloat(r.FormValue("left_sph"), 64)
	if err != nil {
		fieldError["left_sph"] = "invalid lerightft eye strength"
	}

	leftCyl, err := strconv.ParseFloat(r.FormValue("left_cyl"), 64)
	if err != nil {
		fieldError["left_cyl"] = "invalid left eye cylinder"
	}
	rightCyl, err := strconv.ParseFloat(r.FormValue("right_cyl"), 64)
	if err != nil {
		fieldError["right_cyl"] = "invalid right eye cylinder"
	}
	leftAxis, err := strconv.ParseFloat(r.FormValue("left_axis"), 64)
	if err != nil {
		fieldError["left_axis"] = "invalid left eye axis"
	}
	rightAxis, err := strconv.ParseFloat(r.FormValue("right_axis"), 64)
	if err != nil {
		fieldError["right_axis"] = "invalid right eye axis"
	}
	leftAdd, err := strconv.ParseFloat(r.FormValue("left_add"), 64)
	if err != nil {
		fieldError["left_add"] = "invalid left eye add"
	}
	rightAdd, err := strconv.ParseFloat(r.FormValue("right_add"), 64)
	if err != nil {
		fieldError["right_add"] = "invalid right eye add"
	}

	g := models.GlassesForm{
		Reference:   r.FormValue("reference"),
		Brand:       r.FormValue("brand"),
		LeftSph:     leftSph,
		LeftCyl:     leftCyl,
		LeftAxis:    leftAxis,
		LeftAdd:     leftAdd,
		RightSph:    rightSph,
		RightCyl:    rightCyl,
		RightAxis:   rightAxis,
		RightAdd:    rightAdd,
		Color:       r.FormValue("color"),
		Type:        r.FormValue("type"),
		Feature:     r.FormValue("features"),
		UserID:      user.ID,
		FieldErrors: make(map[string]string),
	}

	if len(g.Reference) == 0 {
		g.FieldErrors["reference"] = "reference cannot be empty"
	}

	if len(g.Type) == 0 {
		g.FieldErrors["type"] = "type cannot be empty"
	}

	// TO DO

	if len(g.FieldErrors) > 0 {
		return err
	}

	if err = h.service.InsertGlasses(context.Background(), g); err != nil {
		HandleError(err, "inserting glasses")
		return err
	}

	w.Header().Set("HX-Redirect", "/glasses")

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
	if err = h.service.DeleteGlasses(context.Background(), glassesID); err != nil {
		http.Error(w, "Failed to delete glasses", http.StatusInternalServerError)
		return err
	}

	w.Header().Set("HX-Redirect", "/glasses")

	return nil
}

func (h *Handler) UpdateGlassesPage(w http.ResponseWriter, r *http.Request) error {
	ctx, span := tracer.Start(r.Context(), "UpdateGlassesPageHandler")
	defer span.End()

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

	ls := strconv.FormatFloat(*g.LeftPrescription.Sph, 'f', 2, 64)
	rs := strconv.FormatFloat(*g.RightPrescription.Sph, 'f', 2, 64)
	form := models.GlassesForm{
		Values: map[string]string{
			"Reference": g.Reference,
			"Brand":     g.Brand,
			"LeftSph":   ls,
			"RightSph":  rs,
			"Color":     g.Color,
			"Type":      g.Type,
			"Features":  g.Feature,
		},
	}

	f := glasses.GlassesUpdateForm(form, glassesIDStr)
	sidebar := h.renderSidebar()
	updatePage := pages.MainLayoutPage("Update Glasses", "form to update glasses", sidebar, f)
	return h.CreateLayout(ctx, w, r, "Update Glasses", updatePage).Render(context.Background(), w)
}

func (h *Handler) UpdateGlasses(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	glassesIDStr := vars["glasses_id"]
	glassesID, err := uuid.Parse(glassesIDStr)
	if err != nil {
		http.Error(w, "Invalid glasses ID", http.StatusBadRequest)
		return err
	}

	if err = r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return err
	}

	g := models.GlassesForm{
		GlassesID:   glassesID,
		Reference:   r.FormValue("reference"),
		Brand:       r.FormValue("brand"),
		LeftSph:     parseFloat(r.FormValue("left_sph")),
		LeftCyl:     parseFloat(r.FormValue("left_cyl")),
		LeftAxis:    parseFloat(r.FormValue("left_axis")),
		LeftAdd:     parseFloat(r.FormValue("left_add")),
		RightSph:    parseFloat(r.FormValue("right_sph")),
		RightCyl:    parseFloat(r.FormValue("right_cyl")),
		RightAxis:   parseFloat(r.FormValue("right_axis")),
		RightAdd:    parseFloat(r.FormValue("right_add")),
		Color:       r.FormValue("color"),
		Type:        r.FormValue("type"),
		Feature:     r.FormValue("features"),
		FieldErrors: map[string]string{},
	}

	fmt.Println("\n", g.Values)

	ref, err := h.service.GetGlassesReference(context.Background(), glassesID)
	if err != nil {
		http.Error(w, "Failed to retrieve glasses", http.StatusInternalServerError)
		return err
	}

	fmt.Println("\n", ref, g.Reference)

	if ref == g.Reference {
		g.FieldErrors["reference"] = "Glasses reference must be unique"
	}

	if len(g.FieldErrors) > 0 {
		form := glasses.GlassesUpdateForm(g, glassesIDStr).Render(context.Background(), w)
		return form
	}

	if err = h.service.UpdateGlasses(context.Background(), g); err != nil {
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
		message := glasses.GlassesEmptyPage()
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
	t := glasses.GlassesByFilter(data)

	return t, nil
}

func (h *Handler) GlassesTypePage(w http.ResponseWriter, r *http.Request) error {
	ctx, span := tracer.Start(r.Context(), "GlassesTypePageHandler")
	defer span.End()

	sidebar := h.renderSidebar()
	renderTable, err := h.renderTypeTable(w, r)
	if err != nil {
		HandleError(err, " rendering glasses table")
	}
	home := pages.MainLayoutPage("Glasses Management Page", "Glasses Management Page", sidebar, renderTable)
	return h.CreateLayout(ctx, w, r, "Glasses Management Page", home).Render(context.Background(), w)
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
		{Title: "L Sph", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "L Cyl", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "L Axis", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "L Add", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},

		// Right Eye Prescription
		{Title: "R Sph", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "R Cyl", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "R Axis", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "R Add", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Type", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Has Stock", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Features", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Created At", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Updated At", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
	}

	if len(g) == 0 {
		message := glasses.GlassesEmptyPage()
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
	ctx, span := tracer.Start(r.Context(), "GlassesStockPageHandler")
	defer span.End()

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

	home := pages.MainLayoutPage("Glasses Inventory Management", "Glasses Inventory Management", sidebar, renderTable)
	return h.CreateLayout(ctx, w, r, "Glasses Inventory Management", home).Render(context.Background(), w)
}
