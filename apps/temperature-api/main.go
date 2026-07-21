package main

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const defaultPort = "8081"

type temperatureResponse struct {
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
	Timestamp   time.Time `json:"timestamp"`
	Location    string    `json:"location"`
	Status      string    `json:"status"`
	SensorID    string    `json:"sensor_id"`
	SensorType  string    `json:"sensor_type"`
	Description string    `json:"description"`
}

type temperatureGenerator struct {
	mu   sync.Mutex
	rand *rand.Rand
	last map[string]int
}

func newTemperatureGenerator() *temperatureGenerator {
	return &temperatureGenerator{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
		last: make(map[string]int),
	}
}

// next returns a random indoor temperature from 15.0 to 30.0 °C.
// It also guarantees that two consecutive values for the same sensor differ.
func (g *temperatureGenerator) next(key string) float64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	const variants = 151 // 15.0, 15.1, ..., 30.0
	value := g.rand.Intn(variants)
	if previous, ok := g.last[key]; ok && value == previous {
		value = (value + 1 + g.rand.Intn(variants-1)) % variants
	}
	g.last[key] = value

	return float64(150+value) / 10
}

func resolveLocationAndSensorID(location, sensorID string) (string, string) {
	location = strings.TrimSpace(location)
	sensorID = strings.TrimSpace(sensorID)

	// If no location is provided, use a default based on sensor ID.
	if location == "" {
		switch sensorID {
		case "1":
			location = "Living Room"
		case "2":
			location = "Bedroom"
		case "3":
			location = "Kitchen"
		default:
			location = "Unknown"
		}
	}

	// If no sensor ID is provided, generate one based on location.
	if sensorID == "" {
		switch strings.ToLower(location) {
		case strings.ToLower("Living Room"):
			sensorID = "1"
		case strings.ToLower("Bedroom"):
			sensorID = "2"
		case strings.ToLower("Kitchen"):
			sensorID = "3"
		default:
			sensorID = "0"
		}
	}

	return location, sensorID
}

func routes(generator *temperatureGenerator) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("GET /temperature", func(w http.ResponseWriter, r *http.Request) {
		location, sensorID := resolveLocationAndSensorID(
			r.URL.Query().Get("location"),
			firstNonEmpty(r.URL.Query().Get("sensorId"), r.URL.Query().Get("sensor_id")),
		)
		writeTemperature(w, generator, location, sensorID)
	})

	// Compatibility endpoint for the original smart_home client.
	mux.HandleFunc("GET /temperature/{sensorID}", func(w http.ResponseWriter, r *http.Request) {
		location, sensorID := resolveLocationAndSensorID("", r.PathValue("sensorID"))
		writeTemperature(w, generator, location, sensorID)
	})

	return mux
}

func writeTemperature(w http.ResponseWriter, generator *temperatureGenerator, location, sensorID string) {
	key := sensorID + ":" + location
	response := temperatureResponse{
		Value:       generator.next(key),
		Unit:        "°C",
		Timestamp:   time.Now().UTC(),
		Location:    location,
		Status:      "active",
		SensorID:    sensorID,
		SensorType:  "temperature",
		Description: "Random temperature reading for " + location,
	}
	writeJSON(w, http.StatusOK, response)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil && !errors.Is(err, http.ErrHandlerTimeout) {
		log.Printf("failed to encode response: %v", err)
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func main() {
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = defaultPort
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	server := &http.Server{
		Addr:              port,
		Handler:           routes(newTemperatureGenerator()),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("temperature-api listening on %s", port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
