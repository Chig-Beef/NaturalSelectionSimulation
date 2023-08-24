package main

import (
	"math"
	"math/rand"
)

type Sheep struct {
	x            int
	y            int
	brain        Brain
	energy       int
	angle        float64
	mateCooldown int
}

// Look between this file and the wolf file for explanations, as the code is very similar.
const sheepEnergy int = 3600
const sheepSpeed float64 = 5
const sheepSize int = 25
const sheepRandWeight float32 = 1 / 100
const sheepMateBarrier int = 1800
const sheepMateLoss int = 900
const sheepMatePartnerCooldown int = 240
const sheepMateChildCooldown int = 360
const sheepEnergyLoss int = 100
const sheepEnergyGive int = 80
const sheepChildEnergy int = 2 * sheepMateLoss
const sheepMaxAmt int = 300
const sheepViewDis float64 = 100

// Methods
func (shp *Sheep) update(state *State) bool {
	shp.energy--
	shp.mateCooldown--

	// Get the inputs
	shp.grassPos(state)
	shp.sheepPos(state)
	shp.wolfPos(state)
	shp.grassAngle(state)
	shp.sheepAngle(state)
	shp.wolfAngle(state)
	shp.brain.inputs[6].num = 1 // constant

	// Calculate
	shp.brain.push()
	boolOutput, floatOutput := shp.brain.output_dump()

	// Create changes
	if boolOutput[0] {
		shp.bite(state)
	}
	shp.move(floatOutput[1])
	shp.turn(floatOutput[2])
	if boolOutput[3] {
		shp.mate(state)
	}

	// Reset
	shp.brain.set()

	return shp.energy > 0
}

// Input
func (shp *Sheep) grassPos(state *State) {
	var tempDis float64
	var x, y int

	// Maximum distance
	var d float64 = sheepViewDis

	for i := 0; i < len(state.allGrass); i++ {
		x = state.allGrass[i].x - shp.x
		y = state.allGrass[i].y - shp.y

		// Pythagoras
		tempDis = math.Sqrt(math.Pow(float64(x), 2) + math.Pow(float64(y), 2))

		// Get the shortest distance
		if tempDis < d {
			d = tempDis
		}
	}

	// Should now be between -1 and 1
	tempDis /= sheepViewDis

	// Range of sight (the sheep can only see so far)
	if tempDis >= 1 || tempDis <= -1 {
		return
	}

	// Put this input in the brain
	shp.brain.inputs[0].num = float32(tempDis)
}

func (shp *Sheep) sheepPos(state *State) {
	var tempDis float64
	var x, y int
	var d float64 = sheepViewDis

	for i := 0; i < len(state.allSheep); i++ {
		x = state.allSheep[i].x - shp.x
		y = state.allSheep[i].y - shp.y
		tempDis = math.Sqrt(math.Pow(float64(x), 2) + math.Pow(float64(y), 2))
		if tempDis < d {
			d = tempDis
		}
	}
	tempDis /= sheepViewDis
	if tempDis >= 1 || tempDis <= -1 {
		return
	}
	shp.brain.inputs[1].num = float32(tempDis)
}

func (shp *Sheep) wolfPos(state *State) {
	var tempDis float64
	var x, y int
	var d float64 = sheepViewDis

	for i := 0; i < len(state.allWolves); i++ {
		x = state.allWolves[i].x - shp.x
		y = state.allWolves[i].y - shp.y
		tempDis = math.Sqrt(math.Pow(float64(x), 2) + math.Pow(float64(y), 2))
		if tempDis < d {
			d = tempDis
		}
	}
	tempDis /= sheepViewDis
	if tempDis >= 1 || tempDis <= -1 {
		return
	}
	shp.brain.inputs[2].num = float32(tempDis)
}

func (shp *Sheep) grassAngle(state *State) {
	var tempDis float64
	var x, y int

	// Maximum distance
	var d float64 = sheepViewDis

	for i := 0; i < len(state.allGrass); i++ {
		x = state.allGrass[i].x - shp.x
		y = state.allGrass[i].y - shp.y

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
	shp.brain.inputs[3].num = float32(ang)
}

func (shp *Sheep) sheepAngle(state *State) {
	var tempDis float64
	var x, y int
	var d float64 = sheepViewDis

	for i := 0; i < len(state.allSheep); i++ {
		x = state.allSheep[i].x - shp.x
		y = state.allSheep[i].y - shp.y
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
	shp.brain.inputs[4].num = float32(ang)
}

func (shp *Sheep) wolfAngle(state *State) {
	var tempDis float64
	var x, y int
	var d float64 = sheepViewDis

	for i := 0; i < len(state.allWolves); i++ {
		x = state.allWolves[i].x - shp.x
		y = state.allWolves[i].y - shp.y
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
	shp.brain.inputs[5].num = float32(ang)
}

// Output
func (shp *Sheep) bite(state *State) {
	for i := 0; i < len(state.allGrass); i++ {
		// Check collision
		if shp.x > state.allGrass[i].x+grassSize {
			continue
		}
		if shp.x+sheepSize < state.allGrass[i].x {
			continue
		}
		if shp.y > state.allGrass[i].y+grassSize {
			continue
		}
		if shp.y+sheepSize < state.allGrass[i].y {
			continue
		}

		// Eat
		shp.energy += state.allGrass[i].giveEnergy()
		break
	}
}

func (shp *Sheep) move(dis float32) {
	shp.x += int(math.Cos(shp.angle) * sheepSpeed * float64(dis))
	shp.y += int(math.Sin(shp.angle) * sheepSpeed * float64(dis))

	// Stay within bounds
	if shp.x > 1_000 {
		shp.x = 1_000
	} else if shp.x < 0 {
		shp.x = 0
	}
	if shp.y > 1_000 {
		shp.y = 1_000
	} else if shp.y < 0 {
		shp.y = 0
	}
}

func (shp *Sheep) turn(ang float32) {
	shp.angle += float64(ang)
}

func (shp *Sheep) mate(state *State) {
	if len(state.allSheep) >= sheepMaxAmt {
		return
	}

	// Restrictions
	if !shp.canMate() {
		return
	}

	foundPartner := false
	var partner *Sheep

	for i := 0; i < len(state.allSheep); i++ {
		// Check collision
		if shp.x > state.allSheep[i].x+sheepSize {
			continue
		}
		if shp.x+sheepSize < state.allSheep[i].x {
			continue
		}
		if shp.y > state.allSheep[i].y+sheepSize {
			continue
		}
		if shp.y+sheepSize < state.allSheep[i].y {
			continue
		}

		partner = state.allSheep[i]
		// Partner restrictions
		if !partner.canMate() {
			continue
		}
		foundPartner = true
		break
	}

	// Unlucky
	if !foundPartner {
		return
	}
	if partner == shp {
		return
	}

	// Take from the parents and stop them from doing it every frame
	// by giving them a cooldown.
	shp.energy -= sheepMateLoss
	partner.energy -= sheepMateLoss
	shp.mateCooldown = sheepMatePartnerCooldown
	partner.mateCooldown = sheepMatePartnerCooldown

	child := &Sheep{
		shp.x + rand.Intn(20) - 10, // A position near the parent
		shp.y + rand.Intn(20) - 10,
		createBlankBrain(),
		sheepChildEnergy,
		0,
		sheepMateChildCooldown,
	}

	// Copying over brain information into the child
	for i := 0; i < len(child.brain.inputs); i++ {
		for j := 0; j < len(child.brain.inputs[i].links); j++ {
			// Which parent this part of the brain comes from
			if rand.Intn(2) == 0 {
				child.brain.inputs[i].weights = append(child.brain.inputs[i].weights, shp.brain.inputs[i].weights[j])
			} else {
				child.brain.inputs[i].weights = append(child.brain.inputs[i].weights, partner.brain.inputs[i].weights[j])
			}

			// Random chance to get something completely different
			if rand.Float32() < sheepRandWeight {
				child.brain.inputs[i].weights[j] = randWeight()
			}
		}
	}

	for i := 0; i < len(child.brain.layers); i++ {
		for j := 0; j < len(child.brain.layers[i].nodes); j++ {
			if child.brain.layers[i].nodes[j].lastLayer {
				for k := 0; k < len(child.brain.layers[i].nodes[j].linksO); k++ {
					if rand.Intn(2) == 0 {
						child.brain.layers[i].nodes[j].weights = append(child.brain.layers[i].nodes[j].weights, shp.brain.layers[i].nodes[j].weights[k])
					} else {
						child.brain.layers[i].nodes[j].weights = append(child.brain.layers[i].nodes[j].weights, partner.brain.layers[i].nodes[j].weights[k])
					}

					if rand.Float32() < sheepRandWeight {
						child.brain.layers[i].nodes[j].weights[k] = randWeight()
					}
				}
			} else { // Last layer
				for k := 0; k < len(child.brain.layers[i].nodes[j].linksN); k++ {
					if rand.Intn(2) == 0 {
						child.brain.layers[i].nodes[j].weights = append(child.brain.layers[i].nodes[j].weights, shp.brain.layers[i].nodes[j].weights[k])
					} else {
						child.brain.layers[i].nodes[j].weights = append(child.brain.layers[i].nodes[j].weights, partner.brain.layers[i].nodes[j].weights[k])
					}

					if rand.Float32() < sheepRandWeight {
						child.brain.layers[i].nodes[j].weights[k] = randWeight()
					}
				}
			}
		}
	}

	state.allSheep = append(state.allSheep, child) // Make this sheep part of the simulation
}

func (shp *Sheep) canMate() bool {
	if shp.energy < sheepMateBarrier {
		return false
	}
	if shp.mateCooldown > 0 {
		return false
	}
	return true
}

func (shp *Sheep) giveEnergy() int {
	shp.energy -= sheepEnergyLoss
	return sheepEnergyGive
}
