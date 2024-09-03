package circuit

func SimulateCircuit(components []Component, connections []Connection) (map[string]float64, error) {
    c := &Circuit{
        Components:  components,
        Connections: connections,
    }
    return SolveCircuit(c)
}