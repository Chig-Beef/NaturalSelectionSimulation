package main

type Grass struct {
	x      int
	y      int
	energy int
}

const grassEnergy int = 10000   // How much energy grass has when made
const grassSize int = 25        // How big grass is
const grassEnergyLoss int = 100 // How much enrgy grass loses when eaten
const grassEnergyGive int = 90  // How much energy grass gives to the sheep when eaten
const grassMaxAmt int = 300     // The maximum amount of grass allowed in the simulation

func (grs *Grass) update(state *State) bool {
	return grs.energy > 0
}

func (grs *Grass) giveEnergy() int {
	grs.energy -= grassEnergyLoss
	return grassEnergyGive
}
