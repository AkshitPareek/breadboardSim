package simulation

import (
	"fmt"
	"strings"
)

type Component struct {
	ID    string
	Type  string
	Value float64 // Resistance for resistors, voltage for sources, etc.
	// Add other properties as needed
}

type Connection struct {
	From string
	To   string
}

type Circuit struct {
	Components  []Component
	Connections []Connection
}

func CalculateVoltageAndCurrent(circuit Circuit) (map[string]float64, map[string]float64, error) {
	voltages := make(map[string]float64)
	currents := make(map[string]float64)

	// Debug: Print input circuit
	fmt.Printf("Input circuit: %+v\n", circuit)

	// Identify the voltage source
	var voltageSource *Component
	for i, comp := range circuit.Components {
		if comp.Type == "battery" {
			if voltageSource != nil {
				return nil, nil, fmt.Errorf("multiple voltage sources found in the circuit")
			}
			voltageSource = &circuit.Components[i]
		}
	}

	if voltageSource == nil {
		return nil, nil, fmt.Errorf("no voltage source found in the circuit")
	}

	// Debug: Print voltage source
	fmt.Printf("Voltage source: %+v\n", voltageSource)

	// Sum the resistances
	var totalResistance float64
	for _, comp := range circuit.Components {
		if comp.Type == "resistor" {
			totalResistance += comp.Value
		}
	}

	if totalResistance == 0 {
		return nil, nil, fmt.Errorf("total resistance is zero")
	}

	// Debug: Print total resistance
	fmt.Printf("Total resistance: %f\n", totalResistance)

	// Calculate current using Ohm's law: I = V / R
	current := voltageSource.Value / totalResistance

	// Debug: Print calculated current
	fmt.Printf("Calculated current: %f\n", current)

	// Calculate voltage drops across each resistor
	for _, comp := range circuit.Components {
		if comp.Type == "resistor" {
			voltages[comp.ID] = current * comp.Value
		} else if comp.Type == "battery" {
			voltages[comp.ID] = comp.Value
		}
		currents[comp.ID] = current
	}

	// Debug: Print final voltages and currents
	fmt.Printf("Final voltages: %+v\n", voltages)
	fmt.Printf("Final currents: %+v\n", currents)

	return voltages, currents, nil
}

func ValidateCircuit(circuit Circuit) error {
	// Ensure there is exactly one voltage source
	var voltageSourceCount int
	for _, comp := range circuit.Components {
		if comp.Type == "battery" {
			voltageSourceCount++
		}
	}
	if voltageSourceCount != 1 {
		return fmt.Errorf("circuit must have exactly one voltage source")
	}

	// Ensure all components are connected
	componentMap := make(map[string]bool)
	for _, comp := range circuit.Components {
		componentMap[comp.ID] = false
	}
	for _, conn := range circuit.Connections {
		from := strings.TrimPrefix(conn.From, "custom-")
		to := strings.TrimPrefix(conn.To, "custom-")
		componentMap[from] = true
		componentMap[to] = true
	}
	for id, connected := range componentMap {
		if !connected {
			return fmt.Errorf("component %s is not connected", id)
		}
	}

	return nil
}