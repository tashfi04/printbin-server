package api

import (
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tashfi04/printbin-server/api/middlewares"
	"github.com/tashfi04/printbin-server/config"
	"github.com/thedevsaddam/renderer"
	"net/http"
)

var router = chi.NewRouter()
var rndr *renderer.Render

func init() {
	rndr = renderer.New()
}

// Router main router
func Router() http.Handler {

	router.Use(middlewares.Recoverer) // Use recoverer middleware at last to fire the defer at first
	router.Use(chiMiddleware.Logger)
	router.Use(chiMiddleware.Heartbeat("/health"))

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   config.App().CorsAllowedHosts,
		AllowedMethods:   []string{"GET", "PATCH", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	router.Mount("/metrics", promhttp.Handler())

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		rndr.JSON(w, http.StatusNotFound, renderer.M{
			"message": "Route not found!",
		})
	})

	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		rndr.JSON(w, http.StatusNotFound, renderer.M{
			"message": "Method not allowed!",
		})
	})

	RegisterRoutes()

	return router
}

// RegisterRoutes registers provided routes
func RegisterRoutes() {

	router.Route("/api/v1/", func(r chi.Router) {
		r.Mount("/auth", authenticationRoutes())

		r.Group(func(r chi.Router) {
			r.Use(middlewares.AuthenticateClient)
			r.Get("/isAuth", isAuthenticated)
			r.Mount("/file", fileRoutes())
			r.Mount("/print", printFileRoutes())
			r.Mount("/admin", adminRoutes())
		})
	})
}

func authenticationRoutes() http.Handler {

	h := chi.NewRouter()

	h.Group(func(r chi.Router) {
		r.Post("/login", login)
	})

	return h
}

func fileRoutes() http.Handler {
	h := chi.NewRouter()

	h.Group(func(r chi.Router) {
		r.Get("/", listUserFiles)
		r.Post("/", submitFile)
	})

	return h
}

func printFileRoutes() http.Handler {
	h := chi.NewRouter()

	h.Group(func(r chi.Router) {
		r.Use(middlewares.AuthenticateAdmin)
		r.Get("/", listFiles)
		r.Patch("/", updateStatus)
		r.Get("/storage/*", serveFile)
		r.Get("/rooms", listRooms)
	})

	return h
}

func adminRoutes() http.Handler {
	h := chi.NewRouter()

	h.Group(func(r chi.Router) {
		r.Use(middlewares.AuthenticateAdmin)
		r.Post("/", uploadUser)
	})

	return h
}
