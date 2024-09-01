package simulation

import (
	"fmt"
	"testing"
)

func TestCalculateVoltageAndCurrent(t *testing.T) {
	tests := []struct {
		name         string
		circuit      Circuit
		wantVoltages map[string]float64
		wantCurrents map[string]float64
		wantErr      bool
	}{
		{
			name: "Simple series circuit",
			circuit: Circuit{
				Components: []Component{
					{ID: "B1", Type: "battery", Value: 9},
					{ID: "R1", Type: "resistor", Value: 100},
					{ID: "R2", Type: "resistor", Value: 200},
				},
				Connections: []Connection{
					{From: "B1_1", To: "R1_1"},
					{From: "R1_2", To: "R2_1"},
					{From: "R2_2", To: "B1_2"},
				},
			},
			wantVoltages: map[string]float64{
				"B1_1": 9, "B1_2": 0,
				"R1_1": 3, "R1_2": 0,
				"R2_1": 6, "R2_2": 0,
			},
			wantCurrents: map[string]float64{
				"B1": 0.03, "R1": 0.03, "R2": 0.03,
			},
			wantErr: false,
		},
		{
			name: "Invalid circuit - no battery",
			circuit: Circuit{
				Components: []Component{
					{ID: "R1", Type: "resistor", Value: 100},
					{ID: "R2", Type: "resistor", Value: 200},
				},
				Connections: []Connection{
					{From: "R1_1", To: "R2_1"},
					{From: "R1_2", To: "R2_2"},
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid circuit - multiple batteries",
			circuit: Circuit{
				Components: []Component{
					{ID: "B1", Type: "battery", Value: 9},
					{ID: "B2", Type: "battery", Value: 3},
					{ID: "R1", Type: "resistor", Value: 100},
				},
				Connections: []Connection{
					{From: "B1_1", To: "R1_1"},
					{From: "R1_2", To: "B2_1"},
					{From: "B2_2", To: "B1_2"},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			voltages, currents, err := CalculateVoltageAndCurrent(tt.circuit)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateVoltageAndCurrent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				compareResults(t, "Voltages", voltages, tt.wantVoltages)
				compareResults(t, "Currents", currents, tt.wantCurrents)
			}
		})
	}
}

func compareResults(t *testing.T, name string, got, want map[string]float64) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("%s: got %d results, want %d", name, len(got), len(want))
		return
	}
	for k, v := range want {
		if gotV, ok := got[k]; !ok {
			t.Errorf("%s: missing key %s", name, k)
		} else if !almostEqual(gotV, v, 1e-6) {
			t.Errorf("%s: for %s, got %v, want %v", name, k, gotV, v)
		}
	}
}

func TestPreprocessCircuit(t *testing.T) {
	circuit := Circuit{
		Components: []Component{
			{ID: "B1", Type: "battery", Value: 9},
			{ID: "R1", Type: "resistor", Value: 100},
			{ID: "R2", Type: "resistor", Value: 200},
		},
		Connections: []Connection{
			{From: "B1_1", To: "R1_1"},
			{From: "R1_2", To: "R2_1"},
			{From: "R2_2", To: "B1_2"},
		},
	}

	err := preprocessCircuit(&circuit)
	if err != nil {
		t.Errorf("preprocessCircuit() error = %v", err)
		return
	}

	// Check number of nodes
	expectedNodes := 3
	fmt.Printf("Expected nodes: %d\n", expectedNodes)
	fmt.Printf("Actual nodes: %d\n", len(circuit.Nodes))
	if len(circuit.Nodes) != expectedNodes {
		t.Errorf("Expected %d nodes, got %d", expectedNodes, len(circuit.Nodes))
	}

	// Check specific nodes
	expectedNodeNames := map[string]bool{
		"B1_1": true,
		"B1_2": true,
		"R1_2": true,
	}
	fmt.Println("Expected node names:", expectedNodeNames)
	fmt.Println("Actual nodes:")
	for _, node := range circuit.Nodes {
		fmt.Printf("  - %s\n", node.Name)
		if !expectedNodeNames[node.Name] {
			t.Errorf("Unexpected node: %s", node.Name)
		}
		delete(expectedNodeNames, node.Name)
	}
	for name := range expectedNodeNames {
		t.Errorf("Missing expected node: %s", name)
	}

	// Check number of branches
	expectedBranches := 3 // B1, R1, R2
	fmt.Printf("Expected branches: %d\n", expectedBranches)
	fmt.Printf("Actual branches: %d\n", len(circuit.Branches))
	if len(circuit.Branches) != expectedBranches {
		t.Errorf("Expected %d branches, got %d", expectedBranches, len(circuit.Branches))
	}

	// Check if all components are in branches
	componentIDs := make(map[string]bool)
	fmt.Println("Actual branches:")
	for _, branch := range circuit.Branches {
		fmt.Printf("  - Component: %s, From: %d, To: %d\n", branch.Component.ID, branch.NodeFrom, branch.NodeTo)
		componentIDs[branch.Component.ID] = true
	}
	for _, comp := range circuit.Components {
		if !componentIDs[comp.ID] {
			t.Errorf("Component %s not found in branches", comp.ID)
		}
	}

	// Print debug information
	fmt.Println("Nodes:", circuit.Nodes)
	fmt.Println("Branches:", circuit.Branches)
}
