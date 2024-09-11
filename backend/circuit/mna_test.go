package circuit

import (
	"math"
	"testing"
	"fmt"
	"gonum.org/v1/gonum/mat"
	// "strconv"
	// "strings"
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
		{From: "ground", To: "V1", Polarity: "-"},
		{From: "V1", To: "n1", Polarity: "+"},
		{From: "n1", To: "R1", Polarity: ""},
		{From: "R1", To: "n2", Polarity: ""},
		{From: "n2", To: "R2", Polarity: ""},
		{From: "R2", To: "n3", Polarity: ""},
		{From: "n3", To: "V2", Polarity: "+"},
		{From: "V2", To: "ground", Polarity: "-"},
		{From: "n2", To: "R3", Polarity: ""},
		{From: "R3", To: "ground", Polarity: ""},
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
		{From: "ground", To: "R1", Polarity: ""},
		{From: "R1", To: "n1", Polarity: ""},
		{From: "n1", To: "V1", Polarity: "-"},
		{From: "V1", To: "n2", Polarity: "+"},
		{From: "n2", To: "R2", Polarity: ""},
		{From: "R2", To: "n3", Polarity: ""},
		{From: "n3", To: "V2", Polarity: "+"},
		{From: "V2", To: "ground", Polarity: "-"},
		{From: "n2", To: "R3", Polarity: ""},
		{From: "R3", To: "ground", Polarity: ""},
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
		{From: "ground", To: "I1", Polarity: "-"},
		{From: "I1", To: "n1", Polarity: "+"},
		{From: "n1", To: "V1", Polarity: "-"},
		{From: "V1", To: "n2", Polarity: "+"},
		{From: "n2", To: "R1", Polarity: ""},
		{From: "R1", To: "ground", Polarity: ""},
		{From: "n2", To: "R2", Polarity: ""},
		{From: "R2", To: "ground", Polarity: ""},
		{From: "n2", To: "R3", Polarity: ""},
		{From: "R3", To: "ground", Polarity: ""},
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
	expectedB := map[string]*mat.Dense{
		"TestCircuit1": mat.NewDense(3, 2, []float64{
			1,  0,
			0,  0,
			0,  1,
		}),
		"TestCircuit2": mat.NewDense(3, 2, []float64{
			-1, 0,
			 1, 0,
			 0, 1,
		}),
		"TestCircuit3": mat.NewDense(2, 1, []float64{
			-1,
			 1,
		}),
	}

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

			expected := expectedB[tc.name]
			if !mat.EqualApprox(B, expected, 1e-6) {
				t.Errorf("B matrix for %s does not match expected.\nGot:\n%v\nExpected:\n%v",
					tc.name, mat.Formatted(B), mat.Formatted(expected))
			}
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
	expectedZ := map[string]*mat.VecDense{
		"TestCircuit1": mat.NewVecDense(5, []float64{0, 0, 0, 32, 20}),
		"TestCircuit2": mat.NewVecDense(5, []float64{0, 0, 0, 32, 20}),
		"TestCircuit3": mat.NewVecDense(3, []float64{10, 0, 32}),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			nodeNumbers, nodeComponents := assignNodeNumbers(tc.circuit)
			z := buildzMatrix(tc.circuit, nodeNumbers, nodeComponents)
			
			fmt.Printf("%s z Matrix:\n", tc.name)
			if z == nil {
				t.Fatal("z matrix is nil")
			}
			
			r, _ := z.Dims()
			fmt.Printf("Matrix dimensions: %d x 1\n", r)
			fmt.Printf("%v\n", mat.Formatted(z, mat.Prefix("    "), mat.Squeeze()))

			expected := expectedZ[tc.name]
			if !mat.EqualApprox(z, expected, 1e-6) {
				t.Errorf("z matrix for %s does not match expected.\nGot:\n%v\nExpected:\n%v",
					tc.name, mat.Formatted(z), mat.Formatted(expected))
			}
		})
	}
}

func TestBuildxMatrix(t *testing.T) {
	expectedX := map[string]*mat.VecDense{
		"TestCircuit1": mat.NewVecDense(5, []float64{0, 0, 0, 0, 0}),
		"TestCircuit2": mat.NewVecDense(5, []float64{0, 0, 0, 0, 0}),
		"TestCircuit3": mat.NewVecDense(3, []float64{0, 0, 0}),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			nodeNumbers, _ := assignNodeNumbers(tc.circuit)
			x := buildxMatrix(tc.circuit, nodeNumbers)
			
			fmt.Printf("%s x Matrix:\n", tc.name)
			if x == nil {
				t.Fatal("x matrix is nil")
			}
			
			r, _ := x.Dims()
			fmt.Printf("Matrix dimensions: %d x 1\n", r)
			fmt.Printf("%v\n", mat.Formatted(x, mat.Prefix("    "), mat.Squeeze()))

			expected := expectedX[tc.name]
			if !mat.EqualApprox(x, expected, 1e-6) {
				t.Errorf("x matrix for %s does not match expected.\nGot:\n%v\nExpected:\n%v",
					tc.name, mat.Formatted(x), mat.Formatted(expected))
			}
		})
	}
}

func TestSolveCircuit(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			results, err := SolveCircuit(tc.circuit)
			if err != nil {
				t.Fatalf("Error solving circuit: %v", err)
			}

			fmt.Printf("%s Solution:\n", tc.name)
			for key, value := range results {
				fmt.Printf("%s: %.6f\n", key, value)
			}
			fmt.Println() // Add a blank line between test cases for readability
		})
	}
}

// Add more test functions for other matrices or operations as needed

func isClose(a, b float64) bool {
	const tolerance = 1e-6
	return math.Abs(a-b) < tolerance
}
