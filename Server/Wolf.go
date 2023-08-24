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
}

const wolfEnergy int = 3600                     // How much energy a wolf has when made at the start of the simulation
const wolfSpeed float64 = 5                     // How fast the wolf is
const wolfSize int = 25                         // How big a wolf is
const wolfRandWeight float32 = float32(1) / 100 // The chance that a new weight is random rather than inherited
const wolfMateBarrier int = 1800                // How much energy a wolf needs to breed
const wolfMateLoss int = 900                    // How much energy a wolf loses from breeding
const wolfMatePartnerCooldown int = 240         // How long from giving birth a wolf can attempt mating again
const wolfMateChildCooldown int = 360           // How long from being born a wolf can attempt to mate
const wolfChildEnergy int = 2 * wolfMateLoss    // How much energy a wolf has when it is born
const wolfMaxAmt int = 300                      // Maximum amount of wolves
const wolfViewDis float64 = 100                 // How far a wolf can see

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

	wlf.move(floatOutput[1])

	wlf.turn(floatOutput[2])

	if boolOutput[3] {
		wlf.mate(state)
	}

	// Reset
	wlf.brain.set()

	return wlf.energy > 0
}

// Input
func (wlf *Wolf) grassPos(state *State) {
	var tempDis float64
	var x, y int
	var d float64 = wolfViewDis

	for i := 0; i < len(state.allGrass); i++ {
		x = state.allGrass[i].x - wlf.x
		y = state.allGrass[i].y - wlf.y
		tempDis = math.Sqrt(math.Pow(float64(x), 2) + math.Pow(float64(y), 2))

		if tempDis < d {
			d = tempDis
		}
	}
	tempDis /= wolfViewDis

	if tempDis >= 1 || tempDis <= -1 {
		return
	}
	wlf.brain.inputs[0].num = float32(tempDis)
}

func (wlf *Wolf) sheepPos(state *State) {
	var tempDis float64
	var x, y int
	var d float64 = wolfViewDis

	for i := 0; i < len(state.allSheep); i++ {
		x = state.allSheep[i].x - wlf.x
		y = state.allSheep[i].y - wlf.y
		tempDis = math.Sqrt(math.Pow(float64(x), 2) + math.Pow(float64(y), 2))

		if tempDis < d {
			d = tempDis
		}
	}
	tempDis /= wolfViewDis

	if tempDis >= 1 || tempDis <= -1 {
		return
	}
	wlf.brain.inputs[1].num = float32(tempDis)
}

func (wlf *Wolf) wolfPos(state *State) {
	var tempDis float64
	var x, y int
	var d float64 = wolfViewDis

	for i := 0; i < len(state.allWolves); i++ {
		x = state.allWolves[i].x - wlf.x
		y = state.allWolves[i].y - wlf.y
		tempDis = math.Sqrt(math.Pow(float64(x), 2) + math.Pow(float64(y), 2))

		if tempDis < d {
			d = tempDis
		}
	}
	tempDis /= wolfViewDis

	if tempDis >= 1 || tempDis <= -1 {
		return
	}
	wlf.brain.inputs[2].num = float32(tempDis)
}

func (wlf *Wolf) grassAngle(state *State) {
	var tempDis float64
	var x, y int

	// Maximum distance
	var d float64 = wolfViewDis

	for i := 0; i < len(state.allGrass); i++ {
		x = state.allGrass[i].x - wlf.x
		y = state.allGrass[i].y - wlf.y

		// Pythagoras
		tempDis = math.Sqrt(math.Pow(float64(x), 2) + math.Pow(float64(y), 2))

		// Get the shortest distance
		if tempDis < d {
			d = tempDis
		}
	}

	// Should now be between -1 and 1
	tempDis /= sheepViewDis

	// Range of sight
	if tempDis >= 1 || tempDis <= -1 {
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
	var d float64 = wolfViewDis

	for i := 0; i < len(state.allSheep); i++ {
		x = state.allSheep[i].x - wlf.x
		y = state.allSheep[i].y - wlf.y
		tempDis = math.Sqrt(math.Pow(float64(x), 2) + math.Pow(float64(y), 2))
		if tempDis < d {
			d = tempDis
		}
	}
	tempDis /= sheepViewDis

	if tempDis >= 1 || tempDis <= -1 {
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
	var d float64 = sheepViewDis

	for i := 0; i < len(state.allWolves); i++ {
		x = state.allWolves[i].x - wlf.x
		y = state.allWolves[i].y - wlf.y
		tempDis = math.Sqrt(math.Pow(float64(x), 2) + math.Pow(float64(y), 2))
		if tempDis < d {
			d = tempDis
		}
	}
	tempDis /= sheepViewDis

	if tempDis >= 1 || tempDis <= -1 {
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
	for i := 0; i < len(state.allSheep); i++ {
		if wlf.x > state.allSheep[i].x+sheepSize {
			continue
		}
		if wlf.x+wolfSize < state.allSheep[i].x {
			continue
		}
		if wlf.y > state.allSheep[i].y+sheepSize {
			continue
		}
		if wlf.y+wolfSize < state.allSheep[i].y {
			continue
		}

		wlf.energy += state.allSheep[i].giveEnergy()

		break
	}
}

func (wlf *Wolf) move(dis float32) {
	wlf.x += int(math.Cos(wlf.angle) * wolfSpeed * float64(dis))
	wlf.y += int(math.Sin(wlf.angle) * wolfSpeed * float64(dis))
	if wlf.x > 1_000 {
		wlf.x = 1_000
	} else if wlf.x < 0 {
		wlf.x = 0
	}
	if wlf.y > 1_000 {
		wlf.y = 1_000
	} else if wlf.y < 0 {
		wlf.y = 0
	}
}

func (wlf *Wolf) turn(ang float32) {
	wlf.angle += float64(ang)
}

func (wlf *Wolf) mate(state *State) {
	if len(state.allWolves) >= wolfMaxAmt {
		return
	}

	// Restrictions
	if !wlf.canMate() {
		return
	}

	foundPartner := false
	var partner *Wolf

	for i := 0; i < len(state.allWolves); i++ {
		if wlf.x > state.allWolves[i].x+wolfSize {
			continue
		}
		if wlf.x+wolfSize < state.allWolves[i].x {
			continue
		}
		if wlf.y > state.allWolves[i].y+wolfSize {
			continue
		}
		if wlf.y+wolfSize < state.allWolves[i].y {
			continue
		}

		partner = state.allWolves[i]

		// Partner restrictions
		if !partner.canMate() {
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

	wlf.energy -= wolfMateLoss
	partner.energy -= wolfMateLoss
	wlf.mateCooldown = wolfMatePartnerCooldown
	partner.mateCooldown = wolfMatePartnerCooldown

	child := &Wolf{
		wlf.x + rand.Intn(20) - 10,
		wlf.y + rand.Intn(20) - 10,
		createBlankBrain(),
		wolfChildEnergy,
		0,
		wolfMateChildCooldown,
	}

	for i := 0; i < len(child.brain.inputs); i++ {
		for j := 0; j < len(child.brain.inputs[i].links); j++ {
			if rand.Intn(2) == 0 {
				child.brain.inputs[i].weights = append(child.brain.inputs[i].weights, wlf.brain.inputs[i].weights[j])
			} else {
				child.brain.inputs[i].weights = append(child.brain.inputs[i].weights, partner.brain.inputs[i].weights[j])
			}

			if rand.Float32() < wolfRandWeight {
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

					if rand.Float32() < wolfRandWeight {
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

					if rand.Float32() < wolfRandWeight {
						child.brain.layers[i].nodes[j].weights[k] = randWeight()
					}
				}
			}
		}
	}

	state.allWolves = append(state.allWolves, child)
}

func (wlf *Wolf) canMate() bool {
	if wlf.energy < wolfMateBarrier {
		return false
	}
	if wlf.mateCooldown > 0 {
		return false
	}
	return true
}
