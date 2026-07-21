package main

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestResolveLocationAndSensorID(t *testing.T) {
	tests := []struct {
		name             string
		location         string
		sensorID         string
		expectedLocation string
		expectedSensorID string
	}{
		{"location to id", "Living Room", "", "Living Room", "1"},
		{"id to location", "", "2", "Bedroom", "2"},
		{"unknown location", "Garage", "", "Garage", "0"},
		{"unknown id", "", "99", "Unknown", "99"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			location, sensorID := resolveLocationAndSensorID(tc.location, tc.sensorID)
			if location != tc.expectedLocation || sensorID != tc.expectedSensorID {
				t.Fatalf("got (%q, %q), want (%q, %q)", location, sensorID, tc.expectedLocation, tc.expectedSensorID)
			}
		})
	}
}

func TestTemperatureChangesBetweenCalls(t *testing.T) {
	handler := routes(newTemperatureGenerator())
	values := make([]float64, 0, 2)

	for range 2 {
		req := httptest.NewRequest("GET", "/temperature?location=Living%20Room", nil)
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, req)
		if resp.Code != 200 {
			t.Fatalf("unexpected status: %d", resp.Code)
		}

		var payload temperatureResponse
		if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		values = append(values, payload.Value)
		if payload.Location != "Living Room" || payload.SensorID != "1" {
			t.Fatalf("unexpected mapping: location=%q sensor_id=%q", payload.Location, payload.SensorID)
		}
	}

	if values[0] == values[1] {
		t.Fatalf("consecutive values must differ: %v", values)
	}
}
