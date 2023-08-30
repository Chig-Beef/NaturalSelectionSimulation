package main

type Grass struct {
	x      int
	y      int
	energy int
	alive  bool
}

const grassEnergy int = 10000 // How much energy grass has when the simulation starts

func (grs *Grass) update(state *State) bool {
	grs.alive = grs.energy > 0
	return grs.alive
}

func (grs *Grass) giveEnergy(state *State) int {
	grs.energy -= state.config.grassEnergyLoss
	return state.config.grassEnergyGive
}
