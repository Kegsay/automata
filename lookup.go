package automata

// LookupTable stores mappings of:
//  - Neuron IDs to Neurons
//  - Connection IDs to Connections
// It can be thought of as a global hash map, but is implemented slightly differently for
// performance reasons.
//
// Rationale: Neurons need references to other neurons/connections (e.g. neighbours). The simplest
// way to do this is to store a map[ID]*Thing in the Neuron struct itself. This ends up being slow because
// this is called a lot and each time incurs hashing overheads. It would be better if this was done
// as a slice, especially since network topologies don't tend to change at runtime so there are no
// resizing overheads. This means the IDs would be indexes. LookupTable is a massive global
// slice which provides Neurons/Layers/Connections access to each other via their ID. It's not a true
// global variable as it is dependency injected at the point of use, allowing the ability of running
// multiple disconnected networks without sharing the same ID space.
type LookupTable struct {
	Neurons     []*Neuron
	Connections []*Connection
}

var GlobalLookupTable = &LookupTable{} // TODO: Dependency inject this

// NewLookupTable creates a new LookupTable. There are no existing mappings initially.
func NewLookupTable() *LookupTable {
	return &LookupTable{}
}

func (t *LookupTable) SetNeuron(neuron *Neuron) NeuronID {
	t.Neurons = append(t.Neurons, neuron)
	return NeuronID(len(t.Neurons) - 1)
}

func (t *LookupTable) SetNeuronWithID(id NeuronID, neuron *Neuron) {
	if int(id) > (len(t.Neurons) - 1) {
		// pad out the slice
		diff := int(id) - (len(t.Neurons) - 1)
		t.Neurons = append(t.Neurons, make([]*Neuron, diff)...)
	}
	t.Neurons[id] = neuron
}

func (t *LookupTable) GetNeuron(id NeuronID) *Neuron {
	if int(id) > (len(t.Neurons) - 1) {
		return nil
	}
	return t.Neurons[id]
}

func (t *LookupTable) SetConnection(conn *Connection) ConnID {
	t.Connections = append(t.Connections, conn)
	return ConnID(len(t.Connections) - 1)
}

func (t *LookupTable) SetConnectionWithID(id ConnID, conn *Connection) {
	if int(id) > (len(t.Connections) - 1) {
		// pad out the slice
		diff := int(id) - (len(t.Connections) - 1)
		t.Connections = append(t.Connections, make([]*Connection, diff)...)
	}
	t.Connections[id] = conn
}

func (t *LookupTable) GetConnection(id ConnID) *Connection {
	if int(id) > (len(t.Connections) - 1) {
		return nil
	}
	return t.Connections[id]
}