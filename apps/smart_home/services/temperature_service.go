package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// TemperatureService handles fetching temperature data from temperature-api.
type TemperatureService struct {
	BaseURL    string
	HTTPClient *http.Client
}

// TemperatureResponse represents the response from temperature-api.
type TemperatureResponse struct {
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
	Timestamp   time.Time `json:"timestamp"`
	Location    string    `json:"location"`
	Status      string    `json:"status"`
	SensorID    string    `json:"sensor_id"`
	SensorType  string    `json:"sensor_type"`
	Description string    `json:"description"`
}

// NewTemperatureService creates a new temperature service.
func NewTemperatureService(baseURL string) *TemperatureService {
	return &TemperatureService{
		BaseURL: strings.TrimRight(baseURL, "/"),
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetTemperature fetches temperature data for a location.
func (s *TemperatureService) GetTemperature(location string) (*TemperatureResponse, error) {
	return s.GetTemperatureForSensor(location, "")
}

// GetTemperatureByID fetches temperature data using a sensor identifier.
func (s *TemperatureService) GetTemperatureByID(sensorID string) (*TemperatureResponse, error) {
	return s.GetTemperatureForSensor("", sensorID)
}

// GetTemperatureForSensor fetches a fresh value using both the stored location
// and sensor identifier. temperature-api fills in either value when it is absent.
func (s *TemperatureService) GetTemperatureForSensor(location, sensorID string) (*TemperatureResponse, error) {
	endpoint, err := url.Parse(s.BaseURL + "/temperature")
	if err != nil {
		return nil, fmt.Errorf("build temperature API URL: %w", err)
	}

	query := endpoint.Query()
	if location != "" {
		query.Set("location", location)
	}
	if sensorID != "" {
		query.Set("sensorId", sensorID)
	}
	endpoint.RawQuery = query.Encode()

	resp, err := s.HTTPClient.Get(endpoint.String())
	if err != nil {
		return nil, fmt.Errorf("fetch temperature data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("temperature API returned status %d", resp.StatusCode)
	}

	var temperatureResp TemperatureResponse
	if err := json.NewDecoder(resp.Body).Decode(&temperatureResp); err != nil {
		return nil, fmt.Errorf("decode temperature response: %w", err)
	}

	return &temperatureResp, nil
}
