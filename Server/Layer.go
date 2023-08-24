package main

type Layer struct {
	nodes []Node
}

func (layer *Layer) set() {
	for i := 0; i < len(layer.nodes); i++ {
		layer.nodes[i].set()
	}
}
