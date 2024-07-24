package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"context"

	"github.com/FACorreiaa/glasses-management-platform/app/models"
	"github.com/FACorreiaa/glasses-management-platform/app/repository"
	"github.com/FACorreiaa/glasses-management-platform/app/view/user"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
)

// form errors
func (h *Handler) formErrors(err error) []string {
	var decodeErrors form.DecodeErrors
	isDecodeError := errors.As(err, &decodeErrors)
	if isDecodeError {
		var errs []string
		for _, decodeError := range decodeErrors {
			errs = append(errs, decodeError.Error())
		}

		return errs
	}

	// validateErrors, isValidateError := err.(validator.ValidationErrors)

	var validateErrors validator.ValidationErrors
	isValidateError := errors.As(err, &validateErrors)
	if isValidateError {
		var errs []string
		for _, validateError := range validateErrors {
			errs = append(errs, validateError.Translate(h.translator))
		}
		return errs
	}

	return []string{err.Error()}
}

// login

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) error {
	login := user.LoginPage(models.LoginPage{})
	return h.CreateLayout(w, r, "Login", login).Render(context.Background(), w)
}

func (h *Handler) LoginPost(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	var f models.LoginForm
	var token *repository.Token

	err := h.formDecoder.Decode(&f, r.PostForm)
	if err == nil {
		token, err = h.service.Login(r.Context(), f)
	}

	if err != nil {
		login := user.LoginPage(models.LoginPage{Errors: h.formErrors(err)})

		return h.CreateLayout(w, r, "Sign In", login).Render(context.Background(), w)
	}

	s, _ := h.sessions.Get(r, "auth")
	s.Values["token"] = token

	if err := s.Save(r, w); err != nil {
		return errors.New("error saving session")
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) error {
	s, _ := h.sessions.Get(r, "auth")
	token := s.Values["token"]

	if token, ok := token.(string); ok {
		_ = h.service.Logout(r.Context(), token)
	}

	s.Values["token"] = ""
	delete(s.Values, "token")
	delete(s.Values, "user")
	s.Options.MaxAge = -1
	if err := s.Save(r, w); err != nil {
		slog.Error("failed to clear auth session", "err", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}
