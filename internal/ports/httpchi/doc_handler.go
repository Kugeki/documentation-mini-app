package httpchi

import (
	"context"
	"documentation-mini-app/internal/domain/doc"
	"documentation-mini-app/internal/views/htmlview"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type DocUsecase interface {
	GetDocByID(ctx context.Context, docID int) (*doc.Documentation, error)
	CreateDoc(ctx context.Context, d *doc.Documentation) error
	UpdateDoc(ctx context.Context, d *doc.Documentation) error
	DeleteDoc(ctx context.Context, docID int) error
}

type DocHandler struct {
	uc DocUsecase

	getDV    *htmlview.TemplateView
	createDV *htmlview.TemplateView
	editDV   *htmlview.TemplateView
	deleteDV *htmlview.TemplateView
}

func NewDocHandler(uc DocUsecase,
	getDocView *htmlview.TemplateView, createDocView *htmlview.TemplateView,
	editDocView *htmlview.TemplateView, deleteDocView *htmlview.TemplateView,
) *DocHandler {
	return &DocHandler{uc: uc,
		createDV: createDocView, getDV: getDocView,
		editDV: editDocView, deleteDV: deleteDocView}
}

func (h *DocHandler) SetupRoutes(r chi.Router) {
	r.Route("/documentations", func(r chi.Router) {
		r.Get("/create", h.GetCreateDoc())
		r.Post("/create", h.CreateDoc())

		r.Route("/{docID}", func(r chi.Router) {
			r.Get("/", h.GetDoc())

			r.Get("/edit", h.GetEditDoc())
			r.Post("/edit", h.EditDoc())

			r.Get("/delete", h.GetDeleteDoc())
			r.Post("/delete", h.DeleteDoc())
		})
	})
}

func (h *DocHandler) GetDoc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		docID, err := strconv.Atoi(chi.URLParam(r, "docID"))
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		d, err := h.uc.GetDocByID(r.Context(), docID)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = h.getDV.ToWriter(w, d)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h *DocHandler) GetCreateDoc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h.createDV.ToWriter(w, struct{}{})
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h *DocHandler) CreateDoc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b strings.Builder
		_, err := io.Copy(&b, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		q, err := url.ParseQuery(b.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		name := q.Get("name")

		if name == "" {
			http.Error(w, "name can't be empty", http.StatusBadRequest)
			return
		}

		d := doc.Documentation{
			Name: name,
		}

		err = h.uc.CreateDoc(r.Context(), &d)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/documentations/%v", d.ID), http.StatusSeeOther)
	}
}

func (h *DocHandler) GetEditDoc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		docID, err := strconv.Atoi(chi.URLParam(r, "docID"))
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		d, err := h.uc.GetDocByID(r.Context(), docID)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = h.editDV.ToWriter(w, d)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h *DocHandler) EditDoc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		docID, err := strconv.Atoi(chi.URLParam(r, "docID"))
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var b strings.Builder
		_, err = io.Copy(&b, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		q, err := url.ParseQuery(b.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		name := q.Get("name")

		if name == "" {
			http.Error(w, "name can't be empty", http.StatusBadRequest)
			return
		}

		d := doc.Documentation{
			ID:   docID,
			Name: name,
		}

		err = h.uc.UpdateDoc(r.Context(), &d)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/documentations/%v", d.ID), http.StatusSeeOther)
	}
}

func (h *DocHandler) GetDeleteDoc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		docID, err := strconv.Atoi(chi.URLParam(r, "docID"))
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		d, err := h.uc.GetDocByID(r.Context(), docID)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = h.deleteDV.ToWriter(w, d)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h *DocHandler) DeleteDoc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		docID, err := strconv.Atoi(chi.URLParam(r, "docID"))
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = h.uc.DeleteDoc(r.Context(), docID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
