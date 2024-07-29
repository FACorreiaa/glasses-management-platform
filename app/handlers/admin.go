package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	httperror "github.com/FACorreiaa/glasses-management-platform/app/errors"
	"github.com/FACorreiaa/glasses-management-platform/app/models"
	"github.com/FACorreiaa/glasses-management-platform/app/static/svg"
	"github.com/FACorreiaa/glasses-management-platform/app/view/admin"
	"github.com/FACorreiaa/glasses-management-platform/app/view/pages"
	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	SubmitAction          = "submit"
	MinPasswordLength     = "Password must be at least 10 characters long"
	PasswordDoNotMatch    = "Passwords do not match"
	UsernameMinCharLength = "Username must be at least 3 characters long"
)

func (h *Handler) getCollaborators(w http.ResponseWriter, r *http.Request) (int, []models.UserSession, error) {
	pageSize := 10
	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	email := r.FormValue("email")

	u, err := h.service.GetUsers(context.Background(), page, pageSize, orderBy, sortBy, email)
	if err != nil {
		httperror.ErrNotFound.WriteError(w)
		return 0, nil, err
	}

	return page, u, nil
}

func (h *Handler) renderCollaboratorsTable(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	var page int
	var sortAux string
	orderBy := r.FormValue("orderBy")
	sortBy := r.FormValue("sortBy")
	// brand := r.FormValue("brand")

	if sortBy == ASC {
		sortAux = DESC
	} else {
		sortAux = ASC
	}

	columnNames := []models.ColumnItems{
		{Title: "Username", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Email", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Role", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Created At", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
		{Title: "Updated At", Icon: svg.ArrowOrderIcon(), SortParam: sortAux},
	}

	page, u, _ := h.getCollaborators(w, r)

	if len(u) == 0 {
		message := admin.UserEmptyPage()
		return message, nil
	}

	nextPage := page + 1
	prevPage := page - 1
	if prevPage <= 1 {
		prevPage = 1
	}

	lastPage, err := h.service.GetUsersSum()
	if err != nil {
		HandleError(err, " fetching tax")
		return nil, err
	}
	data := models.UsersTable{
		Column:     columnNames,
		Users:      u,
		PrevPage:   prevPage,
		NextPage:   nextPage,
		Page:       page,
		LastPage:   lastPage,
		OrderParam: orderBy,
		SortParam:  sortAux,
	}

	t := admin.UsersTable(data, models.RegisterFormValues{})

	return t, nil
}

// UsersPage users page for admin to manage views TODO
func (h *Handler) UsersPage(w http.ResponseWriter, r *http.Request) error {
	table, err := h.renderCollaboratorsTable(w, r)
	if err != nil {
		HandleError(err, "rendering glasses table")
	}
	sidebar := h.renderSettingsSidebar()
	users := pages.MainLayoutPage("List of collaborators", "List of collaborators", sidebar, table)
	data := h.CreateLayout(w, r, "Users", users).Render(context.Background(), w)
	return data
}

func (h *Handler) UserInsertPage(w http.ResponseWriter, r *http.Request) error {
	register := admin.RegisterPage(models.RegisterFormValues{})
	sidebar := h.renderSettingsSidebar()
	u := pages.MainLayoutPage("List of collaborators", "List of collaborators", sidebar, register)
	return h.CreateLayout(w, r, "Insert new collaborator", u).Render(context.Background(), w)
}

func (h *Handler) UserInsertModalForm(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		HandleError(err, "parsing form")
		return err
	}

	f := models.RegisterFormValues{
		Username:        r.FormValue("username"),
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		PasswordConfirm: r.FormValue("password_confirm"),
		FieldErrors:     make(map[string]string),
	}

	if len(f.Password) < 10 {
		f.FieldErrors["password"] = MinPasswordLength
	}

	if f.Password != f.PasswordConfirm {
		f.FieldErrors["password_confirm"] = PasswordDoNotMatch
	}

	if len(f.Username) < 3 {
		f.FieldErrors["username"] = UsernameMinCharLength
	}

	if len(f.FieldErrors) > 0 {
		// HTMX request, return the form only
		form := admin.UserInsertModal(f).Render(context.Background(), w)
		return form

	}

	if _, err := h.service.InsertUser(context.Background(), f); err != nil {
		return fmt.Errorf("error inserting users: %v", err)
	}

	actionType := r.FormValue("action")

	if actionType == "back" {
		w.Header().Set("HX-Redirect", "/settings/collaborators")
	} else if actionType == SubmitAction {
		w.Header().Set("HX-Redirect", "/settings/collaborators")
	}

	return nil
}

// UserRegisterPost register new user
func (h *Handler) UserRegisterPost(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		HandleError(err, "parsing form")
		return err
	}

	f := models.RegisterFormValues{
		Username:        r.FormValue("username"),
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		PasswordConfirm: r.FormValue("password_confirm"),
		FieldErrors:     make(map[string]string),
	}

	if len(f.Password) < 10 {
		f.FieldErrors["password"] = MinPasswordLength
	}

	if f.Password != f.PasswordConfirm {
		f.FieldErrors["password_confirm"] = PasswordDoNotMatch
	}

	if len(f.Username) < 3 {
		f.FieldErrors["username"] = UsernameMinCharLength
	}

	if len(f.FieldErrors) > 0 {
		sidebar := h.renderSettingsSidebar()
		form := admin.RegisterPage(f)
		register := pages.MainLayoutPage("Insert user Form", "Insert user Form", sidebar, form)
		return h.CreateLayout(w, r, "Register collaborator", register).Render(context.Background(), w)

	}

	if _, err := h.service.InsertUser(context.Background(), f); err != nil {
		return fmt.Errorf("error inserting users: %v", err)
	}

	actionType := r.FormValue("action")

	if actionType == "back" {
		w.Header().Set("HX-Redirect", "/settings/collaborators")
	} else if actionType == SubmitAction {
		w.Header().Set("HX-Redirect", "/settings/collaborators")
	}

	return nil
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	userIDStr := vars["user_id"]
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return err
	}

	// Delete the glasses
	if err = h.service.DeleteUser(context.Background(), userID); err != nil {
		http.Error(w, "Failed to delete glasses", http.StatusInternalServerError)
		return err
	}

	w.Header().Set("HX-Redirect", "/settings/collaborators")

	return nil
}

func (h *Handler) UpdateUserPage(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	userID := vars["user_id"]
	u, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return err
	}

	g, err := h.service.GetUsersByID(context.Background(), u)
	if err != nil {
		http.Error(w, "Failed to retrieve glasses", http.StatusInternalServerError)
		return err
	}

	form := models.UpdateUserForm{
		Values: map[string]string{
			"Username": g.Username,
			"Email":    g.Email,
			"Role":     g.Role,
		},
	}

	f := admin.UserUpdateForm(form, userID)
	sidebar := h.renderSettingsSidebar()

	updatePage := pages.MainLayoutPage("Update users", "form to update users", sidebar, f)
	return h.CreateLayout(w, r, "Update Glasses", updatePage).Render(context.Background(), w)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	userIDStr := vars["user_id"]
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return err
	}

	if err = r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return err
	}

	action := r.FormValue("action")
	if action == "return" {
		w.Header().Set("HX-Redirect", "/settings/collaborators")
		return nil
	}

	g := models.UpdateUserForm{
		UserID:          userID,
		Email:           r.FormValue("email"),
		Username:        r.FormValue("username"),
		Role:            r.FormValue("role"),
		Password:        r.FormValue("password"),
		PasswordConfirm: r.FormValue("password_confirm"),
		FieldErrors:     make(map[string]string),
	}

	if len(g.Password) < 10 {
		g.FieldErrors["password"] = MinPasswordLength
	}

	if g.Password != g.PasswordConfirm {
		g.FieldErrors["password_confirm"] = PasswordDoNotMatch
	}

	if len(g.Username) < 3 {
		g.FieldErrors["username"] = UsernameMinCharLength
	}

	if len(g.FieldErrors) > 0 {
		form := admin.UserUpdateForm(g, userIDStr).Render(context.Background(), w)
		return form
	}

	if err = h.service.UpdateUser(context.Background(), g); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return err
	}

	w.Header().Set("HX-Redirect", "/settings/collaborators")

	return nil
}
