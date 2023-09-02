package main

import (
	"math"
	"math/rand"
)

type Wolf struct {
	x            int
	y            int
	brain        Brain
	energy       int
	angle        float64
	mateCooldown int
	alive        bool
}

const wolfEnergy int = 3600 // How much energy a wolf has when made at the start of the simulation

// Methods
func (wlf *Wolf) update(state *State) bool {
	wlf.energy--
	wlf.mateCooldown--

	// Get the inputs
	wlf.grassPos(state)
	wlf.sheepPos(state)
	wlf.wolfPos(state)
	wlf.grassAngle(state)
	wlf.sheepAngle(state)
	wlf.wolfAngle(state)
	wlf.brain.inputs[3].num = 1

	// Calculate
	wlf.brain.push()
	boolOutput, floatOutput := wlf.brain.output_dump()

	// Create changes
	if boolOutput[0] {
		wlf.bite(state)
	}

	wlf.move(state, floatOutput[1])

	wlf.turn(floatOutput[2])

	if boolOutput[3] {
		wlf.mate(state)
	}

	// Reset
	wlf.brain.set()

	wlf.alive = wlf.energy > 0
	return wlf.alive
}

// Input
func (wlf *Wolf) grassPos(state *State) {
	var tempDis float64
	var x, y int
	var d float64 = state.config.wolfViewDis * state.config.wolfViewDis

	for i := 0; i < len(state.allGrass); i++ {
		x = state.allGrass[i].x - wlf.x
		y = state.allGrass[i].y - wlf.y
		tempDis = math.Pow(float64(x), 2) + math.Pow(float64(y), 2)

		if tempDis < d {
			d = tempDis
		}
	}
	tempDis = math.Sqrt(tempDis)
	tempDis /= state.config.wolfViewDis

	if tempDis == 1 {
		return
	}
	wlf.brain.inputs[0].num = float32(tempDis)
}

func (wlf *Wolf) sheepPos(state *State) {
	var tempDis float64
	var x, y int
	var d float64 = state.config.wolfViewDis * state.config.wolfViewDis

	for i := 0; i < len(state.allSheep); i++ {
		x = state.allSheep[i].x - wlf.x
		y = state.allSheep[i].y - wlf.y
		tempDis = math.Pow(float64(x), 2) + math.Pow(float64(y), 2)

		if tempDis < d {
			d = tempDis
		}
	}
	tempDis = math.Sqrt(tempDis)
	tempDis /= state.config.wolfViewDis

	if tempDis >= 1 {
		return
	}
	wlf.brain.inputs[1].num = float32(tempDis)
}

func (wlf *Wolf) wolfPos(state *State) {
	var tempDis float64
	var x, y int
	var d float64 = state.config.wolfViewDis * state.config.wolfViewDis

	for i := 0; i < len(state.allWolves); i++ {
		x = state.allWolves[i].x - wlf.x
		y = state.allWolves[i].y - wlf.y
		tempDis = math.Pow(float64(x), 2) + math.Pow(float64(y), 2)

		if tempDis < d {
			d = tempDis
		}
	}
	tempDis = math.Sqrt(tempDis)
	tempDis /= state.config.wolfViewDis

	if tempDis >= 1 {
		return
	}
	wlf.brain.inputs[2].num = float32(tempDis)
}

func (wlf *Wolf) grassAngle(state *State) {
	var tempDis float64
	var x, y int

	// Maximum distance
	var d float64 = state.config.wolfViewDis * state.config.wolfViewDis

	for i := 0; i < len(state.allGrass); i++ {
		x = state.allGrass[i].x - wlf.x
		y = state.allGrass[i].y - wlf.y

		// Pythagoras
		tempDis = math.Pow(float64(x), 2) + math.Pow(float64(y), 2)

		// Get the shortest distance
		if tempDis < d {
			d = tempDis
		}
	}

	// Should now be between -1 and 1
	tempDis /= state.config.wolfViewDis
	tempDis = math.Sqrt(tempDis)

	// Range of sight
	if tempDis >= 1 {
		return
	}

	ang := math.Acos(1 / tempDis)
	sign := math.Asin(1 / tempDis)

	ang /= math.Pi
	if sign < 0 {
		ang = -ang
	}

	// Put this input in the brain
	wlf.brain.inputs[3].num = float32(ang)
}

func (wlf *Wolf) sheepAngle(state *State) {
	var tempDis float64
	var x, y int
	var d float64 = state.config.wolfViewDis * state.config.wolfViewDis

	for i := 0; i < len(state.allSheep); i++ {
		x = state.allSheep[i].x - wlf.x
		y = state.allSheep[i].y - wlf.y
		tempDis = math.Pow(float64(x), 2) + math.Pow(float64(y), 2)
		if tempDis < d {
			d = tempDis
		}
	}
	tempDis = math.Sqrt(tempDis)
	tempDis /= state.config.wolfViewDis

	if tempDis >= 1 {
		return
	}
	ang := math.Acos(1 / tempDis)
	sign := math.Asin(1 / tempDis)
	ang /= math.Pi
	if sign < 0 {
		ang = -ang
	}
	wlf.brain.inputs[4].num = float32(ang)
}

func (wlf *Wolf) wolfAngle(state *State) {
	var tempDis float64
	var x, y int
	var d float64 = state.config.wolfViewDis * state.config.wolfViewDis

	for i := 0; i < len(state.allWolves); i++ {
		x = state.allWolves[i].x - wlf.x
		y = state.allWolves[i].y - wlf.y
		tempDis = math.Pow(float64(x), 2) + math.Pow(float64(y), 2)
		if tempDis < d {
			d = tempDis
		}
	}
	tempDis = math.Sqrt(tempDis)
	tempDis /= state.config.wolfViewDis

	if tempDis >= 1 {
		return
	}
	ang := math.Acos(1 / tempDis)
	sign := math.Asin(1 / tempDis)
	ang /= math.Pi
	if sign < 0 {
		ang = -ang
	}
	wlf.brain.inputs[5].num = float32(ang)
}

// Output
func (wlf *Wolf) bite(state *State) {
	wide := wlf.x + state.config.wolfSize
	high := wlf.y + state.config.wolfSize
	for i := 0; i < len(state.allSheep); i++ {
		if wlf.x > state.allSheep[i].x+state.config.sheepSize {
			continue
		}
		if wide < state.allSheep[i].x {
			continue
		}
		if wlf.y > state.allSheep[i].y+state.config.sheepSize {
			continue
		}
		if high < state.allSheep[i].y {
			continue
		}

		wlf.energy += state.allSheep[i].giveEnergy(state)

		break
	}
}

func (wlf *Wolf) move(state *State, dis float32) {
	degAng := wlf.angle * math.Pi / 180
	wlf.x += int(math.Cos(degAng) * state.config.wolfSpeed * float64(dis))
	wlf.y += int(math.Sin(degAng) * state.config.wolfSpeed * float64(dis))
	if wlf.x > 1_000-state.config.wolfSize {
		wlf.x = 1_000 - state.config.wolfSize
	} else if wlf.x < 0 {
		wlf.x = 0
	}
	if wlf.y > 1_000-state.config.wolfSize {
		wlf.y = 1_000 - state.config.wolfSize
	} else if wlf.y < 0 {
		wlf.y = 0
	}
}

func (wlf *Wolf) turn(ang float32) {
	wlf.angle += float64(ang) * 5
}

func (wlf *Wolf) mate(state *State) {
	if len(state.allWolves) >= state.config.wolfMaxAmt {
		return
	}

	// Restrictions
	if !wlf.canMate(state) {
		return
	}

	foundPartner := false
	var partner *Wolf

	wide := wlf.x + state.config.wolfSize
	high := wlf.y + state.config.wolfSize
	for i := 0; i < len(state.allWolves); i++ {
		if wlf.x > state.allWolves[i].x+state.config.wolfSize {
			continue
		}
		if wide < state.allWolves[i].x {
			continue
		}
		if wlf.y > state.allWolves[i].y+state.config.wolfSize {
			continue
		}
		if high < state.allWolves[i].y {
			continue
		}

		partner = state.allWolves[i]

		// Partner restrictions
		if !partner.canMate(state) {
			continue
		}

		foundPartner = true

		break
	}

	if !foundPartner {
		return
	}
	if partner == wlf {
		return
	}

	wlf.energy -= state.config.wolfMateLoss
	partner.energy -= state.config.wolfMateLoss
	wlf.mateCooldown = state.config.wolfMatePartnerCooldown
	partner.mateCooldown = state.config.wolfMatePartnerCooldown

	child := &Wolf{
		wlf.x + rand.Intn(20) - 10,
		wlf.y + rand.Intn(20) - 10,
		createBlankBrain(),
		state.config.wolfChildEnergy,
		0,
		state.config.wolfMateChildCooldown,
		true,
	}

	for i := 0; i < len(child.brain.inputs); i++ {
		for j := 0; j < len(child.brain.inputs[i].links); j++ {
			if rand.Intn(2) == 0 {
				child.brain.inputs[i].weights = append(child.brain.inputs[i].weights, wlf.brain.inputs[i].weights[j])
			} else {
				child.brain.inputs[i].weights = append(child.brain.inputs[i].weights, partner.brain.inputs[i].weights[j])
			}

			if rand.Float32() < state.config.wolfRandWeight {
				child.brain.inputs[i].weights[j] = randWeight()
			}
		}
	}

	for i := 0; i < len(child.brain.layers); i++ {
		for j := 0; j < len(child.brain.layers[i].nodes); j++ {
			if child.brain.layers[i].nodes[j].lastLayer {
				for k := 0; k < len(child.brain.layers[i].nodes[j].linksO); k++ {
					if rand.Intn(2) == 0 {
						child.brain.layers[i].nodes[j].weights = append(child.brain.layers[i].nodes[j].weights, wlf.brain.layers[i].nodes[j].weights[k])
					} else {
						child.brain.layers[i].nodes[j].weights = append(child.brain.layers[i].nodes[j].weights, partner.brain.layers[i].nodes[j].weights[k])
					}

					if rand.Float32() < state.config.wolfRandWeight {
						child.brain.layers[i].nodes[j].weights[k] = randWeight()
					}
				}
			} else {
				for k := 0; k < len(child.brain.layers[i].nodes[j].linksN); k++ {
					if rand.Intn(2) == 0 {
						child.brain.layers[i].nodes[j].weights = append(child.brain.layers[i].nodes[j].weights, wlf.brain.layers[i].nodes[j].weights[k])
					} else {
						child.brain.layers[i].nodes[j].weights = append(child.brain.layers[i].nodes[j].weights, partner.brain.layers[i].nodes[j].weights[k])
					}

					if rand.Float32() < state.config.wolfRandWeight {
						child.brain.layers[i].nodes[j].weights[k] = randWeight()
					}
				}
			}
		}
	}

	state.allWolves = append(state.allWolves, child)
}

func (wlf *Wolf) canMate(state *State) bool {
	if wlf.energy < state.config.wolfMateBarrier {
		return false
	}
	if wlf.mateCooldown > 0 {
		return false
	}
	return true
}
