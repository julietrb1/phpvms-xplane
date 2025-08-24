package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/julietrb1/phpvms-xplane/models"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
	Logger     *slog.Logger
}

type DataResponse[T any] struct {
	Data T `json:"data"`
}

type PaginatedResponse[T any] struct {
	Data  []T `json:"data"`
	Links struct {
		First *string `json:"first"`
		Last  *string `json:"last"`
		Prev  *string `json:"prev"`
		Next  *string `json:"next"`
	} `json:"links"`
	Meta struct {
		CurrentPage int     `json:"current_page"`
		From        int     `json:"from"`
		LastPage    int     `json:"last_page"`
		Path        string  `json:"path"`
		PerPage     int     `json:"per_page"`
		To          int     `json:"to"`
		Total       int     `json:"total"`
		PrevPage    *string `json:"prev_page"`
		NextPage    *string `json:"next_page"`
	} `json:"meta"`
}

type PrefilePIREPRequest struct {
	AirlineID          int                    `json:"airline_id"`
	AircraftID         int                    `json:"aircraft_id"`
	FlightType         string                 `json:"flight_type"`
	FlightNumber       string                 `json:"flight_number"`
	DepartureAirportID string                 `json:"dpt_airport_id"`
	ArrivalAirportID   string                 `json:"arr_airport_id"`
	AlternateAirportID string                 `json:"alt_airport_id"`
	Route              string                 `json:"route"`
	Level              int                    `json:"level"`
	PlannedDistance    int                    `json:"planned_distance"`
	PlannedFlightTime  int                    `json:"planned_flight_time"`
	BlockFuel          int                    `json:"block_fuel"`
	Source             int                    `json:"source"`
	SourceName         string                 `json:"source_name"`
	Fields             map[string]interface{} `json:"fields"`
}

type FlightUpdateRequest struct {
	Status      string `json:"status"`
	Distance    int    `json:"distance"`
	FuelUsedLbs int    `json:"fuel_used"`
	FlightTime  int    `json:"flight_time"`
}

type PositionUpdateRequest struct {
	Lat        float64 `json:"lat"`
	Lon        float64 `json:"lon"`
	AltMSL     int     `json:"altitude_msl"`
	AltAGL     int     `json:"altitude_agl"`
	GS         int     `json:"gs"`
	SimTime    string  `json:"sim_time"`
	DistanceNM int     `json:"distance"`
	Heading    int     `json:"heading"`
	IAS        int     `json:"ias"`
	VSFPM      int     `json:"vs"`
}

type FilePIREPRequest struct {
	FlightTime  int `json:"flight_time"`
	FuelUsedLbs int `json:"fuel_used"`
	Distance    int `json:"distance"`
}

func NewClient(baseURL, apiKey string, logger *slog.Logger) *Client {
	if logger == nil {
		logger = slog.Default()
	}

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	return &Client{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		HTTPClient: httpClient,
		Logger:     logger,
	}
}

func (c *Client) doGenericRequest(ctx context.Context, method, url string, body interface{}, result interface{}) error {
	return c.doRequest(ctx, method, url, body, result, false)
}

func (c *Client) doACARSRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	url := fmt.Sprintf("%s%s", c.BaseURL, path)
	return c.doRequest(ctx, method, url, body, result, true)
}

func (c *Client) doRequest(ctx context.Context, method, url string, body interface{}, result interface{}, includeAPIKey bool) error {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "phpVMS-Go-Client/1.0")
	if includeAPIKey {
		req.Header.Set("X-API-Key", c.APIKey)
	}

	c.Logger.Debug("API request",
		"method", method,
		"url", url,
		"has_body", body != nil,
	)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API error: %s (status %d)", url, resp.StatusCode)
	}

	if result == nil {
		return nil
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	c.Logger.Debug("API response",
		"method", method,
		"url", url,
		"status", resp.StatusCode,
	)

	return nil
}

func (c *Client) PrefilePIREP(ctx context.Context, data PrefilePIREPRequest) (*DataResponse[models.ListedPIREP], error) {
	var result DataResponse[models.ListedPIREP]
	err := c.doACARSRequest(ctx, http.MethodPost, "/api/pireps/prefile", data, &result)
	return &result, err
}

func (c *Client) ListPIREPs(ctx context.Context) (PaginatedResponse[models.ListedPIREP], error) {
	var response PaginatedResponse[models.ListedPIREP]
	err := c.doACARSRequest(ctx, http.MethodGet, "/api/pireps", nil, &response)
	return response, err
}

func (c *Client) UpdatePIREP(ctx context.Context, id string, data FlightUpdateRequest) error {
	path := fmt.Sprintf("/api/pireps/%s", id)
	return c.doACARSRequest(ctx, http.MethodPut, path, data, nil)
}

func (c *Client) FilePIREP(ctx context.Context, id string, data FilePIREPRequest) error {
	path := fmt.Sprintf("/api/pireps/%s/file", id)
	return c.doACARSRequest(ctx, http.MethodPost, path, data, nil)
}

func (c *Client) CancelPIREP(ctx context.Context, id string) error {
	path := fmt.Sprintf("/api/pireps/%s/cancel", id)
	return c.doACARSRequest(ctx, http.MethodDelete, path, nil, nil)
}

func (c *Client) GetPIREP(ctx context.Context, id string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/pireps/%s", id)
	var result map[string]interface{}
	err := c.doACARSRequest(ctx, http.MethodGet, path, nil, &result)
	return result, err
}

func (c *Client) PostACARSPosition(ctx context.Context, id string, position PositionUpdateRequest) error {
	path := fmt.Sprintf("/api/pireps/%s/acars/position", id)
	body := map[string]interface{}{
		"positions": []PositionUpdateRequest{position},
	}
	return c.doACARSRequest(ctx, http.MethodPost, path, body, nil)
}

func (c *Client) PostACARSLog(ctx context.Context, id string, log map[string]interface{}) error {
	path := fmt.Sprintf("/api/acars/%s/logs", id)
	return c.doACARSRequest(ctx, http.MethodPost, path, log, nil)
}

func (c *Client) PostACARSEvent(ctx context.Context, id string, event map[string]interface{}) error {
	path := fmt.Sprintf("/api/acars/%s/events", id)
	return c.doACARSRequest(ctx, http.MethodPost, path, event, nil)
}

func (c *Client) GetACARSData(ctx context.Context, id string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/acars/%s", id)
	var result map[string]interface{}
	err := c.doACARSRequest(ctx, http.MethodGet, path, nil, &result)
	return result, err
}

func (c *Client) GetFlight(ctx context.Context, id int) (map[string]interface{}, error) {
	path := fmt.Sprintf("/api/flights/%d", id)
	var result map[string]interface{}
	err := c.doACARSRequest(ctx, http.MethodGet, path, nil, &result)
	return result, err
}

func (c *Client) GetFlightAircraft(ctx context.Context, id int) ([]map[string]interface{}, error) {
	path := fmt.Sprintf("/api/flights/%d/aircraft", id)
	var result []map[string]interface{}
	err := c.doACARSRequest(ctx, http.MethodGet, path, nil, &result)
	return result, err
}

func (c *Client) GetCurrentUser(ctx context.Context) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := c.doACARSRequest(ctx, http.MethodGet, "/api/user", nil, &result)
	return result, err
}

func (c *Client) GetUserFleet(ctx context.Context) (DataResponse[[]models.AircraftFleet], error) {
	var result DataResponse[[]models.AircraftFleet]
	err := c.doACARSRequest(ctx, http.MethodGet, "/api/user/fleet", nil, &result)
	return result, err
}

func (c *Client) GetAirlines(ctx context.Context) (*PaginatedResponse[models.Airline], error) {
	var result PaginatedResponse[models.Airline]
	err := c.doACARSRequest(ctx, http.MethodGet, "/api/airlines", nil, &result)
	return &result, err
}

func (c *Client) GetSimbriefOFP(ctx context.Context, simbriefUserID string) (*models.SimBriefOFP, error) {
	var result models.SimBriefOFP
	url := fmt.Sprintf("https://www.simbrief.com/api/xml.fetcher.php?userid=%s&json=1", simbriefUserID)
	err := c.doGenericRequest(ctx, http.MethodGet, url, nil, &result)
	return &result, err
}
