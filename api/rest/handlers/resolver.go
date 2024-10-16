package handlers

import (
	"context"
	"net/http"

	"vk-test-task/internal/service/auth"
	"vk-test-task/internal/service/filmoteka"
	"vk-test-task/pkg/logger"
)

const (
	APIPrefix       = "/api"
	APIVersion      = "v1"
	pathPrefix      = APIPrefix + "/" + APIVersion
	filmotekaPrefix = pathPrefix + "/filmoteka"
)

type Resolver struct {
	serverHost       string
	server           *http.Server
	filmotekaService filmoteka.Service
	authService      auth.Service
}

func NewResolver(serverHost string, filmotekaService filmoteka.Service, authService auth.Service) *Resolver {
	resolver := &Resolver{
		serverHost:       serverHost,
		filmotekaService: filmotekaService,
		authService:      authService,
	}

	mux := http.NewServeMux()
	mux.HandleFunc(pathPrefix+"/auth/login", resolver.login)
	mux.HandleFunc(pathPrefix+"/auth/signup", resolver.signup)

	mux.HandleFunc(filmotekaPrefix+"/stars", resolver.jwtMiddleware(resolver.handleStars))
	mux.HandleFunc(filmotekaPrefix+"/star/", resolver.jwtMiddleware(resolver.handleStar))
	mux.HandleFunc(filmotekaPrefix+"/movies", resolver.jwtMiddleware(resolver.handleMovies))
	mux.HandleFunc(filmotekaPrefix+"/movie/", resolver.jwtMiddleware(resolver.handleMovie))

	loggedRouter := loggingMiddleware(mux)

	server := &http.Server{
		Addr:    serverHost,
		Handler: loggedRouter,
	}

	resolver.server = server

	return resolver
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Log.Info("New request", "method", r.Method, "url", r.URL, "user-agent", r.UserAgent())
		next.ServeHTTP(w, r)
	})
}

func (r *Resolver) jwtMiddleware(nextFunc func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		userData, ok := r.authService.Verify(w, req)
		if !ok {
			return
		}

		ctx := context.WithValue(req.Context(), "user_role", userData.Role) //nolint
		nextFunc(w, req.WithContext(ctx))
	}
}

func (r *Resolver) Run() error {
	return r.server.ListenAndServe()
}

func (r *Resolver) GetAddr() string {
	return r.serverHost
}
