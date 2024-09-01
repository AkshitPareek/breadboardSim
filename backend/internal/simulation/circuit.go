package simulation

import (
	"fmt"
	"math"
)

type Circuit struct {
	Components  []Component
	Connections []Connection
	Nodes       []Node
	Branches    []Branch
}

type Component struct {
	ID    string
	Type  string
	Value float64
}

type Connection struct {
	From string
	To   string
}

type Node struct {
	ID    int
	Name  string
	Connections []string
}

type Branch struct {
	ID        int
	Component Component
	NodeFrom  int
	NodeTo    int
}

func CalculateVoltageAndCurrent(circuit Circuit) (map[string]float64, map[string]float64, error) {
	if err := validateCircuit(circuit); err != nil {
		return nil, nil, err
	}

	if err := preprocessCircuit(&circuit); err != nil {
		return nil, nil, err
	}

	// TODO: Implement MNA algorithm here
	// For now, we'll keep the existing simple calculation

	battery := findBattery(circuit)
	if battery == nil {
		return nil, nil, fmt.Errorf("no battery found in the circuit")
	}

	totalResistance := calculateTotalResistance(circuit)
	totalCurrent := battery.Value / totalResistance

	voltages := make(map[string]float64)
	currents := make(map[string]float64)

	voltages[battery.ID+"_1"] = battery.Value
	voltages[battery.ID+"_2"] = 0
	currents[battery.ID] = totalCurrent

	for _, comp := range circuit.Components {
		if comp.Type == "resistor" {
			voltage := comp.Value * totalCurrent
			voltages[comp.ID+"_1"] = voltage
			voltages[comp.ID+"_2"] = 0
			currents[comp.ID] = totalCurrent
		}
	}

	return voltages, currents, nil
}

func validateCircuit(circuit Circuit) error {
	if len(circuit.Components) == 0 {
		return fmt.Errorf("empty circuit")
	}

	batteryCount := 0
	for _, comp := range circuit.Components {
		if comp.Type == "battery" {
			batteryCount++
		}
		if comp.Type == "resistor" && comp.Value == 0 {
			return fmt.Errorf("zero resistance detected in component %s", comp.ID)
		}
	}

	if batteryCount == 0 {
		return fmt.Errorf("no battery found in the circuit")
	}
	if batteryCount > 1 {
		return fmt.Errorf("multiple batteries found in the circuit")
	}

	return nil
}

func findBattery(circuit Circuit) *Component {
	for _, comp := range circuit.Components {
		if comp.Type == "battery" {
			return &comp
		}
	}
	return nil
}

func calculateTotalResistance(circuit Circuit) float64 {
	totalResistance := 0.0
	for _, comp := range circuit.Components {
		if comp.Type == "resistor" {
			totalResistance += comp.Value
		}
	}
	return totalResistance
}

// Helper function to compare float64 values with a tolerance
func almostEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

// New function to preprocess the circuit
func preprocessCircuit(circuit *Circuit) error {
	fmt.Println("Preprocessing circuit...")

	// Step 1: Identify and number nodes
	nodeMap := make(map[string]int)
	nodeCounter := 0

	fmt.Println("Initial connections:", circuit.Connections)

	// First pass: create nodes for unique connection points
	for _, comp := range circuit.Components {
		node1 := comp.ID + "_1"
		node2 := comp.ID + "_2"
		if _, exists := nodeMap[node1]; !exists {
			nodeMap[node1] = nodeCounter
			circuit.Nodes = append(circuit.Nodes, Node{ID: nodeCounter, Name: node1, Connections: []string{}})
			nodeCounter++
			fmt.Printf("Created node: %s (ID: %d)\n", node1, nodeCounter-1)
		}
		if _, exists := nodeMap[node2]; !exists {
			nodeMap[node2] = nodeCounter
			circuit.Nodes = append(circuit.Nodes, Node{ID: nodeCounter, Name: node2, Connections: []string{}})
			nodeCounter++
			fmt.Printf("Created node: %s (ID: %d)\n", node2, nodeCounter-1)
		}
	}

	fmt.Println("After first pass - Nodes:", circuit.Nodes)
	fmt.Println("NodeMap:", nodeMap)

	// Second pass: merge nodes based on connections
	for _, conn := range circuit.Connections {
		fromID := nodeMap[conn.From]
		toID := nodeMap[conn.To]
		if fromID != toID {
			// Merge the higher ID into the lower ID
			if fromID > toID {
				fromID, toID = toID, fromID
			}
			fmt.Printf("Merging node %d into node %d\n", toID, fromID)
			circuit.Nodes[fromID].Connections = append(circuit.Nodes[fromID].Connections, circuit.Nodes[toID].Connections...)
			circuit.Nodes[fromID].Connections = append(circuit.Nodes[fromID].Connections, conn.To)
			// Update nodeMap
			for k, v := range nodeMap {
				if v == toID {
					nodeMap[k] = fromID
				}
			}
			// Remove the merged node
			circuit.Nodes = append(circuit.Nodes[:toID], circuit.Nodes[toID+1:]...)
			// Update IDs of nodes after the removed one
			for i := toID; i < len(circuit.Nodes); i++ {
				circuit.Nodes[i].ID--
				for k, v := range nodeMap {
					if v == i+1 {
						nodeMap[k] = i
					}
				}
			}
		}
	}

	fmt.Println("After second pass - Nodes:", circuit.Nodes)
	fmt.Println("Final NodeMap:", nodeMap)

	// Step 3: Create branches
	for i, comp := range circuit.Components {
		fromNode := nodeMap[comp.ID+"_1"]
		toNode := nodeMap[comp.ID+"_2"]
		
		circuit.Branches = append(circuit.Branches, Branch{
			ID:        i,
			Component: comp,
			NodeFrom:  fromNode,
			NodeTo:    toNode,
		})
		fmt.Printf("Created branch: Component %s, From: %d, To: %d\n", comp.ID, fromNode, toNode)
	}

	fmt.Println("Nodes identified:", circuit.Nodes)
	fmt.Println("Branches created:", circuit.Branches)

	return nil
}