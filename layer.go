package automata

import "fmt"

// Layer represents a group of neurons which activate together.
type Layer struct {
	Size        int
	List        []Neuron
	ConnectedTo []LayerConnection
}

func NewLayer(size int) Layer {
	neurons := make([]Neuron, size)
	for i := 0; i < size; i++ {
		neurons[i] = NewNeuron()
	}
	return Layer{
		Size: size,
		List: neurons,
	}
}

// Activate all neurons in the layer.
func (l *Layer) Activate(neurons []Neuron) ([]float64, error) {
	var activations []float64

	// Activate without an input
	if neurons == nil {
		for i := 0; i < len(l.List); i++ {
			activation := l.List[i].Activate(nil)
			activations = append(activations, activation)
		}
	} else if len(neurons) != len(l.List) {
		return nil, fmt.Errorf("input and layer size mismatch: cannot activate")
	} else { // Activate with input
		for i := 0; i < len(l.List); i++ {
			activation := l.List[i].Activate(&neurons[i])
			activations = append(activations, activation)
		}
	}

	return activations, nil
}

// Propagate an error on all neurons in this layer.
func (l *Layer) Propagate() {}

// Project a connection from this layer to another one.
func (l *Layer) Project(toLayer Layer, ltype LayerType, weights []float64) *LayerConnection {
	if l.isConnected(toLayer) {
		return nil
	}
	lc := NewLayerConnection(*l, toLayer, ltype, weights)
	return &lc
}

// Gate a connection between two layers.
func (l *Layer) Gate() {}

// isConnected returns true if this layer is connected to the target layer already.
func (l *Layer) isConnected(targetLayer Layer) bool {
	return false // TODO
}