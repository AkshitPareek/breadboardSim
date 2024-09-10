package circuit

import (
	"gonum.org/v1/gonum/mat"
	// "strconv"
	"strings"
)

var globalNodeMapping = make(map[string]string)

// sourcePolarity maps voltage source IDs to their polarity and connected nodes
// The int represents the node number, and the string represents the polarity ('+' or '-')
var sourcePolarity = make(map[string]struct {
	NodeNumber int
	Polarity   string
})

var sourceIndex = make(map[string]int)

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
	nodeNumbers := make(map[string]int)
	nodeComponents := make(map[string][]string)
	currentNodeNumber := 1

	// Assign ground node
	nodeNumbers["ground"] = 0
	nodeComponents["ground"] = []string{}

	// First pass: identify all nodes and assign numbers
	for _, conn := range c.Connections {
		for _, node := range []string{conn.From, conn.To} {
			if node != "ground" && !strings.HasPrefix(node, "n") {
				continue
			}
			if _, exists := nodeNumbers[node]; !exists {
				nodeNumbers[node] = currentNodeNumber
				nodeComponents[node] = []string{}
				currentNodeNumber++
			}
		}
	}

	// Second pass: populate nodeComponents
	for _, conn := range c.Connections {
		fromNode := getNodeName(conn.From)
		toNode := getNodeName(conn.To)
		
		if fromNode != "" {
			component := getComponentName(conn.To)
			if component != "" {
				nodeComponents[fromNode] = appendUnique(nodeComponents[fromNode], component)
			}
		}
		
		if toNode != "" {
			component := getComponentName(conn.From)
			if component != "" {
				nodeComponents[toNode] = appendUnique(nodeComponents[toNode], component)
			}
		}
	}

	return nodeNumbers, nodeComponents
}

func getNodeName(s string) string {
	if s == "ground" || strings.HasPrefix(s, "n") {
		return s
	}
	return ""
}

func getComponentName(s string) string {
	if s != "ground" && !strings.HasPrefix(s, "n") {
		return s
	}
	return ""
}

func appendUnique(slice []string, item string) []string {
	for _, element := range slice {
		if element == item {
			return slice
		}
	}
	return append(slice, item)
}

func findConnectedNode(c *Circuit, componentID string) string {
	for _, conn := range c.Connections {
		if conn.From == componentID {
			if strings.HasPrefix(conn.To, "n") || conn.To == "ground" {
				return conn.To
			}
		}
		if conn.To == componentID {
			if strings.HasPrefix(conn.From, "n") || conn.From == "ground" {
				return conn.From
			}
		}
	}
	return ""
}

func buildMNAMatrices(c *Circuit, nodeNumbers map[string]int, nodeComponents map[string][]string) (*mat.Dense, *mat.VecDense, *mat.VecDense) {
	n := len(nodeNumbers)
	m := countVoltageSources(c)
	
	A := mat.NewDense(n+m-1, n+m-1, nil)
	x := buildxMatrix(c, nodeNumbers)
	z := buildzMatrix(c, nodeNumbers, nodeComponents)
	
	// Copy the G matrix into the top-left corner of A
	G := buildGMatrix(c, nodeNumbers, nodeComponents)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-1; j++ {
			A.Set(i, j, G.At(i, j))
		}
	}
	// Build B matrix
	B := buildBMatrix(c, nodeNumbers, nodeComponents)
	
	// Copy B matrix into the top-right corner of A
	for i := 0; i < n-1; i++ {
		for j := 0; j < m; j++ {
			A.Set(i, n-1+j, B.At(i, j))
		}
	}
	
	// Build C matrix
	C := buildCMatrix(c, nodeNumbers, nodeComponents)
	
	// Copy C matrix into the bottom-left corner of A
	for i := 0; i < m; i++ {
		for j := 0; j < n-1; j++ {
			A.Set(n-1+i, j, C.At(i, j))
		}
	}
	
	return A, x, z
}

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

func buildBMatrix(c *Circuit, nodeNumbers map[string]int, nodeComponents map[string][]string) *mat.Dense {
	m := countVoltageSources(c)
	n := len(nodeNumbers)
	B := mat.NewDense(n-1, m, nil) // Initialize B matrix with zeros

	voltageSourceIndex := 0
	for _, comp := range c.Components {
		if comp.Type == Battery  {
			posNode := ""
			negNode := ""
			for _, conn := range c.Connections {
				if conn.From == comp.ID && conn.Polarity == "+" {
					posNode = conn.To
				} else if conn.To == comp.ID && conn.Polarity == "+" {
					posNode = conn.From
				} else if conn.From == comp.ID && conn.Polarity == "-" {
					negNode = conn.To
				} else if conn.To == comp.ID && conn.Polarity == "-" {
					negNode = conn.From
				}
			}

			if posNode != "" && posNode != "ground" {
				B.Set(nodeNumbers[posNode]-1, voltageSourceIndex, 1)
			}
			if negNode != "" && negNode != "ground" {
				B.Set(nodeNumbers[negNode]-1, voltageSourceIndex, -1)
			}

			voltageSourceIndex++
		}
	}

	return B
}

func buildCMatrix(c *Circuit, nodeNumbers map[string]int, nodeComponents map[string][]string) *mat.Dense {
	B := buildBMatrix(c, nodeNumbers, nodeComponents)
	rows, cols := B.Dims()
	C := mat.NewDense(cols, rows, nil)
	C.Copy(B.T())
	return C
}

func buildDMatrix(c *Circuit) *mat.Dense {
	m := countVoltageSources(c)
	D := mat.NewDense(m, m, nil)
	return D
}

func buildxMatrix(c *Circuit, nodeNumbers map[string]int) *mat.VecDense {
	m := countVoltageSources(c)
	n := len(nodeNumbers)
	x := mat.NewVecDense(m+n, nil)

	// First n rows of x are matrix v (node voltages)
	v := buildvMatrix(c, nodeNumbers)
	for k := 0; k < n-1; k++ {
		x.SetVec(k, v.AtVec(k))
	}

	// Next m rows of x are matrix e (voltage sources)
	e := buildeMatrix(c)
	for k := 0; k < m; k++ {
		x.SetVec(n-1+k, e.AtVec(k))
	}
	return x
}

func buildvMatrix(c *Circuit, nodeNumbers map[string]int) *mat.VecDense {
	n := len(nodeNumbers)
	v := mat.NewVecDense(n, nil)
	return v
}

func buildjMatrix(c *Circuit) *mat.VecDense {
	m := countVoltageSources(c)
	j := mat.NewVecDense(m, nil)
	return j
}
func buildzMatrix(c *Circuit, nodeNumbers map[string]int, nodeComponents map[string][]string) *mat.VecDense {
	m := countVoltageSources(c)
	n := len(nodeNumbers)
	z := mat.NewVecDense(m+n-1, nil)

	// First n rows of z are matrix i (currents)
	i := buildiMatrix(c, nodeNumbers, nodeComponents)
	for k := 0; k < n-1; k++ {
		z.SetVec(k, i.AtVec(k))
	}

	// Next m rows of z are matrix e (voltage sources)
	e := buildeMatrix(c)
	for k := 0; k < m; k++ {
		z.SetVec(n-1+k, e.AtVec(k))
	}

	return z
}

func buildiMatrix(c *Circuit, nodeNumbers map[string]int, nodeComponents map[string][]string) *mat.VecDense {
	n := len(nodeNumbers)
	i := mat.NewVecDense(n, nil)
	// TODO: Implement i matrix (currents through voltage sources: none for now)
	for nodeName, components := range nodeComponents {
		if nodeName != "ground" {
			comp := findComponentByID(c, globalNodeMapping[nodeName])
			if comp.Type == CurrentSource {
				currentValue := i.AtVec(nodeNumbers[nodeName]-1)
				i.SetVec(nodeNumbers[nodeName]-1, currentValue+comp.Value)
			}
			for _, compID := range components {
				compIn := findComponentByID(c, compID)
				if compIn.Type == CurrentSource {
					currentValue := i.AtVec(nodeNumbers[nodeName]-1)
					i.SetVec(nodeNumbers[nodeName]-1, currentValue-compIn.Value)
				}
			}
		}
	}
	return i
}

func buildeMatrix(c *Circuit) *mat.VecDense {
	m := countVoltageSources(c)
	e := mat.NewVecDense(m, nil)
	voltIndex := 0
	for _, comp := range c.Components {
		if comp.Type == Battery {
			e.SetVec(voltIndex, comp.Value)
			voltIndex++
		}
	}
	return e
}

func getNodePair(compID string, connections []Connection, nodeMap map[string]int) (int, int) {
	for _, conn := range connections {
		if conn.From == compID {
			return nodeMap[conn.From], nodeMap[conn.To]
		}
		if conn.To == compID {
			return nodeMap[conn.To], nodeMap[conn.From]
		}
	}
	return 0, 0 // Return 0 for ground node if not found
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

func contains(slice []string, item string) bool {
	for _, element := range slice {
		if element == item {
			return true
		}
	}
	return false
}
