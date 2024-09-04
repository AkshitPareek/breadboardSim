package circuit

import (
	"math"
	"reflect"
	"testing"
	"fmt"
	"gonum.org/v1/gonum/mat"
)

var testCircuit = &Circuit{
	Components: []Component{
		{ID: "V1", Type: Battery, Value: 32},
		{ID: "V2", Type: Battery, Value: 20},
		{ID: "R1", Type: Resistor, Value: 2},
		{ID: "R2", Type: Resistor, Value: 4},
		{ID: "R3", Type: Resistor, Value: 8},
	},
	Connections: []Connection{
		{From: "ground", To: "V1"},
		{From: "V1", To: "R1"},
		{From: "R1", To: "R2"},
		{From: "R1", To: "R3"},
		{From: "R3", To: "ground"},
		{From: "V2", To: "R2"},
		{From: "ground", To: "V2"},
	},
}

func TestBuildGMatrix(t *testing.T) {
	nodeNumbers, nodeComponents := assignNodeNumbers(testCircuit)
	G := buildGMatrix(testCircuit, nodeNumbers, nodeComponents)
	
	fmt.Println("G Matrix:")
	if G == nil {
		t.Fatal("G matrix is nil")
	}
	
	// Print matrix dimensions
	r, c := G.Dims()
	fmt.Printf("Matrix dimensions: %d x %d\n", r, c)
	
	// Print matrix contents
	fmt.Printf("%v\n", mat.Formatted(G, mat.Prefix("    "), mat.Squeeze()))
	
}

func TestBuildMNAMatrices(t *testing.T) {
	nodeNumbers, nodeComponents := assignNodeNumbers(testCircuit)
	A, _, _ := buildMNAMatrices(testCircuit, nodeNumbers, nodeComponents)
	
	fmt.Println("A Matrix:")
	if A == nil {
		t.Fatal("A matrix is nil")
	}

	// Print matrix dimensions
	r, c := A.Dims()
	fmt.Printf("Matrix dimensions: %d x %d\n", r, c)
	
	// Print matrix contents
	fmt.Printf("%v\n", mat.Formatted(A, mat.Prefix("    "), mat.Squeeze()))
}

func TestBuildBMatrix(t *testing.T) {
	nodeNumbers, nodeComponents := assignNodeNumbers(testCircuit)
	B := buildBMatrix(testCircuit, nodeNumbers, nodeComponents)
	
	fmt.Println("B Matrix:")
	if B == nil {
		t.Fatal("B matrix is nil")
	}
	
	// Print matrix dimensions
	r, c := B.Dims()
	fmt.Printf("Matrix dimensions: %d x %d\n", r, c)
	
	// Print matrix contents
	fmt.Printf("%v\n", mat.Formatted(B, mat.Prefix("    "), mat.Squeeze()))
	
	// Expected B matrix
	expectedB := mat.NewDense(3, 2, []float64{
		1, 0,
		0, 0,
		0, 1,
	})
	
	// Compare B with expectedB
	if !mat.EqualApprox(B, expectedB, 1e-10) {
		t.Errorf("B matrix does not match expected values.\nGot:\n%v\nWant:\n%v",
			mat.Formatted(B, mat.Prefix("    "), mat.Squeeze()),
			mat.Formatted(expectedB, mat.Prefix("    "), mat.Squeeze()))
	}
}

func TestAssignNodeNumbers(t *testing.T) {

    nodeNumbers, nodeComponents := assignNodeNumbers(testCircuit)

	expectedNodeComponents := map[string][]string{
        "ground": {"V1", "R3", "V2"},
        "v_1": {"V1", "R1"},
        "v_2": {"R1", "R2", "R3"},
        "v_3": {"V2", "R2"},
    }
    if !reflect.DeepEqual(nodeComponents, expectedNodeComponents) {
        t.Errorf("Node components mismatch. Got %v, want %v", nodeComponents, expectedNodeComponents)
    }

	fmt.Println("Node Numbers:", nodeNumbers)
    fmt.Println("Node Components:", nodeComponents)

}


func isClose(a, b float64) bool {
	const tolerance = 1e-6
	return math.Abs(a-b) < tolerance
}
