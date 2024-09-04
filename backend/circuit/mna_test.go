package circuit

import (
	"math"
	// "reflect"
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
		{From: "R3", To: "R2"},
		{From: "R3", To: "ground"},
		{From: "V2", To: "R2"},
		{From: "ground", To: "V2"},
	},
}

var testCircuit2 = &Circuit{
	Components: []Component{
		{ID: "V1", Type: Battery, Value: 32},
		{ID: "V2", Type: Battery, Value: 20},
		{ID: "R1", Type: Resistor, Value: 2},
		{ID: "R2", Type: Resistor, Value: 4},
		{ID: "R3", Type: Resistor, Value: 8},
	},
	Connections: []Connection{
		{From: "ground", To: "R1"},
		{From: "R1", To: "V1"},
		{From: "V1", To: "R2"},
		{From: "V1", To: "R3"},
		{From: "R3", To: "ground"},
		{From: "V2", To: "R2"},
		{From: "ground", To: "V2"},
	},
}

var testCircuit3 = &Circuit{
	Components: []Component{
		{ID: "V1", Type: Battery, Value: 32},
		{ID: "I1", Type: CurrentSource, Value: 10},
		{ID: "R1", Type: Resistor, Value: 2},
		{ID: "R2", Type: Resistor, Value: 4},
		{ID: "R3", Type: Resistor, Value: 8},
	},
	Connections: []Connection{
		{From: "ground", To: "I1"},
		// {From: "I1", To: "R1"},
		// {From: "I1", To: "R2"},
		{From: "V1", To: "I1"},
		{From: "V1", To: "R1"},
		{From: "V1", To: "R2"},
		{From: "R3", To: "V1"},
		{From: "ground", To: "R3"},
		{From: "R1", To: "ground"},
	},
}

var testCases = []struct {
	name    string
	circuit *Circuit
}{
	{
		name:    "TestCircuit1",
		circuit: testCircuit,
	},
	{
		name:    "TestCircuit2",
		circuit: testCircuit2,
	},
	{
		name:    "TestCircuit3",
		circuit: testCircuit3,
	},
}

func TestAssignNodeNumbers(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			nodeNumbers, nodeComponents := assignNodeNumbers(tc.circuit)
			fmt.Printf("%s Node Numbers: %v\n", tc.name, nodeNumbers)
			fmt.Printf("%s Node Components: %v\n", tc.name, nodeComponents)
		})
	}
}

func TestBuildGMatrix(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			nodeNumbers, nodeComponents := assignNodeNumbers(tc.circuit)
			G := buildGMatrix(tc.circuit, nodeNumbers, nodeComponents)
			
			fmt.Printf("%s G Matrix:\n", tc.name)
			if G == nil {
				t.Fatal("G matrix is nil")
			}
			
			r, c := G.Dims()
			fmt.Printf("Matrix dimensions: %d x %d\n", r, c)
			fmt.Printf("%v\n", mat.Formatted(G, mat.Prefix("    "), mat.Squeeze()))
		})
	}
}

func TestBuildMNAMatrices(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			nodeNumbers, nodeComponents := assignNodeNumbers(tc.circuit)
			A, x, z := buildMNAMatrices(tc.circuit, nodeNumbers, nodeComponents)
			
			fmt.Printf("%s A Matrix:\n", tc.name)
			if A == nil {
				t.Fatal("A matrix is nil")
			}
			
			r, c := A.Dims()
			fmt.Printf("Matrix dimensions: %d x %d\n", r, c)
			fmt.Printf("%v\n", mat.Formatted(A, mat.Prefix("    "), mat.Squeeze()))
			
			fmt.Printf("%s x Vector:\n", tc.name)
			fmt.Printf("%v\n", mat.Formatted(x, mat.Prefix("    "), mat.Squeeze()))
			
			fmt.Printf("%s z Vector:\n", tc.name)
			fmt.Printf("%v\n", mat.Formatted(z, mat.Prefix("    "), mat.Squeeze()))
		})
	}
}

func TestBuildBMatrix(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			nodeNumbers, nodeComponents := assignNodeNumbers(tc.circuit)
			B := buildBMatrix(tc.circuit, nodeNumbers, nodeComponents)
			
			fmt.Printf("%s B Matrix:\n", tc.name)
			if B == nil {
				t.Fatal("B matrix is nil")
			}
			
			r, c := B.Dims()
			fmt.Printf("Matrix dimensions: %d x %d\n", r, c)
			fmt.Printf("%v\n", mat.Formatted(B, mat.Prefix("    "), mat.Squeeze()))
		})
	}
}

func TestBuildCMatrix(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			nodeNumbers, nodeComponents := assignNodeNumbers(tc.circuit)
			C := buildCMatrix(tc.circuit, nodeNumbers, nodeComponents)
			
			fmt.Printf("%s C Matrix:\n", tc.name)
			if C == nil {
				t.Fatal("C matrix is nil")
			}
			
			r, c := C.Dims()
			fmt.Printf("Matrix dimensions: %d x %d\n", r, c)
			fmt.Printf("%v\n", mat.Formatted(C, mat.Prefix("    "), mat.Squeeze()))
		})
	}
}

func TestBuildzMatrix(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			nodeNumbers, nodeComponents := assignNodeNumbers(tc.circuit)
			Z := buildzMatrix(tc.circuit, nodeNumbers, nodeComponents)
			
			fmt.Printf("%s Z Matrix:\n", tc.name)
			if Z == nil {
				t.Fatal("Z matrix is nil")
			}
			
			r, _ := Z.Dims()
			fmt.Printf("Matrix dimensions: %d x 1\n", r)
			fmt.Printf("%v\n", mat.Formatted(Z, mat.Prefix("    "), mat.Squeeze()))
		})
	}
}

// Add more test functions for other matrices or operations as needed

func isClose(a, b float64) bool {
	const tolerance = 1e-6
	return math.Abs(a-b) < tolerance
}
