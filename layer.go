package automata

import "fmt"

// Layer represents a group of neurons which activate together.
type Layer struct {
	List        []*Neuron
	ConnectedTo []LayerConnection
	LookupTable *LookupTable
}

func NewLayer(table *LookupTable, size int) Layer {
	neurons := make([]*Neuron, size)
	for i := 0; i < size; i++ {
		neurons[i] = NewNeuron(table)
	}
	return Layer{
		List:        neurons,
		LookupTable: table,
	}
}

// Activate all neurons in the layer.
func (l *Layer) Activate(inputs []float64) ([]float64, error) {
	var activations []float64

	// Activate without an input
	if inputs == nil {
		for i := 0; i < len(l.List); i++ {
			activation := l.List[i].Activate(nil)
			activations = append(activations, activation)
		}
	} else if len(inputs) != len(l.List) {
		return nil, fmt.Errorf("input and layer size mismatch: cannot activate")
	} else { // Activate with input
		for i := 0; i < len(l.List); i++ {
			activation := l.List[i].Activate(&inputs[i])
			activations = append(activations, activation)
		}
	}

	return activations, nil
}

// Propagate an error on all neurons in this layer.
func (l *Layer) Propagate(rate float64, target []float64) error {
	if target != nil {
		if len(target) != len(l.List) {
			return fmt.Errorf("target and layer size mismatch: cannot propagate")
		}
		for i := len(l.List) - 1; i >= 0; i-- {
			l.List[i].Propagate(rate, &target[i])
		}
	} else {
		for i := len(l.List) - 1; i >= 0; i-- {
			l.List[i].Propagate(rate, nil)
		}
	}
	return nil
}

// Project a connection from this layer to another one.
func (l *Layer) Project(toLayer *Layer, ltype LayerType) *LayerConnection {
	if l.isConnected(toLayer) {
		return nil
	}
	lc := NewLayerConnection(l, toLayer, ltype)
	return &lc
}

// Gate a connection between two layers.
func (l *Layer) Gate(conn *LayerConnection, gateType GateType) error {
	switch gateType {
	case GateTypeInput:
		if len(conn.To.List) != len(l.List) {
			return fmt.Errorf("Cannot gate connection to layer - neuron count mismatch, %d != %d", len(conn.To.List), len(l.List))
		}
		for i, neuron := range conn.To.List {
			gater := l.List[i]
			for _, gatedID := range neuron.Inputs {
				gated := l.LookupTable.GetConnection(gatedID)
				if _, ok := conn.Connections[gated.ID]; ok {
					gater.Gate(gated)
				}
			}
		}
	case GateTypeOutput:
		if len(conn.From.List) != len(l.List) {
			return fmt.Errorf("Cannot gate connection to layer - neuron count mismatch, %d != %d", len(conn.From.List), len(l.List))
		}
		for i, neuron := range conn.From.List {
			gater := l.List[i]
			for _, gatedID := range neuron.Projected {
				gated := l.LookupTable.GetConnection(gatedID)
				if _, ok := conn.Connections[gated.ID]; ok {
					gater.Gate(gated)
				}
			}
		}
	case GateTypeOneToOne:
		if len(conn.List) != len(l.List) {
			return fmt.Errorf("Cannot gate connection to layer - neuron count mismatch, %d != %d", len(conn.List), len(l.List))
		}
		for i := range l.List {
			gater := l.List[i]
			gated := conn.List[i]
			gater.Gate(gated)
		}
	default:
		return fmt.Errorf("unknown GateType: %d", gateType)
	}
	return nil
}

// SetBias sets the bias of all neurons in this layer to the given value.
func (l *Layer) SetBias(bias float64) {
	for i := range l.List {
		l.List[i].Bias = bias
	}
}

// isConnected returns true if this layer is connected to the target layer already.
func (l *Layer) isConnected(targetLayer *Layer) bool {
	connCount := 0
	for _, neuron := range l.List {
		for _, target := range targetLayer.List {
			if neuron.ConnectionForNeuron(target) != nil {
				connCount++
			}
		}
	}
	if connCount == (len(l.List) * len(targetLayer.List)) {
		return true // all to all layer
	}
	return false
}
