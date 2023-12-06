package httpchi

import (
	"context"
	"documentation-mini-app/internal/domain/example"
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

type ExampleUsecase interface {
	GetExampleByID(ctx context.Context, id int) (*example.Example, error)
	CreateExample(ctx context.Context, exa *example.Example) error
	AddExampleToArticle(ctx context.Context, exaID int, artID int) error
	UpdateExample(ctx context.Context, exa *example.Example) error
	DeleteExample(ctx context.Context, id int) error
}

type ExampleHandler struct {
	uc ExampleUsecase

	getEV    *htmlview.TemplateView
	createEV *htmlview.TemplateView
	editEV   *htmlview.TemplateView
	deleteEV *htmlview.TemplateView
}

func NewExampleHandler(uc ExampleUsecase,
	getExampleView *htmlview.TemplateView, createExampleView *htmlview.TemplateView,
	editExampleView *htmlview.TemplateView, deleteExampleView *htmlview.TemplateView,
) *ExampleHandler {
	return &ExampleHandler{uc: uc,
		getEV: getExampleView, createEV: createExampleView,
		editEV: editExampleView, deleteEV: deleteExampleView,
	}
}

func (h *ExampleHandler) SetupRoutes(r chi.Router) {
	r.Route("/examples", func(r chi.Router) {
		r.Route("/{exaID}", func(r chi.Router) {
			r.Get("/", h.GetExample())

			r.Get("/edit", h.GetEditExample())
			r.Post("/edit", h.EditExample())

			r.Get("/delete", h.GetDeleteExample())
			r.Post("/delete", h.DeleteExample())
		})
	})

	r.Route("/articles/{artID}/examples/create", func(r chi.Router) {
		r.Get("/", h.GetCreateExample())
		r.Post("/", h.CreateExample())
	})
}

func (h *ExampleHandler) GetExample() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		exaID, err := strconv.Atoi(chi.URLParam(r, "exaID"))
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		d, err := h.uc.GetExampleByID(r.Context(), exaID)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = h.getEV.ToWriter(w, d)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h *ExampleHandler) GetCreateExample() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h.createEV.ToWriter(w, struct{}{})
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h *ExampleHandler) CreateExample() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		artID, err := strconv.Atoi(chi.URLParam(r, "artID"))
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
		desc := q.Get("description")
		code := q.Get("code")
		outp := q.Get("output")

		if name == "" {
			http.Error(w, "name can't be empty", http.StatusBadRequest)
			return
		}

		exa := example.Example{
			Name:        name,
			Description: desc,
			Code:        code,
			Output:      outp,
			Priority:    0,
		}

		err = h.uc.CreateExample(r.Context(), &exa)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = h.uc.AddExampleToArticle(r.Context(), exa.ID, artID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/examples/%v", exa.ID), http.StatusSeeOther)
	}
}

func (h *ExampleHandler) GetEditExample() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		exaID, err := strconv.Atoi(chi.URLParam(r, "exaID"))
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		exa, err := h.uc.GetExampleByID(r.Context(), exaID)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = h.editEV.ToWriter(w, exa)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h *ExampleHandler) EditExample() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		exaID, err := strconv.Atoi(chi.URLParam(r, "exaID"))
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
		desc := q.Get("description")
		code := q.Get("code")
		outp := q.Get("output")

		if name == "" {
			http.Error(w, "name can't be empty", http.StatusBadRequest)
			return
		}

		exa := example.Example{
			ID:          exaID,
			Name:        name,
			Description: desc,
			Code:        code,
			Output:      outp,
			Priority:    0,
		}

		err = h.uc.UpdateExample(r.Context(), &exa)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/examples/%v", exa.ID), http.StatusSeeOther)
	}
}

func (h *ExampleHandler) GetDeleteExample() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		exaID, err := strconv.Atoi(chi.URLParam(r, "exaID"))
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		exa, err := h.uc.GetExampleByID(r.Context(), exaID)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = h.deleteEV.ToWriter(w, exa)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h *ExampleHandler) DeleteExample() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		exaID, err := strconv.Atoi(chi.URLParam(r, "exaID"))
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = h.uc.DeleteExample(r.Context(), exaID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
