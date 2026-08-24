# ev_routing

EV routing algorithms — an HTTP service that finds the cheapest route for an
electric vehicle between two points, routing through charging stations along
the way when the vehicle's battery can't make the trip on a single charge.

Five search strategies are exposed over the same cost model: Dijkstra, a
genetic algorithm, Variable Neighborhood Search, Branch and Bound, and Ant
Colony Optimization.

## How it works

Given a start point, a finish point, and a pool of candidate charging
stations, the service:

1. Selects a subset of candidate stations within a 5 km corridor of the
   direct start→finish line, capped at a count that shrinks as the vehicle's
   remaining range grows relative to the trip distance.
2. Builds a directed graph over `{start, finish} ∪ candidates`, where each
   edge's cost is the driving cost between the two points plus, if the
   destination is a charging stop, the cost of the cheapest usable slot
   there (tariff × energy, plus a wait-time cost).
3. Searches that graph for the minimum-cost path from start to finish,
   using Dijkstra's algorithm, a genetic algorithm, Variable Neighborhood
   Search, Branch and Bound, or Ant Colony Optimization.

The vehicle's energy consumption is temperature-dependent: it's derived from
the vehicle's passport consumption (`spendOpt`, kWh/100km at 25°C — `E0` in
the cost model) and the trip's ambient temperature, via a piecewise model
that penalizes both cold and hot extremes relative to the 25°C optimum.

The full cost model — the station-selection formula, the temperature/
consumption relationship, and the `Z = Z_s + Z_w + Z_ch` objective (driving
cost, wait cost, charging cost) — is documented alongside its implementation
in `internal/service/cost_model.go`.

## Requirements

- Go 1.26+
- An [OpenRouteService](https://openrouteservice.org/) API key, used to
  compute driving distance/duration between points.

## Configuration

`config/config.yaml` holds the OpenRouteService request URL:

```yaml
openrouteservice:
  request:
    url: https://api.openrouteservice.org/v2/directions/driving-car
```

The API key itself is read from the `OPENROUTESERVICE_API_KEY` environment
variable (see `.env` for local development — do not commit real keys to
version control).

## Running

```sh
go run ./cmd/server
```

The server listens on `:8080`. Each search strategy logs its own progress
and total duration to stdout as it runs.

### Docker

```sh
docker compose up --build
```

Reads `OPENROUTESERVICE_API_KEY` from `.env` (see Configuration above) and
publishes the server on `:8080`.

## API

All endpoints accept the same request body and return the same response
shape; they differ only in which search strategy is used.

### `POST /route/genetic`

### `POST /route/dijkstra`

### `POST /route/vns`

### `POST /route/branch-and-bound`

### `POST /route/aco`

**Request body:**

```json
{
  "startCoords": { "longitude": 0, "latitude": 0 },
  "finishCoords": { "longitude": 0, "latitude": 0 },
  "accLevel": 0,
  "accMax": 0,
  "spendOpt": 0,
  "temperature": 0,
  "filteredStations": [
    {
      "id": 0,
      "coords": { "longitude": 0, "latitude": 0 },
      "plug": "TYPE_2",
      "slots": [
        {
          "id": 0,
          "price": 0,
          "power": 0,
          "waitTime": 0,
          "isActive": true
        }
      ]
    }
  ]
}
```

| Field                        | Description                                                    |
| ----------------------------- | ---------------------------------------------------------------|
| `accLevel`                   | Battery charge level at departure, kWh                          |
| `accMax`                     | Battery's maximum capacity, kWh                                 |
| `spendOpt`                   | Passport specific consumption at 25°C, kWh/100km (E0)            |
| `temperature`                | Ambient temperature, °C                                         |
| `filteredStations[].plug`    | One of `TYPE_2`, `CCS`, `CHADEMO`                                |
| `filteredStations[].slots[]` | The station's charging slots, each with its own tariff/power/wait|

**Response body:** an ordered list of route stops — the start, zero or more
charging stations the route stops at, and the finish — each with the
cumulative distance, cost, and duration to reach it. `routeNode` is an echo
of that stop's station exactly as it was sent in `filteredStations`
(including *all* of its slots, not just the one used); `slotId` says which
one the route actually used there — it's only unique within `routeNode.id`,
not globally. A single route can stop to charge at more than one station:

```json
[
  {
    "routeNode": { "id": -1, "coords": { "longitude": 0, "latitude": 0 }, "plug": "", "slots": [] },
    "distance": 0,
    "cost": 0,
    "chargeDuration": 0,
    "reachDuration": 0,
    "slotId": 0
  },
  {
    "routeNode": {
      "id": 1,
      "coords": { "longitude": 0, "latitude": 0 },
      "plug": "TYPE_2",
      "slots": [
        { "id": 1, "price": 25, "power": 50, "waitTime": 0, "isActive": true },
        { "id": 2, "price": 22, "power": 150, "waitTime": 0.25, "isActive": true }
      ]
    },
    "distance": 120,
    "cost": 450,
    "chargeDuration": 1800000000000,
    "reachDuration": 5400000000000,
    "slotId": 2
  },
  {
    "routeNode": {
      "id": 2,
      "coords": { "longitude": 0, "latitude": 0 },
      "plug": "CCS",
      "slots": [
        { "id": 0, "price": 30, "power": 100, "waitTime": 0, "isActive": true }
      ]
    },
    "distance": 260,
    "cost": 900,
    "chargeDuration": 1200000000000,
    "reachDuration": 9600000000000,
    "slotId": 0
  },
  {
    "routeNode": { "id": -2, "coords": { "longitude": 0, "latitude": 0 }, "plug": "", "slots": [] },
    "distance": 340,
    "cost": 1050,
    "chargeDuration": 0,
    "reachDuration": 12000000000000,
    "slotId": 0
  }
]
```

The start/finish stops have `slots: []` because they're synthetic vertices
(not real stations from `filteredStations`), and `slotId: 0` since they're
never charging stops.

## Project layout

```
cmd/server           entry point
config/              config loading and config.yaml
internal/dto/        request/response types
internal/controller/ HTTP handlers
internal/service/    cost model, routing graph, Dijkstra/genetic/VNS/branch-
                      and-bound/ACO search, OpenRouteService client
Dockerfile           container build
docker-compose.yml   local run via Docker
```
