package circuit

type ComponentType string

const (
    Battery  ComponentType = "battery"
    Resistor ComponentType = "resistor"
    // Add more component types as needed
)

type Component struct {
    ID    string
    Type  ComponentType
    Value float64
}

type Connection struct {
    From string
    To   string
}

type Circuit struct {
    Components  []Component
    Connections []Connection
}