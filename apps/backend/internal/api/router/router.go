package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/okzmo/kyob/internal/api/handlers"
	mid "github.com/okzmo/kyob/internal/api/middleware"
)

func Setup() {
	handlers.SetupValidation()
	handlers.SetupWebsocket()

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
			r.Get("/connect/{user_id}", handlers.WS)
			r.Get("/setup", handlers.Setup)
			r.Post("/server", handlers.CreateServer)
			r.Post("/server/join", handlers.JoinServer)
			r.Post("/server/{id}/leave", handlers.LeaveServer)
			r.Get("/server/create_invite/{id}", handlers.CreateServerInvite)
			r.Patch("/servers/{id}", handlers.EditServer)
			r.Delete("/servers/{id}", handlers.DeleteServer)
			r.Post("/channels/{server_id}", handlers.CreateChannel)
			r.Patch("/channels/{channel_id}", handlers.EditChannel)
			r.Delete("/channels/{server_id}/{channel_id}", handlers.DeleteChannel)
			r.Post("/messages/{server_id}/{channel_id}", handlers.CreateMessage)
			r.Get("/messages/{channel_id}", handlers.GetMessages)
			r.Patch("/messages/{message_id}", handlers.EditMessage)
			r.Delete("/messages/{message_id}", handlers.DeleteMessage)
			r.Post("/logout", handlers.Logout)
		})
	})

	err := http.ListenAndServe(":3000", r)
	if err != nil {
		panic(err)
	}
}
