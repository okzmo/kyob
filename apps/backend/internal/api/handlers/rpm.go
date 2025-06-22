package handlers

import (
	"fmt"
	"net/http"

	queries "github.com/okzmo/kyob/db/gen_queries"
	services "github.com/okzmo/kyob/internal/service"
	"github.com/okzmo/kyob/internal/utils"
)

func GetRPMAssets(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(queries.User)

	rpmService := services.NewRPMService()
	assets, err := rpmService.GetAssets(user.RpmToken.String, user.RpmID.String)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to get assets")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, assets)
}

func UpdateRPMAvatar(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AssetID   string `json:"asset_id"`
		AssetType string `json:"asset_type"`
	}

	if err := utils.ParseAndValidate(r, validate, &req); err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	user := r.Context().Value("user").(queries.User)
	rpmService := services.NewRPMService()

	err := rpmService.EquipAsset(r.Context(), user.RpmToken.String, "685470d6259f6d6e519467ee", req.AssetID, req.AssetType)
	fmt.Println(err)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, "failed to equip asset")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, DefaultResponse{Message: "ok"})
}
