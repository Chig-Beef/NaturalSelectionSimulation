package main

type Layer struct {
	nodes []Node
}

// Reset to default
func (layer *Layer) set() {
	for i := 0; i < len(layer.nodes); i++ {
		layer.nodes[i].set()
	}
}

// Calculate output
func (layer *Layer) push() {
	for i := 0; i < len(layer.nodes); i++ {
		layer.nodes[i].push()
	}
}
