package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/anthdm/hollywood/actor"
	"github.com/go-chi/chi/v5"
	"github.com/okzmo/kyob/db"
	"github.com/okzmo/kyob/internal/api/actors"
	services "github.com/okzmo/kyob/internal/service"
	"github.com/okzmo/kyob/internal/utils"
	proto "github.com/okzmo/kyob/types"
)

func CreateServer(w http.ResponseWriter, r *http.Request) {
	var body services.CreateServerBody

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		slog.Error(err.Error())
		utils.RespondWithError(w, http.StatusBadRequest, "Failed to parse given image.")
		return
	}

	file, fileHeader, err := r.FormFile("avatar")
	if err != nil {
		slog.Error(err.Error())
		utils.RespondWithError(w, http.StatusBadRequest, "Failed to get image.")
		return
	}
	defer file.Close()

	body.Name = r.FormValue("name")
	body.Description = r.FormValue("description")
	body.Private = r.FormValue("private") == "true"
	x, errX := strconv.Atoi(r.FormValue("x"))
	y, errY := strconv.Atoi(r.FormValue("x"))
	if errX != nil || errY != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid coordinates.")
	}

	body.X = x
	body.Y = y

	cropJSON := r.FormValue("crop")
	var crop services.Crop
	if err := json.Unmarshal([]byte(cropJSON), &crop); err != nil {
		slog.Error(err.Error())
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid  crop data.")
		return
	}
	body.Crop = crop

	err = validate.Struct(body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	server, err := services.CreateServer(r.Context(), file, fileHeader, &body)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user := r.Context().Value("user").(db.User)
	serverPID := actors.ServersEngine.Spawn(actors.NewServer, "server", actor.WithID(strconv.Itoa(int(server.ID))))
	userPID := actors.UsersEngine.Registry.GetPID("user", strconv.Itoa(int(user.ID)))
	actors.UsersEngine.Send(userPID, &proto.NewServerCreated{
		ServerPID: serverPID.ID,
		Address:   serverPID.Address,
	})

	utils.RespondWithJSON(w, http.StatusCreated, server)
}

func EditServer(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var body services.EditServerBody
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = validate.Struct(body)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = services.EditServer(r.Context(), id, &body)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUnauthorizedServerEdition):
			utils.RespondWithError(w, http.StatusUnauthorized, "You cannot edit this server.")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusContinue, &DefaultResponse{Message: "success"})
}

func DeleteServer(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = services.DeleteServer(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUnauthorizedServerDeletion):
			utils.RespondWithError(w, http.StatusUnauthorized, "You cannot delete this server.")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusContinue, &DefaultResponse{Message: "success"})
}
