package main

type Output struct {
	num   float32
	fired bool
	infs  int // Influences
}

func (output *Output) calc() (bool, float32) {
	if output.infs != 0 {
		output.num /= float32(output.infs)
	}
	if output.num > 0 {
		output.fired = true
	}
	return output.fired, output.num
}

func (output *Output) set() {
	output.num = 0
	output.infs = 0
	output.fired = false
}
