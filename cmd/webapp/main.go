package main

import (
	"context"
	"documentation-mini-app/internal/adapters/pgstore"
	"documentation-mini-app/internal/config"
	"documentation-mini-app/internal/ports/httpchi"
	"documentation-mini-app/internal/usecase/appuc"
	"documentation-mini-app/internal/usecase/articleuc"
	"documentation-mini-app/internal/usecase/docuc"
	"documentation-mini-app/internal/usecase/exampleuc"
	"documentation-mini-app/internal/views/htmlview"
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5"

	"github.com/go-chi/chi/v5"
)

func parseConfig(path string) *config.Config {
	configFile, err := os.Open(path)
	if err != nil {
		log.Panicf("config file open: %v\n", err)
	}

	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			log.Panicf("file close: %v\n", err)
		}
	}(configFile)

	return config.Parse(configFile)
}

func main() {
	var configPath string
	flag.StringVar(&configPath, "config-path", "configs/server_config.json", "path to config file")
	flag.Parse()

	dbURL := os.Getenv("DOC_DATABASE_URL")
	if dbURL == "" {
		log.Fatalln("Need DOC_DATABASE_URL env variable.")
	}

	conf := parseConfig(configPath)
	fmt.Println(conf.Addr)

	ctx := context.TODO()

	store, err := pgstore.New(ctx, dbURL)
	if err != nil {
		log.Fatalln(err)
	}
	defer store.Close()

	contentView, err := htmlview.New("templates/contents.html")
	if err != nil {
		log.Panicf("contentView create: %v\n", err)
	}

	crossedView, err := htmlview.New("templates/crossed.html")
	if err != nil {
		log.Panicf("contentView create: %v\n", err)
	}

	getDocView, err := htmlview.New("templates/docs/get_doc.html")
	if err != nil {
		log.Panicf("contentView create: %v\n", err)
	}

	createDocView, err := htmlview.New("templates/docs/create_doc.html")
	if err != nil {
		log.Panicf("contentView create: %v\n", err)
	}

	editDocView, err := htmlview.New("templates/docs/edit_doc.html")
	if err != nil {
		log.Panicf("contentView create: %v\n", err)
	}

	deleteDocView, err := htmlview.New("templates/docs/delete_doc.html")
	if err != nil {
		log.Panicf("contentView create: %v\n", err)
	}

	getArticleView, err := htmlview.New("templates/articles/get_article.html")
	if err != nil {
		log.Panicf("contentView create: %v\n", err)
	}

	createArticleView, err := htmlview.New("templates/articles/create_article.html")
	if err != nil {
		log.Panicf("contentView create: %v\n", err)
	}

	editArticleView, err := htmlview.New("templates/articles/edit_article.html")
	if err != nil {
		log.Panicf("contentView create: %v\n", err)
	}

	deleteArticleView, err := htmlview.New("templates/articles/delete_article.html")
	if err != nil {
		log.Panicf("contentView create: %v\n", err)
	}

	getExampleView, err := htmlview.New("templates/examples/get_example.html")
	if err != nil {
		log.Panicf("contentView create: %v\n", err)
	}

	createExampleView, err := htmlview.New("templates/examples/create_example.html")
	if err != nil {
		log.Panicf("contentView create: %v\n", err)
	}

	editExampleView, err := htmlview.New("templates/examples/edit_example.html")
	if err != nil {
		log.Panicf("contentView create: %v\n", err)
	}

	deleteExampleView, err := htmlview.New("templates/examples/delete_example.html")
	if err != nil {
		log.Panicf("contentView create: %v\n", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	docUC := docuc.New(store)
	docHandler := httpchi.NewDocHandler(docUC,
		getDocView, createDocView, editDocView, deleteDocView)

	artUC := articleuc.New(store)
	artHandler := httpchi.NewArticleHandler(artUC,
		getArticleView, createArticleView, editArticleView, deleteArticleView)

	exaUC := exampleuc.New(store)
	exaHandler := httpchi.NewExampleHandler(exaUC,
		getExampleView, createExampleView, editExampleView, deleteExampleView)

	appUC := appuc.New(store)
	appHandler := httpchi.NewAppHandler(r, appUC,
		artHandler, docHandler, exaHandler,
		contentView, crossedView)

	server := http.Server{
		Addr:         conf.Addr,
		Handler:      appHandler,
		ReadTimeout:  time.Duration(10) * time.Second,
		WriteTimeout: time.Duration(10) * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Panicln(err)
	}
}
