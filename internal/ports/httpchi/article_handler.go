package httpchi

import (
	"context"
	"documentation-mini-app/internal/domain/article"
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

type ArticleUsecase interface {
	GetArticleByID(ctx context.Context, id int) (*article.Article, error)
	CreateArticle(ctx context.Context, art *article.Article) error
	AddArticleToDoc(ctx context.Context, artID int, docID int) error
	UpdateArticle(ctx context.Context, art *article.Article) error
	DeleteArticle(ctx context.Context, artID int) error
}

type ArticleHandler struct {
	uc ArticleUsecase

	getAV    *htmlview.TemplateView
	createAV *htmlview.TemplateView
	editAV   *htmlview.TemplateView
	deleteAV *htmlview.TemplateView
}

func NewArticleHandler(uc ArticleUsecase,
	getArticleView *htmlview.TemplateView, createArticleView *htmlview.TemplateView,
	editArticleView *htmlview.TemplateView, deleteArticleView *htmlview.TemplateView,
) *ArticleHandler {
	return &ArticleHandler{uc: uc,
		getAV: getArticleView, createAV: createArticleView,
		editAV: editArticleView, deleteAV: deleteArticleView}
}

func (h *ArticleHandler) SetupRoutes(r chi.Router) {
	r.Route("/articles/{articleID}", func(r chi.Router) {
		r.Get("/", h.GetArticle())

		r.Get("/edit", h.GetEditArticle())
		r.Post("/edit", h.EditArticle())

		r.Get("/delete", h.GetDeleteArticle())
		r.Post("/delete", h.DeleteArticle())
	})

	r.Route("/documentations/{docID}/articles/create", func(r chi.Router) {
		r.Get("/", h.GetCreateArticle())
		r.Post("/", h.CreateArticle())
	})
}

func (h *ArticleHandler) GetArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		artID, err := strconv.Atoi(chi.URLParam(r, "articleID"))
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		art, err := h.uc.GetArticleByID(r.Context(), artID)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = h.getAV.ToWriter(w, art)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h *ArticleHandler) GetCreateArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h.createAV.ToWriter(w, struct{}{})
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h *ArticleHandler) CreateArticle() http.HandlerFunc {
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
		desc := q.Get("description")

		if name == "" {
			http.Error(w, "name can't be empty", http.StatusBadRequest)
			return
		}

		art := article.Article{
			Name:        name,
			Description: desc,
		}

		err = h.uc.CreateArticle(r.Context(), &art)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = h.uc.AddArticleToDoc(r.Context(), art.ID, docID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/articles/%v", art.ID), http.StatusSeeOther)
	}
}

func (h *ArticleHandler) GetEditArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		artID, err := strconv.Atoi(chi.URLParam(r, "articleID"))
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		art, err := h.uc.GetArticleByID(r.Context(), artID)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = h.editAV.ToWriter(w, art)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h *ArticleHandler) EditArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		artID, err := strconv.Atoi(chi.URLParam(r, "articleID"))
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

		if name == "" {
			http.Error(w, "name can't be empty", http.StatusBadRequest)
			return
		}

		art := article.Article{
			ID:          artID,
			Name:        name,
			Description: desc,
		}

		err = h.uc.UpdateArticle(r.Context(), &art)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/articles/%v", art.ID), http.StatusSeeOther)
	}
}

func (h *ArticleHandler) GetDeleteArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		artID, err := strconv.Atoi(chi.URLParam(r, "articleID"))
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		art, err := h.uc.GetArticleByID(r.Context(), artID)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = h.deleteAV.ToWriter(w, art)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h *ArticleHandler) DeleteArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		artID, err := strconv.Atoi(chi.URLParam(r, "articleID"))
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = h.uc.DeleteArticle(r.Context(), artID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
