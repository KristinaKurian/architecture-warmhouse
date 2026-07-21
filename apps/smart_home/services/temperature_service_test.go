package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTemperatureForSensorUsesQueryParameters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/temperature" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("location"); got != "Living Room" {
			t.Fatalf("unexpected location: %q", got)
		}
		if got := r.URL.Query().Get("sensorId"); got != "1" {
			t.Fatalf("unexpected sensorId: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"value":       22.4,
			"unit":        "°C",
			"timestamp":   "2026-07-21T00:00:00Z",
			"location":    "Living Room",
			"status":      "active",
			"sensor_id":   "1",
			"sensor_type": "temperature",
			"description": "test",
		})
	}))
	defer server.Close()

	service := NewTemperatureService(server.URL)
	result, err := service.GetTemperatureForSensor("Living Room", "1")
	if err != nil {
		t.Fatalf("GetTemperatureForSensor returned error: %v", err)
	}
	if result.Value != 22.4 || result.Location != "Living Room" || result.SensorID != "1" {
		t.Fatalf("unexpected response: %+v", result)
	}
}
