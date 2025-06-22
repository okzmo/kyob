package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ReadyPlayerMeService struct {
	BaseURL       string
	ApplicationID string
	Subdomain     string
	HTTPClient    *http.Client
}

type RPMResponse[T any] struct {
	Data T `json:"data"`
}

type RPMUser struct {
	ID          string `json:"id"`
	AccessToken string `json:"token"`
}

type RPMTemplate struct {
	ImageURL string `json:"imageUrl"`
	Gender   string `json:"gender"`
	ID       string `json:"id"`
}

type RPMAssignAvatar struct {
	ID       string `json:"id"`
	Partner  string `json:"partner"`
	Gender   string `json:"gender"`
	BodyType string `json:"bodyType"`
	Assets   any    `json:"assets"`
	Favorite bool   `json:"favorite"`
	Default  bool   `json:"default"`
}

func NewRPMService() *ReadyPlayerMeService {
	return &ReadyPlayerMeService{
		BaseURL:       os.Getenv("RPM_BASEURL"),
		ApplicationID: os.Getenv("RPM_APPID"),
		Subdomain:     os.Getenv("RPM_SUBDOMAIN"),
		HTTPClient:    &http.Client{},
	}
}

func (rpm *ReadyPlayerMeService) CreateDefaultBody(ctx context.Context) (string, *RPMUser, error) {
	rpmUser, err := rpm.CreateAnonymousUser(ctx)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create anonymous user: %w", err)
	}

	templates, err := rpm.GetTemplates(ctx, rpmUser.AccessToken)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get templates on RPM: %w", err)
	}

	template := templates[0]

	avatarID, err := rpm.AssignAvatar(ctx, rpmUser.AccessToken, template.ID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to assign avatar: %w", err)
	}

	err = rpm.SaveAvatar(ctx, rpmUser.AccessToken, avatarID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to save avatar: %w", err)
	}

	return avatarID, rpmUser, nil
}

func (rpm *ReadyPlayerMeService) CreateAnonymousUser(ctx context.Context) (*RPMUser, error) {
	response, err := makeRPMRequest[RPMUser](ctx, fmt.Sprintf("https://%s.readyplayer.me", rpm.Subdomain), rpm.HTTPClient, http.MethodPost, "/api/users", "", nil, "")
	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func (rpm *ReadyPlayerMeService) GetTemplates(ctx context.Context, token string) ([]RPMTemplate, error) {
	response, err := makeRPMRequest[[]RPMTemplate](ctx, rpm.BaseURL, rpm.HTTPClient, http.MethodGet, "/v2/avatars/templates", token, nil, "")
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (rpm *ReadyPlayerMeService) AssignAvatar(ctx context.Context, token string, templateID string) (string, error) {
	payload := map[string]any{
		"data": map[string]string{
			"partner":  rpm.Subdomain,
			"bodyType": "fullbody",
		},
	}

	assignment, err := makeRPMRequest[RPMAssignAvatar](ctx, rpm.BaseURL, rpm.HTTPClient, http.MethodPost, "/v2/avatars/templates/"+templateID, token, payload, "")
	if err != nil {
		return "", err
	}

	return assignment.Data.ID, nil
}

func (rpm *ReadyPlayerMeService) SaveAvatar(ctx context.Context, token, avatarID string) error {
	_, err := makeRPMRequest[any](ctx, rpm.BaseURL, rpm.HTTPClient, http.MethodPut, "/v2/avatars/"+avatarID, token, nil, "")
	if err != nil {
		return err
	}

	return nil
}

func (rpm *ReadyPlayerMeService) GetAssets(token, userID string) (any, error) {
	path := fmt.Sprintf("/v1/assets?filter=usable-by-user-and-app&filterApplicationId=%s&filterUserId=%s", rpm.ApplicationID, userID)

	response, err := makeRPMRequest[any](context.TODO(), rpm.BaseURL, rpm.HTTPClient, http.MethodGet, path, token, nil, rpm.ApplicationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get rpm assets: %w", err)
	}

	return response.Data, nil
}

func (rpm *ReadyPlayerMeService) EquipAsset(ctx context.Context, token, avatarID, assetID, assetType string) error {
	payload := map[string]any{
		"data": map[string]string{
			"assetId": assetID,
		},
	}

	_, err := makeRPMRequest[any](ctx, rpm.BaseURL, rpm.HTTPClient, http.MethodPut, "/v1/avatars/"+avatarID+"/equip", token, payload, "")
	if err != nil {
		return fmt.Errorf("failed to update avatar asset: %w", err)
	}

	return nil
}

func makeRPMRequest[T any](ctx context.Context, baseURL string, client *http.Client, method, path, token string, payload any, appIDHeader string) (*RPMResponse[T], error) {
	var body io.Reader

	if payload != nil {
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(jsonPayload)
	}

	req, err := http.NewRequestWithContext(ctx, method, baseURL+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	if appIDHeader != "" {
		req.Header.Set("X-APP-ID", appIDHeader)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("RPM API error: %d - %s", res.StatusCode, string(bodyBytes))
	}

	if res.StatusCode == 204 {
		return &RPMResponse[T]{}, nil
	}

	if res.ContentLength == 0 {
		return &RPMResponse[T]{}, nil
	}

	var result RPMResponse[T]
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
