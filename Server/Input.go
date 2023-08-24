package main

type Input struct {
	num     float32
	links   []*Node
	weights []float32
}

func (input *Input) push() {
	for i := 0; i < len(input.links); i++ {
		input.links[i].num += input.num * input.weights[i]
		input.links[i].infs++
	}
}

func (input *Input) set() {
	input.num = 0
}
