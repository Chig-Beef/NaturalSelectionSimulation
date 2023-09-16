package main

import (
	"fmt"
	"math"
)

type Output struct {
	num   float32
	fired bool
	infs  int // Influences
}

// Calculate output
func (output *Output) calc() (bool, float32) {
	// Avoids division by 0
	if output.infs != 0 {
		output.num /= float32(output.infs)
	}
	if output.num > 0 {
		output.fired = true
	}
	if math.IsNaN(float64(output.num)) {
		fmt.Println("Node Issue")
	}
	return output.fired, output.num
}

// Reset to default
func (output *Output) set() {
	output.num = 0
	output.infs = 0
	output.fired = false
}
