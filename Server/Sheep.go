package main

import (
	"fmt"
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
	alive        bool
}

// Look between this file and the wolf file for explanations, as the code is very similar.
const sheepEnergy int = 3600

// Methods
func (shp *Sheep) update(state *State) bool {
	shp.energy--
	shp.mateCooldown--

	// Get the inputs
	shp.grassPosAndAngle(state)
	shp.sheepPosAndAngle(state)
	shp.wolfPosAndAngle(state)
	shp.brain.inputs[6].num = 1 // constant

	// Calculate
	shp.brain.push()
	boolOutput, floatOutput := shp.brain.output_dump()

	// Create changes
	if boolOutput[0] {
		shp.bite(state)
	}
	shp.move(state, floatOutput[1])
	shp.turn(floatOutput[2])
	if boolOutput[3] {
		shp.mate(state)
	}

	// Reset
	shp.brain.set()

	shp.alive = shp.energy > 0
	return shp.alive
}

// Input
func (shp *Sheep) grassPosAndAngle(state *State) {
	var tempDis float64
	var x, y int
	var minX int

	// Maximum distance
	var d float64 = state.config.sheepViewDis * state.config.sheepViewDis

	for i := 0; i < len(state.allGrass); i++ {
		x = state.allGrass[i].x - shp.x
		y = state.allGrass[i].y - shp.y

		// Pythagoras
		tempDis = math.Pow(float64(x), 2) + math.Pow(float64(y), 2)

		// Get the shortest distance
		if tempDis < d {
			d = tempDis
			minX = x
		}
	}

	// Should now be between -1 and 1
	tempDis = math.Sqrt(tempDis)
	tempDis /= state.config.sheepViewDis

	// Range of sight (the sheep can only see so far)
	if tempDis >= 1 {
		return
	}

	// Calculating angle as weight
	cos := float64(minX) / tempDis
	if tempDis == 0 {
		cos = 0
	}
	// Put this input in the brain
	shp.brain.inputs[0].num = float32(tempDis)
	shp.brain.inputs[3].num = float32(cos)
}

func (shp *Sheep) sheepPosAndAngle(state *State) {
	var tempDis float64
	var x, y int
	var minX int
	var d float64 = state.config.sheepViewDis * state.config.sheepViewDis

	for i := 0; i < len(state.allSheep); i++ {
		if state.allSheep[i] == shp {
			continue
		}

		x = state.allSheep[i].x - shp.x
		y = state.allSheep[i].y - shp.y
		tempDis = math.Pow(float64(x), 2) + math.Pow(float64(y), 2)
		if tempDis < d {
			d = tempDis
			minX = x
		}
	}
	tempDis = math.Sqrt(tempDis)
	tempDis /= state.config.sheepViewDis
	if tempDis >= 1 {
		return
	}
	cos := float64(minX) / tempDis
	if tempDis == 0 {
		cos = 0
	}
	shp.brain.inputs[1].num = float32(tempDis)
	shp.brain.inputs[4].num = float32(cos)
}

func (shp *Sheep) wolfPosAndAngle(state *State) {
	var tempDis float64
	var x, y int
	var minX int
	var d float64 = state.config.sheepViewDis * state.config.sheepViewDis

	for i := 0; i < len(state.allWolves); i++ {
		x = state.allWolves[i].x - shp.x
		y = state.allWolves[i].y - shp.y
		tempDis = math.Pow(float64(x), 2) + math.Pow(float64(y), 2)
		if tempDis < d {
			d = tempDis
			minX = x
		}
	}
	tempDis = math.Sqrt(tempDis)
	tempDis /= state.config.sheepViewDis
	if tempDis >= 1 {
		return
	}
	cos := float64(minX) / tempDis
	if tempDis == 0 {
		cos = 0
	}
	if math.IsNaN(cos) {
		fmt.Println("What the Heck")
	}
	shp.brain.inputs[2].num = float32(tempDis)
	shp.brain.inputs[5].num = float32(cos)
}

// Output
func (shp *Sheep) bite(state *State) {
	wide := shp.x + state.config.sheepSize
	high := shp.y + state.config.sheepSize
	for i := 0; i < len(state.allGrass); i++ {
		// Check collision
		if shp.x > state.allGrass[i].x+state.config.grassSize {
			continue
		}
		if wide < state.allGrass[i].x {
			continue
		}
		if shp.y > state.allGrass[i].y+state.config.grassSize {
			continue
		}
		if high < state.allGrass[i].y {
			continue
		}

		// Eat
		shp.energy += state.allGrass[i].giveEnergy(state)
		break
	}
}

func (shp *Sheep) move(state *State, dis float32) {
	degAng := shp.angle * math.Pi / 180
	shp.x += int(math.Cos(degAng) * state.config.sheepSpeed * float64(dis))
	shp.y += int(math.Sin(degAng) * state.config.sheepSpeed * float64(dis))

	// Stay within bounds
	if shp.x > 1_000-state.config.sheepSize {
		shp.x = 1_000 - state.config.sheepSize
	} else if shp.x < 0 {
		shp.x = 0
	}
	if shp.y > 1_000-state.config.sheepSize {
		shp.y = 1_000 - state.config.sheepSize
	} else if shp.y < 0 {
		shp.y = 0
	}
}

func (shp *Sheep) turn(ang float32) {
	shp.angle += float64(ang) * 5
}

func (shp *Sheep) mate(state *State) {
	if len(state.allSheep) >= state.config.sheepMaxAmt {
		return
	}

	// Restrictions
	if !shp.canMate(state) {
		return
	}

	foundPartner := false
	var partner *Sheep

	wide := shp.x + state.config.wolfSize
	high := shp.y + state.config.wolfSize
	for i := 0; i < len(state.allSheep); i++ {
		// Check collision
		if shp.x > state.allSheep[i].x+state.config.sheepSize {
			continue
		}
		if wide < state.allSheep[i].x {
			continue
		}
		if shp.y > state.allSheep[i].y+state.config.sheepSize {
			continue
		}
		if high < state.allSheep[i].y {
			continue
		}
		if state.allSheep[i] == shp {
			continue
		}

		partner = state.allSheep[i]
		// Partner restrictions
		if !partner.canMate(state) {
			continue
		}
		foundPartner = true
		break
	}

	// Unlucky
	if !foundPartner {
		return
	}

	// Take from the parents and stop them from doing it every frame
	// by giving them a cooldown.
	shp.energy -= state.config.sheepMateLoss
	partner.energy -= state.config.sheepMateLoss
	shp.mateCooldown = state.config.sheepMatePartnerCooldown
	partner.mateCooldown = state.config.sheepMatePartnerCooldown

	child := &Sheep{
		shp.x + rand.Intn(20) - 10, // A position near the parent
		shp.y + rand.Intn(20) - 10,
		createBlankBrain(),
		state.config.sheepChildEnergy,
		0,
		state.config.sheepMateChildCooldown,
		true,
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
			if rand.Float32() < state.config.sheepRandWeight {
				child.brain.inputs[i].weights[j] = randWeight()
			}
		}
	}

	for i := 0; i < len(child.brain.layers); i++ {
		for j := 0; j < len(child.brain.layers[i].nodes); j++ {
			if !child.brain.layers[i].nodes[j].lastLayer { // Last layer
				for k := 0; k < len(child.brain.layers[i].nodes[j].linksN); k++ {
					if rand.Intn(2) == 0 {
						child.brain.layers[i].nodes[j].weights = append(child.brain.layers[i].nodes[j].weights, shp.brain.layers[i].nodes[j].weights[k])
					} else {
						child.brain.layers[i].nodes[j].weights = append(child.brain.layers[i].nodes[j].weights, partner.brain.layers[i].nodes[j].weights[k])
					}

					if rand.Float32() < state.config.sheepRandWeight {
						child.brain.layers[i].nodes[j].weights[k] = randWeight()
					}
				}
			} else {
				for k := 0; k < len(child.brain.layers[i].nodes[j].linksO); k++ {
					if rand.Intn(2) == 0 {
						child.brain.layers[i].nodes[j].weights = append(child.brain.layers[i].nodes[j].weights, shp.brain.layers[i].nodes[j].weights[k])
					} else {
						child.brain.layers[i].nodes[j].weights = append(child.brain.layers[i].nodes[j].weights, partner.brain.layers[i].nodes[j].weights[k])
					}

					if rand.Float32() < state.config.sheepRandWeight {
						child.brain.layers[i].nodes[j].weights[k] = randWeight()
					}
				}
			}
		}
	}

	state.allSheep = append(state.allSheep, child) // Make this sheep part of the simulation
}

func (shp *Sheep) canMate(state *State) bool {
	if shp.energy < state.config.sheepMateBarrier {
		return false
	}
	if shp.mateCooldown > 0 {
		return false
	}
	return true
}

func (shp *Sheep) giveEnergy(state *State) int {
	shp.energy -= state.config.sheepEnergyLoss
	return state.config.sheepEnergyGive
}
