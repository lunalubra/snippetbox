package main

import (
	"database/sql"
	"errors"
	"flag"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lunalubra/snippetbox/internal/models"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type application struct {
	logger         *slog.Logger
	snippets       models.SnippetModelInterface
	users          models.UserModelInterface
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	defaultAddr := ":4000"
	if port := os.Getenv("PORT"); port != "" {
		defaultAddr = ":" + port
	}

	defaultDSN := "web:acosta@/snippetbox?parseTime=true"
	if envDSN := os.Getenv("DSN"); envDSN != "" {
		defaultDSN = envDSN
	}

	// Upsun's router terminates TLS, so the app must serve plain HTTP there.
	defaultTLS := os.Getenv("PLATFORM_APPLICATION") == ""

	addr := flag.String("addr", defaultAddr, "HTTP network address")
	dsn := flag.String("dsn", defaultDSN, "MySQL data source name")
	useTLS := flag.Bool("tls", defaultTLS, "Serve HTTPS with the local self-signed certificate")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	application := application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	router := application.routes()

	logger.Info("starting server", "addr", *addr)
	s := http.Server{
		Addr:         *addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      router,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
	if *useTLS {
		err = s.ListenAndServeTLS("./assets/tls/localhost.pem", "./assets/tls/localhost-key.pem")
	} else {
		err = s.ListenAndServe()
	}
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
