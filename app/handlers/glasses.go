package handlers

import (
	"context"
	"fmt"
	"log"
	"log/slog"
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

func (h *Handler) renderGlassesTable(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	_, span := tracer.Start(r.Context(), "renderGlassesTable")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	referenceFilter := r.FormValue("reference")
	leftSphFilter := r.FormValue("left_sph")
	rightSphFilter := r.FormValue("right_sph")

	var page int
	var sortAux string
	var sortIcon templ.Component
	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")
	brand := r.FormValue("brand")
	if sortBy == ASC {
		sortAux = DESC
		sortIcon = svg.ArrowUp()
	} else {
		sortAux = ASC
		sortIcon = svg.ArrowDownIcon()
	}

	columnNames := []models.ColumnItems{
		// General Info
		{Title: "Reference", Icon: sortIcon, SortParam: "reference"},
		{Title: "Brand", Icon: sortIcon, SortParam: "brand"},
		{Title: "Type", Icon: sortIcon, SortParam: "type"},
		{Title: "Color", Icon: sortIcon, SortParam: "color"},

		// Left Eye Prescription
		{Title: "L Sph"},
		{Title: "L Cyl"},
		{Title: "L Axis"},
		{Title: "L Add"},

		// Right Eye Prescription
		{Title: "R Sph"},
		{Title: "R Cyl"},
		{Title: "R Axis"},
		{Title: "R Add"},

		// Status & Details
		{Title: "Stock"}, // Renamed from "Has Stock" for brevity
		{Title: "Features", Icon: sortIcon, SortParam: "feature"}, // Assuming struct field is 'Feature' or db col is 'feature'

		// Timestamps
		{Title: "Created"}, // Shortened title
		// {Title: "Updated At", Icon: svg.ArrowOrderIcon(), SortParam: "updated_at"}, // Often Updated At isn't shown by default unless needed
	}

	page, g, err := h.getGlasses(w, r)

	if err != nil {
		message := glasses.GlassesEmptyPage()
		return message, nil
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

	lastPage, err := h.service.GetSum()
	if err != nil {
		HandleError(err, " fetching glasses")
		return nil, err
	}
	data := models.GlassesTable{
		Column:          columnNames,
		Glasses:         g,
		PrevPage:        prevPage,
		NextPage:        nextPage,
		Page:            page,
		LastPage:        lastPage,
		FilterBrand:     brand,
		OrderParam:      orderBy,
		SortParam:       sortAux,
		FieldErrors:     map[string]string{},
		FilterReference: referenceFilter,
		FilterLeftSph:   leftSphFilter,
		FilterRightSph:  rightSphFilter,
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
	renderTable, err := h.renderGlassesTable(w, r)
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
	ctx := r.Context() // Use request context

	// 1. Get User (keep your existing logic)
	var user *models.UserSession
	userCtx := ctx.Value(models.CtxKeyAuthUser)
	if userCtx != nil {
		switch u := userCtx.(type) {
		case *models.UserSession:
			user = u
		default:
			log.Printf("Unexpected type in userCtx: %T", userCtx)
		}
	}

	// 2. Parse Form
	if err := r.ParseForm(); err != nil {
		slog.ErrorContext(ctx, "Failed to parse insert form", "err", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return fmt.Errorf("parsing form: %w", err)
	}

	// 3. Create the Form Struct FIRST
	form := models.GlassesForm{
		Values:      make(map[string]string), // To repopulate form on error
		FieldErrors: make(map[string]string),
		UserID:      user.ID, // Assign user ID now
	}

	// 4. Populate Values and Typed fields, Validate as you go
	form.Reference = r.FormValue("reference")
	form.Values["reference"] = form.Reference // Keep submitted value
	if form.Reference == "" {
		form.FieldErrors["reference"] = "Reference is required."
	} // Add uniqueness check if needed

	form.Brand = r.FormValue("brand")
	form.Values["brand"] = form.Brand // Keep submitted value

	form.Color = r.FormValue("color")
	form.Values["color"] = form.Color // Keep submitted value

	form.Type = r.FormValue("type")
	form.Values["type"] = form.Type // Keep submitted value
	if form.Type != "adult" && form.Type != "children" {
		form.FieldErrors["type"] = "Please select a valid type."
	}

	form.Feature = r.FormValue("features")
	form.Values["features"] = form.Feature // Keep submitted value

	// Parse numeric fields, add errors directly to form.FieldErrors
	form.Values["left_sph"] = r.FormValue("left_sph") // Store original string value
	form.LeftSph = parseFloatPointerAndError(form.Values["left_sph"], "left_sph", "Invalid number format for Left Sphere.", form.FieldErrors)

	form.Values["left_cyl"] = r.FormValue("left_cyl")
	form.LeftCyl = parseFloatPointerAndError(form.Values["left_cyl"], "left_cyl", "Invalid number format for Left Cylinder.", form.FieldErrors)

	form.Values["left_axis"] = r.FormValue("left_axis")
	form.LeftAxis = parseFloatPointerAndError(form.Values["left_axis"], "left_axis", "Invalid number format for Left Axis.", form.FieldErrors) // Adapt if Axis is int

	form.Values["left_add"] = r.FormValue("left_add")
	form.LeftAdd = parseFloatPointerAndError(form.Values["left_add"], "left_add", "Invalid number format for Left Addition.", form.FieldErrors)

	form.Values["right_sph"] = r.FormValue("right_sph")
	form.RightSph = parseFloatPointerAndError(form.Values["right_sph"], "right_sph", "Invalid number format for Right Sphere.", form.FieldErrors)

	form.Values["right_cyl"] = r.FormValue("right_cyl")
	form.RightCyl = parseFloatPointerAndError(form.Values["right_cyl"], "right_cyl", "Invalid number format for Right Cylinder.", form.FieldErrors)

	form.Values["right_axis"] = r.FormValue("right_axis")
	form.RightAxis = parseFloatPointerAndError(form.Values["right_axis"], "right_axis", "Invalid number format for Right Axis.", form.FieldErrors) // Adapt if Axis is int

	form.Values["right_add"] = r.FormValue("right_add")
	form.RightAdd = parseFloatPointerAndError(form.Values["right_add"], "right_add", "Invalid number format for Right Addition.", form.FieldErrors)

	// 5. Check if any errors occurred during parsing/validation
	if len(form.FieldErrors) > 0 {
		slog.WarnContext(ctx, "Validation errors on glasses insert form", "errors", form.FieldErrors)

		// --- Re-render the MODAL component with errors ---
		// Use status 422 Unprocessable Entity for validation errors
		w.WriteHeader(http.StatusUnprocessableEntity)
		// Render the modal component, passing the form with Values and Errors
		component := glasses.GlassesInsertModal(form)
		// Render directly to the response writer
		// HTMX will swap this into the target specified on the form/button
		return component.Render(ctx, w) // Use request context
	}

	// 6. If validation passed, call the service
	slog.InfoContext(ctx, "Validation passed, attempting to insert glasses", "reference", form.Reference)
	err := h.service.InsertGlasses(ctx, form) // Use request context
	if err != nil {
		slog.ErrorContext(ctx, "Failed to insert glasses via service", "err", err)
		// You could potentially parse the error and add to FieldErrors
		// and re-render the form again, or just return a server error
		http.Error(w, "Failed to save glasses data.", http.StatusInternalServerError)
		return err // Return the error from the service call
	}

	// 7. Success: Tell HTMX to redirect (or close modal and refresh table)
	slog.InfoContext(ctx, "Successfully inserted glasses", "reference", form.Reference) // Use the form's reference
	w.Header().Set("HX-Redirect", "/glasses")                                           // Standard redirect
	// OR, to close modal and refresh table (assuming table has id="glasses-table"):
	// w.Header().Set("HX-Trigger", `{"closeModal": "true", "refreshTable": "#glasses-table"}`)
	// Your Alpine modal would need to listen for `closeModal` event
	// Your table container would need hx-get="/path/to/table/content" hx-trigger="refreshTable from:body"

	w.WriteHeader(http.StatusOK) // Important for HX-Redirect to be processed
	return nil                   // Indicate success
}

func parseFloatPointerAndError(valueStr string, fieldName string, errorMsg string, errors map[string]string) *float64 {
	if valueStr == "" {
		return nil // Treat empty string as valid null value
	}
	f, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		errors[fieldName] = errorMsg
		return nil
	}
	return &f
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
	// Use request context for tracing and potentially database calls
	ctx, span := tracer.Start(r.Context(), "UpdateGlassesPageHandler")
	defer span.End()

	vars := mux.Vars(r)
	glassesIDStr := vars["glasses_id"]
	glassesID, err := uuid.Parse(glassesIDStr)
	if err != nil {
		slog.ErrorContext(ctx, "Invalid glasses ID format in URL", "id", glassesIDStr, "err", err)
		http.Error(w, "Invalid glasses ID", http.StatusBadRequest)
		return fmt.Errorf("parsing glasses ID: %w", err) // Return error for handler wrapper
	}

	// 1. Fetch the current glasses data
	// Use ctx from request for service call
	currentGlasses, err := h.service.GetGlassesByID(ctx, glassesID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get glasses by ID for edit page", "err", err, "glasses_id", glassesIDStr)
		// Handle appropriately - maybe show a "not found" page
		http.Error(w, "Glasses not found", http.StatusNotFound)
		return fmt.Errorf("getting glasses by ID: %w", err)
	}

	// 2. Create a new GlassesForm and POPULATE its Values map fully
	form := models.GlassesForm{
		Values:      make(map[string]string),
		FieldErrors: make(map[string]string),
		GlassesID:   glassesID,
	}

	// --- Populate the Values map from currentGlasses data ---
	form.Values["reference"] = currentGlasses.Reference
	form.Values["brand"] = currentGlasses.Brand
	form.Values["color"] = currentGlasses.Color
	form.Values["type"] = currentGlasses.Type
	form.Values["features"] = currentGlasses.Feature // Ensure Feature is the correct field name

	// Populate prescription values (handle nil pointers!)
	// Left Eye
	if currentGlasses.LeftPrescription.Sph != nil {
		form.Values["left_sph"] = fmt.Sprintf("%.2f", *currentGlasses.LeftPrescription.Sph)
	} else {
		form.Values["left_sph"] = "" // Or "0.00" if a default is preferred
	}
	if currentGlasses.LeftPrescription.Cyl != nil {
		form.Values["left_cyl"] = fmt.Sprintf("%.2f", *currentGlasses.LeftPrescription.Cyl)
	} else {
		form.Values["left_cyl"] = ""
	}
	if currentGlasses.LeftPrescription.Axis != nil {
		// Assuming Axis is *float64 based on previous fixes, format as integer
		form.Values["left_axis"] = fmt.Sprintf("%.0f", *currentGlasses.LeftPrescription.Axis)
		// If Axis is *int type: form.Values["left_axis"] = fmt.Sprintf("%d", *currentGlasses.LeftPrescription.Axis)
	} else {
		form.Values["left_axis"] = ""
	}
	if currentGlasses.LeftPrescription.Add != nil {
		form.Values["left_add"] = fmt.Sprintf("%.2f", *currentGlasses.LeftPrescription.Add)
	} else {
		form.Values["left_add"] = ""
	}

	// Right Eye
	if currentGlasses.RightPrescription.Sph != nil {
		form.Values["right_sph"] = fmt.Sprintf("%.2f", *currentGlasses.RightPrescription.Sph)
	} else {
		form.Values["right_sph"] = ""
	}
	if currentGlasses.RightPrescription.Cyl != nil {
		form.Values["right_cyl"] = fmt.Sprintf("%.2f", *currentGlasses.RightPrescription.Cyl)
	} else {
		form.Values["right_cyl"] = ""
	}
	if currentGlasses.RightPrescription.Axis != nil {
		// Assuming Axis is *float64 based on previous fixes, format as integer
		form.Values["right_axis"] = fmt.Sprintf("%.0f", *currentGlasses.RightPrescription.Axis)
		// If Axis is *int type: form.Values["right_axis"] = fmt.Sprintf("%d", *currentGlasses.RightPrescription.Axis)
	} else {
		form.Values["right_axis"] = ""
	}
	if currentGlasses.RightPrescription.Add != nil {
		form.Values["right_add"] = fmt.Sprintf("%.2f", *currentGlasses.RightPrescription.Add)
	} else {
		form.Values["right_add"] = ""
	}

	formComponent := glasses.GlassesUpdateForm(*currentGlasses, form, glassesIDStr)
	sidebar := h.renderSidebar() // Assuming this returns a templ.Component

	pageComponent := pages.MainLayoutPage("Update Glasses", "Form to update glasses", sidebar, formComponent)

	return pageComponent.Render(r.Context(), w)
}

func (h *Handler) UpdateGlasses(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	glassesIDStr := vars["glasses_id"]
	glassesID, err := uuid.Parse(glassesIDStr)
	if err != nil {
		slog.Error("Invalid glasses ID format in URL", "id", glassesIDStr, "err", err)
		http.Error(w, "Invalid glasses ID", http.StatusBadRequest)
		return fmt.Errorf("parsing glasses ID: %w", err)
	}

	// 1. Fetch current glasses data (needed for stock check and potentially re-rendering form)
	currentGlasses, err := h.service.GetGlassesByID(r.Context(), glassesID)
	if err != nil {
		slog.Error("Failed to get glasses by ID for update", "err", err, "glasses_id", glassesIDStr)
		http.Error(w, "Glasses not found", http.StatusNotFound)
		return fmt.Errorf("getting glasses by ID: %w", err)
	}

	// 2. Parse form data
	if err := r.ParseForm(); err != nil {
		slog.Error("Failed to parse update form", "err", err, "glasses_id", glassesIDStr)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return fmt.Errorf("parsing form: %w", err)
	}

	// 3. Populate the GlassesForm struct AND Validate
	form := models.GlassesForm{
		GlassesID: glassesID, // Store the parsed UUID
		Reference: r.FormValue("reference"),
		Brand:     r.FormValue("brand"),
		Color:     r.FormValue("color"),
		Type:      r.FormValue("type"), // Get submitted type
		Feature:   r.FormValue("features"),
		// --- Parse numeric fields robustly ---
		LeftSph:   parseFloatPointer(r.FormValue("left_sph")), // Use helper for pointers
		LeftCyl:   parseFloatPointer(r.FormValue("left_cyl")),
		LeftAxis:  parseFloatPointer(r.FormValue("left_axis")), // Assuming float64 for Axis now
		LeftAdd:   parseFloatPointer(r.FormValue("left_add")),
		RightSph:  parseFloatPointer(r.FormValue("right_sph")),
		RightCyl:  parseFloatPointer(r.FormValue("right_cyl")),
		RightAxis: parseFloatPointer(r.FormValue("right_axis")), // Assuming float64 for Axis now
		RightAdd:  parseFloatPointer(r.FormValue("right_add")),
		// --- End numeric parsing ---
		Values:      make(map[string]string), // Keep for re-rendering form if needed
		FieldErrors: make(map[string]string),
	}

	isValidType := false
	if form.Type == "adult" || form.Type == "children" {
		isValidType = true
	}
	if !isValidType {
		form.FieldErrors["type"] = "Please select a valid type (Adult or Children)."
	}

	// Add other validations as needed (e.g., reference check)
	// ref, err := h.service.GetGlassesReference(context.Background(), glassesID)
	// if err != nil {
	// 	http.Error(w, "Failed to retrieve glasses reference", http.StatusInternalServerError)
	// 	return err
	// }
	// if ref != form.Reference { // Check only if reference *changed*
	//   // Check if the *new* reference already exists for another pair
	//   exists, errCheck := h.service.CheckReferenceExists(r.Context(), form.Reference, glassesID) // Need a service method like this
	//   if errCheck != nil {
	//      http.Error(w, "Failed checking reference uniqueness", http.StatusInternalServerError)
	//      return errCheck
	//   }
	//   if exists {
	//       form.FieldErrors["reference"] = "This reference is already used by another pair of glasses."
	//   }
	// }

	// --- Repopulate form.Values for re-rendering if there are errors ---
	if len(form.FieldErrors) > 0 {
		slog.Warn("Validation errors found during update", "errors", form.FieldErrors, "glasses_id", glassesIDStr)
		// Populate form.Values from the *submitted* (potentially invalid) data
		form.Values["reference"] = form.Reference
		form.Values["brand"] = form.Brand
		form.Values["color"] = form.Color
		form.Values["type"] = form.Type // Use the submitted type for re-selection
		form.Values["features"] = form.Feature
		// Repopulate numeric values (might need nil checks if parseFloatPointer returns nil)
		if form.LeftSph != nil {
			form.Values["left_sph"] = fmt.Sprintf("%.2f", *form.LeftSph)
		} else {
			form.Values["left_sph"] = ""
		}
		if form.LeftCyl != nil {
			form.Values["left_cyl"] = fmt.Sprintf("%.2f", *form.LeftCyl)
		} else {
			form.Values["left_cyl"] = ""
		}
		if form.LeftAxis != nil {
			form.Values["left_axis"] = fmt.Sprintf("%.0f", *form.LeftAxis)
		} else {
			form.Values["left_axis"] = ""
		} // Adjust format if int
		if form.LeftAdd != nil {
			form.Values["left_add"] = fmt.Sprintf("%.2f", *form.LeftAdd)
		} else {
			form.Values["left_add"] = ""
		}
		if form.RightSph != nil {
			form.Values["right_sph"] = fmt.Sprintf("%.2f", *form.RightSph)
		} else {
			form.Values["right_sph"] = ""
		}
		if form.RightCyl != nil {
			form.Values["right_cyl"] = fmt.Sprintf("%.2f", *form.RightCyl)
		} else {
			form.Values["right_cyl"] = ""
		}
		if form.RightAxis != nil {
			form.Values["right_axis"] = fmt.Sprintf("%.0f", *form.RightAxis)
		} else {
			form.Values["right_axis"] = ""
		} // Adjust format if int
		if form.RightAdd != nil {
			form.Values["right_add"] = fmt.Sprintf("%.2f", *form.RightAdd)
		} else {
			form.Values["right_add"] = ""
		}

		// Re-render the form with errors and previously entered values
		component := glasses.GlassesUpdateForm(*currentGlasses, form, glassesIDStr) // Pass currentGlasses for restock check
		w.WriteHeader(http.StatusUnprocessableEntity)                               // Use 422 for validation errors
		return component.Render(r.Context(), w)
		// Note: Rendering directly might not work perfectly with HTMX depending on swap targets.
		// You might need specific HTMX responses for error handling.
	}

	// --- Proceed if validation passed ---
	restockFlag := r.FormValue("restock") == "true"

	// Call the service/repository update function
	if err = h.service.UpdateGlasses(r.Context(), form, restockFlag, currentGlasses.IsInStock); err != nil {
		slog.Error("Failed to update glasses in repository", "err", err, "glasses_id", glassesIDStr, "restock", restockFlag)
		http.Error(w, "Failed to update glasses", http.StatusInternalServerError)
		return err
	}

	slog.Info("Successfully processed glasses update", "glasses_id", glassesIDStr, "restock", restockFlag)

	// Redirect or respond
	w.Header().Set("HX-Redirect", "/glasses")
	w.WriteHeader(http.StatusOK) // Or StatusNoContent
	return nil
}

func parseFloat(value string) float64 {
	f, _ := strconv.ParseFloat(value, 64)
	return f
}

func parseFloatPointer(valueStr string) *float64 {
	if valueStr == "" {
		return nil // Handle empty string as nil
	}
	f, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return nil // Handle parse errors as nil (or log/handle differently)
	}
	return &f
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
