package main

import (
	"math/rand"
	"strconv"
)

type State struct {
	allGrass    []*Grass
	allWolves   []*Wolf
	allSheep    []*Sheep
	grassSince  int
	sheepSince  int
	wolvesSince int
	active      bool
	config      Config
	frame       int
}

var simulations map[int]*State = make(map[int]*State) // All the student's simulations

const startingGrass int = 100
const startingWolf int = 30
const startingSheep int = 50
const grassChance int = 10 // Change that grass spawns every frame. Out of 1,000.

func (state *State) step() {
	state.frame++

	// Steps a State one frame
	if len(state.allGrass) < state.config.grassMaxAmt {
		placeGrass := rand.Intn(1_000)
		if placeGrass < grassChance {

			// Generate coordinates and give the grass energy, then allocate it
			state.allGrass = append(state.allGrass, &Grass{
				rand.Intn(1_000 - state.config.grassSize),
				rand.Intn(1_000 - state.config.grassSize),
				state.config.grassEnergy,
				true,
			})
		}
	}

	var i int
	i = 0

	// These are while loops (golang only has the for keyword)
	// The else cases are if the thing has died, and so they are therefore removed
	for i < len(state.allWolves) {
		if state.allWolves[i].update(state) {
			i++
		} else {
			state.allWolves = append(state.allWolves[:i], state.allWolves[i+1:]...)
		}
	}

	i = 0
	for i < len(state.allSheep) {
		if state.allSheep[i].update(state) {
			i++
		} else {
			state.allSheep = append(state.allSheep[:i], state.allSheep[i+1:]...)
		}
	}

	i = 0
	for i < len(state.allGrass) {
		if state.allGrass[i].update(state) {
			i++
		} else {
			state.allGrass = append(state.allGrass[:i], state.allGrass[i+1:]...)
		}
	}
}

func (state State) toJson() string {
	// This is sent back to the client over http
	// Just a bunch of string concat

	grassText := "["
	wolfText := "["
	sheepText := "["

	// Grass
	if len(state.allGrass) > 0 {
		grassText += "[" + strconv.Itoa(state.allGrass[0].x) + "," + strconv.Itoa(state.allGrass[0].y) + "]"
		for i := 1; i < len(state.allGrass); i++ {
			grassText += ",[" +
				strconv.Itoa(state.allGrass[i].x) + "," +
				strconv.Itoa(state.allGrass[i].y) +
				"]"
		}
	}
	grassText += "]"

	// Wolf
	if len(state.allWolves) > 0 {
		wolfText += "[" + strconv.Itoa(state.allWolves[0].x) + "," + strconv.Itoa(state.allWolves[0].y) + "]"
		for i := 1; i < len(state.allWolves); i++ {
			wolfText += ",[" +
				strconv.Itoa(state.allWolves[i].x) + "," +
				strconv.Itoa(state.allWolves[i].y) +
				"]"
		}
	}
	wolfText += "]"

	// Sheep
	if len(state.allSheep) > 0 {
		sheepText += "[" + strconv.Itoa(state.allSheep[0].x) + "," + strconv.Itoa(state.allSheep[0].y) + "]"
		for i := 1; i < len(state.allSheep); i++ {
			sheepText += ",[" +
				strconv.Itoa(state.allSheep[i].x) + "," +
				strconv.Itoa(state.allSheep[i].y) +
				"]"
		}
	}
	sheepText += "]"

	// Concat and return
	output := "[" + strconv.Itoa(state.frame) + "," + grassText + "," + wolfText + "," + sheepText + "]"

	return output
}

func createNewSimulation(id int) *State {
	// Creates a State with a bunch of random objects
	simulations[id] = &State{
		initializeGrassSlice(),
		initializeWolfSlice(),
		initializeSheepSlice(),
		startingGrass,
		startingWolf,
		startingSheep,
		true,
		makeDefaultConfig(),
		0,
	}
	return simulations[id]
}

func initializeGrassSlice() []*Grass {
	grassSlice := []*Grass{}

	for i := 0; i < startingGrass; i++ {
		grassSlice = append(grassSlice, &Grass{
			rand.Intn(1_000 - 25),
			rand.Intn(1_000 - 25),
			grassEnergy,
			true,
		})
	}

	return grassSlice
}

func initializeWolfSlice() []*Wolf {
	wolfSlice := []*Wolf{}

	for i := 0; i < startingWolf; i++ {
		wolfSlice = append(wolfSlice, &Wolf{
			rand.Intn(1_000 - 25),
			rand.Intn(1_000 - 25),
			createRandomConnections(createBlankBrain()),
			wolfEnergy,
			0,
			0,
			true,
		})
	}

	return wolfSlice
}

func initializeSheepSlice() []*Sheep {
	sheepSlice := []*Sheep{}

	for i := 0; i < startingSheep; i++ {
		sheepSlice = append(sheepSlice, &Sheep{
			rand.Intn(1_000 - 25),
			rand.Intn(1_000 - 25),
			createRandomConnections(createBlankBrain()),
			sheepEnergy,
			0,
			0,
			true,
		})
	}

	return sheepSlice
}

func (state *State) getEnergy() int {
	totalEnergy := 0
	for _, grs := range state.allGrass {
		totalEnergy += grs.energy
	}

	for _, shp := range state.allSheep {
		totalEnergy += shp.energy
	}

	for _, wlf := range state.allWolves {
		totalEnergy += wlf.energy
	}
	return totalEnergy
}
