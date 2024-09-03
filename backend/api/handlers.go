package api

import (
    "encoding/json"
    "net/http"
    "breadboard-simulator/circuit"
)

func SimulateHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Components  []circuit.Component
        Connections []circuit.Connection
    }

    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    results, err := circuit.SimulateCircuit(input.Components, input.Connections)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(results)
}