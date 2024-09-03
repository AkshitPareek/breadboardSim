package circuit

import (
	// "fmt"
	// "sort"
	"strconv"

	"gonum.org/v1/gonum/mat"
)

func SolveCircuit(c *Circuit) (map[string]float64, error) {
	// 1. Assign node numbers
	nodeMap, nodeComponents := assignNodeNumbers(c)

	// 2. Build MNA matrices
	A, x, z := buildMNAMatrices(c, nodeMap, nodeComponents)

	// 3. Solve the system
	err := x.SolveVec(A, z)
	if err != nil {
		return nil, err
	}

	// 4. Extract results
	return nil, nil
}

func assignNodeNumbers(c *Circuit) (map[string]int, map[string][]string) {

	// Initialize node numbers
	nodeNumbers := make(map[string]int)
	nodeComponents := make(map[string][]string)
	nodeNumbers["ground"] = 0 // Ground is always node 0
	nodeComponents["ground"] = []string{}
	nextNode := 1

	visitedFrom := make(map[string]bool)
	visitedTo := make(map[string]bool)

	// Assign node numbers to all nodes except ground
	for _, conn := range c.Connections {
		currNode := conn.From
		if conn.From == "ground" {
			nodeComponents["ground"] = append(nodeComponents["ground"], conn.To)
			visitedTo[conn.To] = true
			continue
		}
		if conn.To == "ground" {
			nodeComponents["ground"] = append(nodeComponents["ground"], conn.From)
			visitedFrom[conn.From] = true
			continue
		}

		if !visitedFrom[conn.From] && !visitedTo[conn.To] {
			if _, exists := nodeNumbers[currNode]; !exists {
				nodeNumbers[currNode] = nextNode
				nextNode++
			}
			nodeComponents[currNode] = append(nodeComponents[currNode], conn.From, conn.To)
			visitedFrom[conn.From] = true
			visitedTo[conn.To] = true
		} else if visitedFrom[conn.From] && !visitedTo[conn.To] {
			nodeComponents[currNode] = append(nodeComponents[currNode], conn.To)
			visitedTo[conn.To] = true
		} else if !visitedFrom[conn.From] && visitedTo[conn.To] {
			nodeComponents[currNode] = append(nodeComponents[currNode], conn.From)
			visitedFrom[conn.From] = true
		}

	}

	newNodeComponents := make(map[string][]string)
	newNodeComponents["ground"] = nodeComponents["ground"]

	for _, conn := range c.Connections {
		if conn.From != "ground" && conn.To != "ground" {
			nodeName := "v_" + strconv.Itoa(nodeNumbers[conn.From])
			if _, exists := newNodeComponents[nodeName]; !exists {
				newNodeComponents[nodeName] = append(newNodeComponents[nodeName], nodeComponents[conn.From]...)
			}
		}
	}

	newNodeNumbers := make(map[string]int)
    newNodeNumbers["ground"] = 0

	for _, nodeIndex := range nodeNumbers {
		if nodeIndex > 0 {
			newNodeNumbers["v_" + strconv.Itoa(nodeIndex)] = nodeIndex
		}
	}

	nodeComponents = newNodeComponents
	nodeNumbers = newNodeNumbers

	return nodeNumbers, nodeComponents
}

// func appendUnique(slice []string, item string) []string {
//     for _, element := range slice {
//         if element == item {
//             return slice
//         }
//     }
//     return append(slice, item)
// }

func contains(slice []string, item string) bool {
	for _, element := range slice {
		if element == item {
			return true
		}
	}
	return false
}

// Make sure this function is available in your package
// func findComponent(c *Circuit, from, to string) Component {
//     for _, comp := range c.Components {
//         if (comp.ID == from && findConnectionTo(c, comp.ID) == to) ||
//            (comp.ID == to && findConnectionFrom(c, comp.ID) == from) {
//             return comp
//         }
//     }
//     return Component{} // Return an empty component if not found
// }

// /////////////////////////////////////////////
func buildGMatrix(circuit *Circuit, nodeNumbers map[string]int, nodeComponents map[string][]string) *mat.Dense {
	// Determine the size of the matrix
	matrixSize := len(nodeNumbers) // Ensure this reflects the actual number of nodes

	// Initialize the G matrix with the correct size
	G := mat.NewDense(matrixSize-1, matrixSize-1, nil)

	for nodeName, components := range nodeComponents {
		if nodeName == "ground" {
			continue
		}

		nodeIndex := nodeNumbers[nodeName] - 1 // Adjust for 0-based indexing

		// Calculate total conductance for the node
		totalConductance := 0.0
		for _, compID := range components {
			for _, comp := range circuit.Components {
				if comp.ID == compID && comp.Type == Resistor {
					totalConductance += 1.0 / comp.Value
				}
			}
		}

		if totalConductance > 0 {
			G.Set(nodeIndex, nodeIndex, totalConductance)
			// fmt.Println(nodeName, nodeIndex, totalConductance)
		}
	}

	// Build off-diagonal elements of G matrix
	for nodeName1, components1 := range nodeComponents {
		for nodeName2, components2 := range nodeComponents {
			if nodeName1 == "ground" || nodeName2 == "ground" || nodeName1 == nodeName2 {
				continue
			}

			for _, compID := range components1 {
				if contains(components2, compID) {
					comp := findComponentByID(circuit, compID)
					if comp.Type == Resistor {
						conductance := 1.0 / comp.Value * (-1.0)
						i := nodeNumbers[nodeName1] - 1
						j := nodeNumbers[nodeName2] - 1

						// Stamp the negative conductance at (i,j) and (j,i)
						G.Set(i, j, conductance)
						G.Set(j, i, conductance)
                        // fmt.Println(i, j, conductance, comp.Value)
						// fmt.Println(i, j, conductance)
					}
				}
			}
		}
	}

	return G
}

func buildMNAMatrices(c *Circuit, nodeNumbers map[string]int, nodeComponents map[string][]string) (*mat.Dense, *mat.VecDense, *mat.VecDense) {
	n := len(nodeNumbers)
	m := countVoltageSources(c)

	A := mat.NewDense(n+m-1, n+m-1, nil)
	x := mat.NewVecDense(n+m-1, nil)
	z := mat.NewVecDense(n+m-1, nil)

	// Build G matrix
	for nodeName, components := range nodeComponents {
		if nodeName == "ground" {
			continue
		}

		nodeIndex := nodeNumbers[nodeName] - 1 // Adjust for 0-based indexing

		// Calculate total conductance for the node
		totalConductance := 0.0
		for _, compID := range components {
			for _, comp := range c.Components {
				if comp.ID == compID && comp.Type == Resistor {
					totalConductance += 1.0 / comp.Value
				}
			}
		}

		if totalConductance > 0 {
			A.Set(nodeIndex, nodeIndex, totalConductance)
		}
	}

	// Build off-diagonal elements of G matrix
	for nodeName1, components1 := range nodeComponents {
		for nodeName2, components2 := range nodeComponents {
			if nodeName1 == "ground" || nodeName2 == "ground" || nodeName1 == nodeName2 {
				continue
			}

			for _, compID := range components1 {
				if contains(components2, compID) {
					comp := findComponentByID(c, compID)
					if comp.Type == Resistor {
						conductance := 1.0 / comp.Value * (-1.0)
						i := nodeNumbers[nodeName1] - 1
						j := nodeNumbers[nodeName2] - 1

						// Stamp the negative conductance at (i,j) and (j,i)
						A.Set(i, j, conductance)
						A.Set(j, i, conductance)
						
					}
				}
			}
		}
	}

	// G Matrix built

	// Build B matrix

	// Build A and z matrices for voltage sources
	// voltIndex := 0
	// for _, comp := range c.Components {
	//     if comp.Type == Battery {
	//         n1, n2 := getNodePair(comp.ID, c.Connections, nodeNumbers)
	//         stampVoltageSource(A, z, n1, n2, voltIndex, comp.Value, n-1)
	//         voltIndex++
	//     }
	// }

	return A, x, z
}

func countVoltageSources(c *Circuit) int {
	count := 0
	for _, comp := range c.Components {
		if comp.Type == Battery {
			count++
		}
	}
	return count
}

func findComponentByID(c *Circuit, id string) Component {
	for _, comp := range c.Components {
		if comp.ID == id {
			return comp
		}
	}
	return Component{} // Return an empty component if not found
}

// func getNodePair(compID string, connections []Connection, nodeMap map[string]int) (int, int) {
//     for _, conn := range connections {
//         if conn.From == compID {
//             return nodeMap[conn.From], nodeMap[conn.To]
//         }
//         if conn.To == compID {
//             return nodeMap[conn.To], nodeMap[conn.From]
//         }
//     }
//     return 0, 0 // Return 0 for ground node if not found
// }

// func getVoltageSourceIndex(voltageID string, c *Circuit) int {
//     index := 0
//     for _, comp := range c.Components {
//         if comp.Type == Battery {
//             if comp.ID == voltageID {
//                 return index
//             }
//             index++
//         }
//     }
//     return -1 // Error case
// }

// func stampResistor(G *mat.Dense, n1, n2 int, conductance float64) {
//     if n1 >= 0 {
//         G.Set(n1, n1, G.At(n1, n1)+conductance)
//     }
//     if n2 >= 0 {
//         G.Set(n2, n2, G.At(n2, n2)+conductance)
//     }
//     if n1 >= 0 && n2 >= 0 {
//         G.Set(n1, n2, G.At(n1, n2)-conductance)
//         G.Set(n2, n1, G.At(n2, n1)-conductance)
//     }
// }

// // Update the function signature to include numNodes
// func stampVoltageSource(A *mat.Dense, z *mat.VecDense, n1, n2, voltIndex int, voltage float64, numNodes int) {
//     if n1 > 0 {
//         A.Set(n1-1, numNodes-1+voltIndex, 1)
//         A.Set(numNodes-1+voltIndex, n1-1, 1)
//     }
//     if n2 > 0 {
//         A.Set(n2-1, numNodes-1+voltIndex, -1)
//         A.Set(numNodes-1+voltIndex, n2-1, -1)
//     }
//     z.SetVec(numNodes-1+voltIndex, voltage)
// }

// func solveLinearSystem(A *mat.Dense, z *mat.VecDense) *mat.VecDense {
//     var x mat.VecDense
//     err := x.SolveVec(A, z)
//     if err != nil {
//         // Handle error
//         return mat.NewVecDense(1, []float64{0}) // Return a dummy vector in case of error
//     }
//     return &x
// }

// func extractResults(x *mat.VecDense, nodeMap map[string]int, c *Circuit) map[string]float64 {
//     results := make(map[string]float64)
//     numNodes := len(nodeMap) - 1

//     // Extract node voltages
//     for _, index := range nodeMap {
//         if index > 0 { // Skip ground node
//             results[fmt.Sprintf("v_%d", index)] = x.AtVec(index - 1)
//         }
//     }
// }

//     // Extract currents through voltage sources
//     voltageSourceIndex := 0
//     for _, comp := range c.Components {
//         if comp.Type == Battery {
//             current := x.AtVec(numNodes + voltageSourceIndex)
//             results[fmt.Sprintf("I_%s", comp.ID)] = current
//             voltageSourceIndex++
//         }
//     }

//     return results
// }

// func solveMNA(A *mat.Dense, z *mat.VecDense) *mat.VecDense {
//     var LU mat.LU
//     LU.Factorize(A)
//     x := mat.NewVecDense(z.Len(), nil)
//     LU.SolveVecTo(x, false, z)
//     return x
// }

// func findConnectionTo(c *Circuit, compID string) string {
//     for _, conn := range c.Connections {
//         if conn.From == compID {
//             return conn.To
//         }
//     }
//     return "" // Return an empty string if not found
// }

// func findConnectionFrom(c *Circuit, compID string) string {
//     for _, conn := range c.Connections {
//         if conn.To == compID {
//             return conn.From
//         }
//     }
//     return "" // Return an empty string if not found
// }
