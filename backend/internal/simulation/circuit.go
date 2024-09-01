package simulation

import (
	"math"
)

type Component struct {
	ID    string
	Type  string
	Value float64 // Resistance for resistors, voltage for sources, etc.
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

	// For now, let's assume a simple series circuit with a voltage source and resistors
	var totalResistance float64
	var voltageSource *Component

	for _, comp := range circuit.Components {
		if comp.Type == "resistor" {
			totalResistance += comp.Value
		} else if comp.Type == "voltage_source" {
			voltageSource = &comp
		}
	}

	if voltageSource == nil {
		return nil, nil, fmt.Errorf("no voltage source found in the circuit")
	}

	// Calculate current using Ohm's law: I = V / R
	current := voltageSource.Value / totalResistance

	// Calculate voltage drops across each component
	for _, comp := range circuit.Components {
		if comp.Type == "resistor" {
			voltages[comp.ID] = current * comp.Value
		} else if comp.Type == "voltage_source" {
			voltages[comp.ID] = comp.Value
		}
		currents[comp.ID] = current
	}

	return voltages, currents, nil
}