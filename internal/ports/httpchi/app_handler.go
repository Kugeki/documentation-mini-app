package httpchi

import (
	"context"
	"documentation-mini-app/internal/domain/article"
	"documentation-mini-app/internal/domain/crossed"
	"documentation-mini-app/internal/domain/doc"
	"documentation-mini-app/internal/views/htmlview"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type AppUsecase interface {
	GetDocByID(ctx context.Context, id int) (*doc.Documentation, error)
	GetAllDoc(ctx context.Context) ([]*doc.Documentation, error)
	GetCrossed(ctx context.Context) (*crossed.Crossed, error)
	GetArticlesWithoutDoc(ctx context.Context) ([]article.Article, error)
}

type AppHandler struct {
	router chi.Router
	uc     AppUsecase

	contentView *htmlview.TemplateView
	crossedView *htmlview.TemplateView

	artHandler     *ArticleHandler
	docHandler     *DocHandler
	exampleHandler *ExampleHandler
}

func NewAppHandler(r chi.Router, uc AppUsecase, ah *ArticleHandler, dh *DocHandler, eh *ExampleHandler,
	contentsView *htmlview.TemplateView, crossedView *htmlview.TemplateView,
) *AppHandler {
	h := &AppHandler{
		router:         r,
		uc:             uc,
		contentView:    contentsView,
		crossedView:    crossedView,
		artHandler:     ah,
		docHandler:     dh,
		exampleHandler: eh,
	}
	h.SetupRoutes()

	return h
}

func (h *AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *AppHandler) SetupRoutes() {
	h.router.Get("/", h.GetContents())
	h.router.Get("/crossed", h.GetCrossed())

	h.router.Route("/", h.setupOtherRoutes)
}

func (h *AppHandler) setupOtherRoutes(r chi.Router) {
	h.artHandler.SetupRoutes(r)
	h.docHandler.SetupRoutes(r)
	h.exampleHandler.SetupRoutes(r)
}

func (h *AppHandler) GetContents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		docs, err := h.uc.GetAllDoc(r.Context())
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		artsWithoutDoc, err := h.uc.GetArticlesWithoutDoc(r.Context())
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		docs = append(docs, &doc.Documentation{
			Name:     "Статьи без документации",
			Articles: artsWithoutDoc,
		})

		w.WriteHeader(http.StatusOK)

		err = h.contentView.ToWriter(w, docs)
		if err != nil {
			log.Println(err)
		}
	}
}

func (h *AppHandler) GetCrossed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		crsd, err := h.uc.GetCrossed(r.Context())
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = h.crossedView.ToWriter(w, crsd)
		if err != nil {
			log.Println(err)
		}
	}
}
