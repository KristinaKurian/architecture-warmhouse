#!/usr/bin/env sh
set -eu

if docker compose version >/dev/null 2>&1; then
  COMPOSE="docker compose"
elif docker-compose version >/dev/null 2>&1; then
  COMPOSE="docker-compose"
else
  echo "Docker Compose is not installed." >&2
  exit 1
fi

echo "Starting Smart Home, PostgreSQL and temperature-api..."
$COMPOSE up --build -d

echo "Waiting for services to become healthy..."
for attempt in $(seq 1 60); do
  postgres_status=$(docker inspect --format='{{if .State.Health}}{{.State.Health.Status}}{{else}}{{.State.Status}}{{end}}' smarthome-postgres 2>/dev/null || true)
  temperature_status=$(docker inspect --format='{{if .State.Health}}{{.State.Health.Status}}{{else}}{{.State.Status}}{{end}}' temperature-api 2>/dev/null || true)
  app_status=$(docker inspect --format='{{.State.Status}}' smarthome-app 2>/dev/null || true)

  if [ "$postgres_status" = "healthy" ] && [ "$temperature_status" = "healthy" ] && [ "$app_status" = "running" ]; then
    echo "All services are ready."
    echo "Smart Home API: http://localhost:8080"
    echo "Temperature API: http://localhost:8081/temperature?location=Living%20Room"
    exit 0
  fi

  echo "Attempt $attempt/60: postgres=$postgres_status temperature-api=$temperature_status app=$app_status"
  sleep 1
done

echo "Services did not become ready. Current logs:" >&2
$COMPOSE logs --tail=100 >&2
exit 1
