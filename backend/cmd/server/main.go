package main

import (
	"encoding/json"
	"fmt"
	// "io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	// Add this import
	"breadboard-simulator/internal/simulation"
)

type BreadboardState struct {
	Components []Component `json:"components"`
	Connections []Connection `json:"connections"`
}

type Component struct {
	ID string `json:"id"`
	Type string `json:"type"`
	Position struct {
		X int `json:"x"`
		Y int `json:"y"`
	} `json:"position"`
}

type Connection struct {
	From string `json:"from"`
	To string `json:"to"`
}

var (
	sessionState BreadboardState
	stateMutex sync.RWMutex
)

// Add this struct
type SimulationRequest struct {
	Components  []simulation.Component  `json:"components"`
	Connections []simulation.Connection `json:"connections"`
}

// Add this function
func handleSimulate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var simReq SimulationRequest
	if err := json.NewDecoder(r.Body).Decode(&simReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	circuit := simulation.Circuit{
		Components:  simReq.Components,
		Connections: simReq.Connections,
	}

	voltages, currents, err := simulation.CalculateVoltageAndCurrent(circuit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"voltages": voltages,
		"currents": currents,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	mux := http.NewServeMux()

	// Add this line
	mux.HandleFunc("/", enableCORS(handleRoot))

	mux.HandleFunc("/api/components", enableCORS(handleComponents))
	mux.HandleFunc("/api/save", enableCORS(handleSave))
	mux.HandleFunc("/api/load", enableCORS(handleLoad))
	mux.HandleFunc("/api/save-file", enableCORS(handleSaveFile))
	mux.HandleFunc("/api/load-file", enableCORS(handleLoadFile))
	mux.HandleFunc("/api/download", enableCORS(handleDownload))
	mux.HandleFunc("/api/upload", enableCORS(handleUpload))
	mux.HandleFunc("/api/simulate", enableCORS(handleSimulate))

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// Add this function
func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Breadboard Simulator API is running")
}

func handleComponents(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling /api/components request")
	components := []string{"resistor", "capacitor", "led", "transistor", "ic"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"components": components,
	})
}

func handleSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var state BreadboardState
	if err := json.NewDecoder(r.Body).Decode(&state); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stateMutex.Lock()
	sessionState = state
	stateMutex.Unlock()

	w.WriteHeader(http.StatusOK)
}

func handleLoad(w http.ResponseWriter, r *http.Request) {
	stateMutex.RLock()
	json.NewEncoder(w).Encode(sessionState)
	stateMutex.RUnlock()
}

func handleSaveFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var state BreadboardState
	if err := json.NewDecoder(r.Body).Decode(&state); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := ioutil.WriteFile("saved_breadboard.json", data, 0644); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleLoadFile(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("saved_breadboard.json")
	if err != nil {
		http.Error(w, "No saved file found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	stateMutex.RLock()
	data, err := json.MarshalIndent(sessionState, "", "  ")
	stateMutex.RUnlock()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=breadboard_state.json")
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	var state BreadboardState
	if err := json.NewDecoder(file).Decode(&state); err != nil {
		http.Error(w, "Invalid JSON file", http.StatusBadRequest)
		return
	}

	stateMutex.Lock()
	sessionState = state
	stateMutex.Unlock()

	fmt.Fprintf(w, "File uploaded and state updated successfully")
}

// Add this function to enable CORS
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	}
}