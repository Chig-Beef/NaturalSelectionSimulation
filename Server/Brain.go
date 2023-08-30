package main

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
)

type Brain struct {
	inputs  []Input
	layers  []Layer
	outputs []Output
}

const nLayers int = 7  // How many layers
const nDepth int = 8   // How many nodes per layer
const nInputs int = 7  // How many inputs
const nOutputs int = 4 // How many outputs

func createBlankBrain() Brain {
	brain := Brain{}

	// Creating all the nodes, inputs, outputs
	for i := 0; i < nInputs; i++ {
		brain.inputs = append(brain.inputs, Input{})
	}

	for i := 0; i < nLayers; i++ {
		brain.layers = append(brain.layers, Layer{})
		for j := 0; j < nDepth; j++ {
			brain.layers[i].nodes = append(brain.layers[i].nodes, Node{})
			brain.layers[i].nodes[j].lastLayer = false
		}
	}

	for i := 0; i < nOutputs; i++ {
		brain.outputs = append(brain.outputs, Output{})
	}

	// Creating the links between these parts
	for i := 0; i < nInputs; i++ {
		for j := 0; j < nDepth; j++ {
			brain.inputs[i].links = append(brain.inputs[i].links, &brain.layers[0].nodes[j])
		}
	}

	for i := 0; i < nLayers-1; i++ {
		for j := 0; j < nDepth; j++ {
			for k := 0; k < nDepth; k++ {
				brain.layers[i].nodes[j].linksN = append(brain.layers[i].nodes[j].linksN, &brain.layers[i+1].nodes[k])
			}
		}
	}

	for i := 0; i < nDepth; i++ {
		brain.layers[len(brain.layers)-1].nodes[i].lastLayer = true
		for j := 0; j < nOutputs; j++ {
			brain.layers[len(brain.layers)-1].nodes[i].linksO = append(brain.layers[len(brain.layers)-1].nodes[i].linksO, &brain.outputs[j])
		}
	}

	return brain
}

func createRandomConnections(brain Brain) Brain {
	for i := 0; i < len(brain.inputs); i++ {
		for j := 0; j < len(brain.inputs[i].links); j++ {
			brain.inputs[i].weights = append(brain.inputs[i].weights, randWeight())
		}
	}
	var top int
	for i := 0; i < len(brain.layers); i++ {
		for j := 0; j < len(brain.layers[i].nodes); j++ {
			if !brain.layers[i].nodes[j].lastLayer {
				top = len(brain.layers[i].nodes[j].linksN)
			} else {
				top = len(brain.layers[i].nodes[j].linksO)
			}

			for k := 0; k < top; k++ {
				brain.layers[i].nodes[j].weights = append(brain.layers[i].nodes[j].weights, randWeight())
			}
		}
	}

	// Outputs don't have links

	return brain
}

func (brain *Brain) push() {
	// Send data through the neural network

	for i := 0; i < len(brain.inputs); i++ {
		brain.inputs[i].push()
	}

	for i := 0; i < len(brain.layers); i++ {
		for j := 0; j < len(brain.layers[i].nodes); j++ {
			brain.layers[i].nodes[j].push()
		}
	}
}

func (brain *Brain) output_dump() ([]bool, []float32) {
	boolOutput := []bool{}
	floatOutput := []float32{}

	for i := 0; i < len(brain.outputs); i++ {
		b, f := brain.outputs[i].calc()
		boolOutput = append(boolOutput, b)
		floatOutput = append(floatOutput, f)
	}
	return boolOutput, floatOutput
}

func (brain *Brain) set() {
	for i := 0; i < len(brain.inputs); i++ {
		brain.inputs[i].set()
	}

	for i := 0; i < len(brain.layers); i++ {
		brain.layers[i].set()
	}

	for i := 0; i < len(brain.outputs); i++ {
		brain.outputs[i].set()
	}
}

func randWeight() float32 {
	// Creates a random float between -1 and 1
	return rand.Float32()*2 - 1
}

func (brain Brain) convToStr() string {
	// Converts a brain into string format for saving

	outputString := ""

	// Inputs
	for _, input := range brain.inputs {
		for _, weight := range input.weights {
			outputString += strconv.FormatFloat(float64(weight), 'f', -1, 32)
			outputString += "_"
		}
		outputString = outputString[:len(outputString)-1] + "+"
	}

	outputString = outputString[:len(outputString)-1] + "="

	// Nodes
	for _, layer := range brain.layers {
		for _, node := range layer.nodes {
			for _, weight := range node.weights {
				outputString += strconv.FormatFloat(float64(weight), 'f', -1, 32)
				outputString += "_"
			}
			outputString = outputString[:len(outputString)-1] + "+"
		}
		outputString = outputString[:len(outputString)-1] + "!"
	}

	outputString = outputString[:len(outputString)-1] + "=" + strconv.Itoa(len(brain.outputs))

	return outputString
}

func convBrainFromStr(data string) (Brain, error) {
	// Create a Brain from string

	brain := Brain{}
	splitData := strings.Split(data, "=")

	if len(splitData) != 3 {
		return brain, errors.New("expected inputs, layers, outputs, did not get the right amount to hold these 3 values")
	}

	inputs := strings.Split(splitData[0], "+")
	for i, input := range inputs {
		obj := strings.Split(input, "_")
		brain.inputs = append(brain.inputs, Input{})

		for _, weight := range obj {
			temp, err := strconv.ParseFloat(weight, 32)
			if err != nil {
				return brain, errors.New("a weight in a brain wasn't in the correct format")
			}

			brain.inputs[i].weights = append(brain.inputs[i].weights, float32(temp))
		}
	}

	layers := strings.Split(splitData[1], "!")
	for i, layer := range layers {
		objL := strings.Split(layer, "+")
		brain.layers = append(brain.layers, Layer{})

		for j, node := range objL {
			objN := strings.Split(node, "_")
			brain.layers[i].nodes = append(brain.layers[i].nodes, Node{})

			for _, weight := range objN {
				temp, err := strconv.ParseFloat(weight, 32)
				if err != nil {
					return brain, errors.New("a weight in a brain wasn't in the correct format")
				}

				brain.layers[i].nodes[j].weights = append(brain.layers[i].nodes[j].weights, float32(temp))
			}
		}
	}

	outputs, err := strconv.Atoi(splitData[2])
	if err != nil {
		return brain, errors.New("the amount of outputs was not a valid number")
	}
	for i := 0; i < outputs; i++ {
		brain.outputs = append(brain.outputs, Output{})
	}

	return brain, nil
}
