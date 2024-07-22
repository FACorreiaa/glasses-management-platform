package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"context"

	httperror "github.com/FACorreiaa/glasses-management-platform/app/errors"
	"github.com/FACorreiaa/glasses-management-platform/app/models"
	"github.com/FACorreiaa/glasses-management-platform/app/static/svg"
	"github.com/FACorreiaa/glasses-management-platform/app/view/admin"
	"github.com/FACorreiaa/glasses-management-platform/app/view/components"
	"github.com/FACorreiaa/glasses-management-platform/app/view/glasses"
	"github.com/FACorreiaa/glasses-management-platform/app/view/pages"
	"github.com/FACorreiaa/glasses-management-platform/app/view/settings"
	"github.com/FACorreiaa/glasses-management-platform/app/view/shipping"
	"github.com/a-h/templ"
	"github.com/gorilla/mux"
)

func (h *Handler) renderSettingsSidebar() []models.SidebarItem {
	sidebar := []models.SidebarItem{
		{Path: "/settings", Label: "Home"},
		{Path: "/settings/admin", Label: "Change details"},
		{Path: "/settings/collaborators", Label: "View collaborators"},
		{Path: "/settings/glasses", Label: "View glasses stock"},
		{Path: "/settings/shipping", Label: "View transactions"},
		{Path: "/log-out", Label: "Log out"},
	}
	return sidebar
}

func (h *Handler) getGlassesDetails(w http.ResponseWriter, r *http.Request) (int, []models.Glasses, error) {
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

	g, err := h.service.GetGlassesDetails(context.Background(), page, pageSize, orderBy, sortBy, reference, leftEye, rightEye)
	if err != nil {
		httperror.ErrNotFound.WriteError(w)
		return 0, nil, err
	}

	return page, g, nil
}

func (h *Handler) renderGlassesTableDetails(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
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

	//
	columnNames := []models.ColumnItems{
		{Title: "Username", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Email", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Left Eye", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Right Eye", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Reference", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Has Stock", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Created At", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Updated At", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
	}

	page, g, _ := h.getGlassesDetails(w, r)

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
	t := settings.GlassesTable(data)

	return t, nil
}

func (h *Handler) UpdateAdmin(w http.ResponseWriter, r *http.Request) error {
	var user *models.UserSession

	userCtx := r.Context().Value(models.CtxKeyAuthUser)
	if userCtx != nil {
		switch u := userCtx.(type) {
		case *models.UserSession:
			user = u
		default:
			log.Printf("Unexpected type in userCtx: %T", userCtx)
		}
	}

	if user == nil || user.Role != "admin" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return nil
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return err
	}

	g := models.UpdateUserForm{
		UserID:          user.ID,
		Email:           r.FormValue("email"),
		Username:        r.FormValue("username"),
		Password:        r.FormValue("password"),
		PasswordConfirm: r.FormValue("password_confirm"),
	}

	err = h.service.UpdateUser(context.Background(), g)
	if err != nil {
		if err.Error() == "username already exists" || err.Error() == "email already exists" {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
		}
		return err
	}

	w.Header().Set("HX-Redirect", "/settings/admin")
	return nil
}

func (h *Handler) UpdateAdminPage(w http.ResponseWriter, r *http.Request) error {
	var user *models.UserSession
	userCtx := r.Context().Value(models.CtxKeyAuthUser)
	if userCtx != nil {
		switch u := userCtx.(type) {
		case *models.UserSession:
			user = u
		default:
			log.Printf("Unexpected type in userCtx: %T", userCtx)
		}
	}

	g, err := h.service.GetAdminID(context.Background(), user.ID)
	if err != nil {
		log.Printf("Failed to retrieve user information: %v", err)
		http.Error(w, "Failed to retrieve glasses", http.StatusInternalServerError)
		return err
	}

	if g == nil {
		log.Printf("No user found for given ID")
		http.Error(w, "No user found for given ID", http.StatusNotFound)
		return nil
	}

	form := models.UpdateUserForm{
		Values: map[string]string{
			"Username": g.Username,
			"Email":    g.Email,
		},
	}

	f := settings.AdminUpdateForm(form, user.ID)
	updatePage := admin.UserLayoutPage("Update users", "form to update users", f)
	return h.CreateLayout(w, r, "Update Glasses", updatePage).Render(context.Background(), w)
}

func (h *Handler) SettingsPage(w http.ResponseWriter, r *http.Request) error {
	s := settings.AdminSettingsPage(models.SettingsPage{})
	page := settings.AdminSettingsLayoutPage("Settings Page", "Admin settings main page", h.renderSettingsSidebar(), s)
	data := h.CreateLayout(w, r, "Settings", page).Render(context.Background(), w)
	return data
}

func (h *Handler) SettingsGlassesPage(w http.ResponseWriter, r *http.Request) error {
	sidebar := h.renderSettingsSidebar()
	renderTable, err := h.renderGlassesTableDetails(w, r)
	if err != nil {
		HandleError(err, "rendering glasses table")
	}
	home := pages.MainLayoutPage("Glasses Management Page", "Check glasses stock details", sidebar, renderTable)
	return h.CreateLayout(w, r, "Glasses Management Page", home).Render(context.Background(), w)
}

// SETTINGS

func (h *Handler) getSettingsShipping(w http.ResponseWriter, r *http.Request) (int, []models.SettingsShippingDetails, error) {
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

	s, err := h.service.GetShippingExpandedDetails(context.Background(), page, pageSize, orderBy, sortBy, reference, leftEye, rightEye)
	if err != nil {
		httperror.ErrNotFound.WriteError(w)
		return 0, nil, err
	}

	return page, s, nil
}

func (h *Handler) renderSettingsShippingTable(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
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
		{Title: "Collaborator Name", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Collaborator Email", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Customer Name", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Customer Card ID", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Customer Email", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Glasses Reference", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Right Eye Strength", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Left Eye Strength", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Created at:", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Updated at:", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
	}

	page, s, _ := h.getSettingsShipping(w, r)

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
	data := models.SettingsShippingDetailsTable{
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
	t := shipping.ShippingDetailsExpanded(data, models.SettingsShippingDetails{})

	return t, nil
}

func (h *Handler) SettingsShippingPage(w http.ResponseWriter, r *http.Request) error {
	sidebar := h.renderSettingsSidebar()
	renderTable, err := h.renderSettingsShippingTable(w, r)
	if err != nil {
		HandleError(err, "rendering glasses table")
	}
	home := pages.MainLayoutPage("Insert Shipping Form", "Insert Shipping Form", sidebar, renderTable)
	return h.CreateLayout(w, r, "Insert Shipping Form", home).Render(context.Background(), w)
}

func (h *Handler) DeleteCustomer(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["card_id_number"]

	// Delete the glasses
	err := h.service.DeleteCustomer(context.Background(), id)
	if err != nil {
		http.Error(w, "Failed to delete customer", http.StatusInternalServerError)
		return err
	}

	// Return a success response
	w.Header().Set("HX-Redirect", "/settings/shipping")

	return nil
}
