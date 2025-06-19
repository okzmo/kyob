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
	Id       string `json:"id"`
}

type RPMAssignAvatar struct {
	Id       string `json:"id"`
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

func (rpm *ReadyPlayerMeService) CreateDefaultBody(ctx context.Context) (string, error) {
	rpmUser, err := rpm.CreateAnonymousUser(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create anonymous user: %w", err)
	}

	templates, err := rpm.GetTemplates(ctx, rpmUser.AccessToken)
	if err != nil {
		return "", fmt.Errorf("failed to get templates on RPM: %w", err)
	}

	template := templates[0]

	avatarId, err := rpm.AssignAvatar(ctx, rpmUser.AccessToken, template.Id)
	if err != nil {
		return "", fmt.Errorf("failed to assign avatar: %w", err)
	}

	err = rpm.SaveAvatar(ctx, rpmUser.AccessToken, avatarId)
	if err != nil {
		return "", fmt.Errorf("failed to save avatar: %w", err)
	}

	return fmt.Sprintf("https://models.readyplayer.me/%s.glb", avatarId), nil
}

func (rpm *ReadyPlayerMeService) CreateAnonymousUser(ctx context.Context) (*RPMUser, error) {
	response, err := makeRPMRequest[RPMUser](ctx, fmt.Sprintf("https://%s.readyplayer.me", rpm.Subdomain), rpm.HTTPClient, http.MethodPost, "/api/users", "", nil)
	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func (rpm *ReadyPlayerMeService) GetTemplates(ctx context.Context, token string) ([]RPMTemplate, error) {
	response, err := makeRPMRequest[[]RPMTemplate](ctx, rpm.BaseURL, rpm.HTTPClient, http.MethodGet, "/v2/avatars/templates", token, nil)
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

	fmt.Println(templateID)
	assignment, err := makeRPMRequest[RPMAssignAvatar](ctx, rpm.BaseURL, rpm.HTTPClient, http.MethodPost, "/v2/avatars/templates/"+templateID, token, payload)
	if err != nil {
		return "", err
	}

	return assignment.Data.Id, nil
}

func (rpm *ReadyPlayerMeService) SaveAvatar(ctx context.Context, token, avatarId string) error {
	_, err := makeRPMRequest[any](ctx, rpm.BaseURL, rpm.HTTPClient, http.MethodPut, "/v2/avatars/"+avatarId, token, nil)
	if err != nil {
		return err
	}

	return nil
}

func makeRPMRequest[T any](ctx context.Context, baseUrl string, client *http.Client, method, path, token string, payload any) (*RPMResponse[T], error) {
	var body io.Reader

	if payload != nil {
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(jsonPayload)
	}

	req, err := http.NewRequestWithContext(ctx, method, baseUrl+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(res.Body)
		fmt.Println(string(bodyBytes))
		return nil, fmt.Errorf("RPM API error: %d", res.StatusCode)
	}

	var result RPMResponse[T]
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
