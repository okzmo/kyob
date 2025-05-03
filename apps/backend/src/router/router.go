package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/okzmo/kyob/src/handlers"
	mid "github.com/okzmo/kyob/src/middleware"
)

func Setup() {
	handlers.SetupValidation()
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/v1", func(r chi.Router) {
		r.Post("/signin", handlers.SignIn)
		r.Post("/signup", handlers.SignUp)
		r.Route("/authenticated", func(r chi.Router) {
			r.Use(mid.Auth)
			r.Post("/logout", handlers.Logout)
			r.Post("/server", handlers.CreateServer)
			r.Patch("/servers/{id}", handlers.EditServer)
			r.Delete("/servers/{id}", handlers.DeleteServer)
			r.Post("/channels/{server_id}", handlers.CreateChannel)
			r.Get("/channels/{server_id}", handlers.GetChannels)
			r.Patch("/channels/{id}", handlers.EditChannel)
			r.Delete("/channels/{id}", handlers.DeleteChannel)
			r.Post("/messages/{channel_id}", handlers.CreateMessage)
			r.Get("/messages/{channel_id}", handlers.GetMessages)
			r.Patch("/messages/{id}", handlers.EditMessage)
			r.Delete("/messages/{id}", handlers.DeleteMessage)
		})
	})

	err := http.ListenAndServe(":3000", r)
	if err != nil {
		panic(err)
	}
}
