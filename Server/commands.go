package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func runCommand() {
	running := true
	reader := bufio.NewReader(os.Stdin)

	// Infinite loop, always checking for the next command
	for running {
		// Get a command
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}

		// Edge case that could cause a crash
		if len(text) < 2 {
			fmt.Println("A command went really wrong to hit this error condition.")
			continue
		}

		// Commands are returned with /r/n appended, so this gets rid of those
		text = text[:len(text)-2]

		switch text {
		case "quit":
			running = false // Breaks the loop
		case "save":
			save()
		case "load":
			load()
		case "help":
			help()
		case "keys":
			getKeys()
		default:
			fmt.Println("Invalid Command.")
		}
	}

	fmt.Println("Exiting")
	os.Exit(0) // Ends the program (error code 0 means no error)
}

func save() {
	if len(simulations) == 0 {
		fmt.Println("There was nothing to save.")
		return
	}
	fmt.Println("Saving...")

	// Putting everything in string variation
	// and concat with various delimeters
	outputString := ""
	for index, sim := range simulations {
		// The key for the simulation
		outputString += strconv.Itoa(index) + "?"

		for _, grs := range sim.allGrass {
			outputString += strconv.Itoa(grs.energy) +
				"|" +
				strconv.Itoa(grs.x) +
				"|" +
				strconv.Itoa(grs.y) +
				";"
		}

		// Doing the len-1 gets rid of that extra ";" from the last thing.
		// Getting rid of it makes it easier to split.
		outputString = outputString[:len(outputString)-1] + "\n"

		for _, shp := range sim.allSheep {
			outputString += strconv.Itoa(shp.energy) +
				"|" +
				strconv.Itoa(shp.x) +
				"|" +
				strconv.Itoa(shp.y) +
				"|" +
				strconv.Itoa(shp.mateCooldown) +
				"|" +
				strconv.FormatFloat(shp.angle, 'f', -1, 64) +
				"|" +
				shp.brain.convToStr() +
				";"
		}

		outputString = outputString[:len(outputString)-1] + "\n"

		for _, wlf := range sim.allWolves {
			outputString += strconv.Itoa(wlf.energy) +
				"|" +
				strconv.Itoa(wlf.x) +
				"|" +
				strconv.Itoa(wlf.y) +
				"|" +
				strconv.Itoa(wlf.mateCooldown) +
				"|" +
				strconv.FormatFloat(wlf.angle, 'f', -1, 64) +
				"|" +
				wlf.brain.convToStr() +
				";"
		}

		outputString = outputString[:len(outputString)-1] + ","
	}

	outputString = outputString[:len(outputString)-1]

	// Write to the file and close it
	f, err := os.Create("save.csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	f.WriteString(outputString)
	err = f.Sync()
	if err != nil {
		fmt.Println(err)
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Saved.")
}

func load() {
	fmt.Println("Loading...")

	// Get the file
	f, err := os.ReadFile("save.csv")
	if err != nil {
		fmt.Println(err)
		return
	}

	rawData := string(f)
	data := strings.Split(rawData, ",")

	for _, sim := range data {

		// Get the key to put this simulation in and save it for later
		length := strings.Index(sim, "?")
		index, err := strconv.Atoi(sim[:length])
		if err != nil {
			fmt.Println("A Simulation key was an invalid integer.")
			continue
		}
		// Then remove this key (don't need it anymore)
		sim := sim[length+1:]

		// This should have grass, sheep, and wolves, no more, no less
		splitSim := strings.Split(sim, "\n")
		if len(splitSim) != 3 {
			fmt.Println("A Simulation was found that did not have a grass, sheep, and wolf dataset (either too few or too many).")
			continue
		}

		// Create a temp State to eventually store
		state := &State{}

		grass := strings.Split(splitSim[0], ";")
		for i, grs := range grass {
			state.allGrass = append(state.allGrass, &Grass{})

			obj := strings.Split(grs, "|")
			if len(obj) != 3 {
				fmt.Println("Grass object found in state that didn't have the correct number of variables.")
				continue
			}

			// Checkign if values are valid integers, and if so, allocating them

			temp, err := strconv.Atoi(obj[0])
			if err != nil {
				fmt.Println(err)
				continue
			}
			state.allGrass[i].energy = temp

			temp, err = strconv.Atoi(obj[1])
			if err != nil {
				fmt.Println(err)
				continue
			}
			state.allGrass[i].x = temp

			temp, err = strconv.Atoi(obj[2])
			if err != nil {
				fmt.Println(err)
				continue
			}
			state.allGrass[i].y = temp
		}

		sheep := strings.Split(splitSim[1], ";")
		for i, shp := range sheep {
			state.allSheep = append(state.allSheep, &Sheep{})

			obj := strings.Split(shp, "|")

			if len(obj) != 6 {
				fmt.Println("Sheep object found in state that didn't have the correct number of variables.")
				continue
			}

			temp, err := strconv.Atoi(obj[0])
			if err != nil {
				fmt.Println(err)
				continue
			}
			state.allSheep[i].energy = temp

			temp, err = strconv.Atoi(obj[1])
			if err != nil {
				fmt.Println(err)
				continue
			}
			state.allSheep[i].x = temp

			temp, err = strconv.Atoi(obj[2])
			if err != nil {
				fmt.Println(err)
				continue
			}
			state.allSheep[i].y = temp

			temp, err = strconv.Atoi(obj[3])
			if err != nil {
				fmt.Println(err)
				continue
			}
			state.allSheep[i].mateCooldown = temp

			tempFloat, err := strconv.ParseFloat(obj[4], 64) // Like int checking, but for floats, the 64 just means 64 bit (like a double)
			if err != nil {
				fmt.Println(err)
				continue
			}
			state.allSheep[i].angle = tempFloat

			tempBrain, err := convFromStr(obj[5]) // Gets the brain into brain form
			if err != nil {
				fmt.Println(err)
				continue
			}
			state.allSheep[i].brain = tempBrain
		}

		wolves := strings.Split(splitSim[2], ";")
		for i, wlf := range wolves {
			state.allWolves = append(state.allWolves, &Wolf{})

			obj := strings.Split(wlf, "|")

			if len(obj) != 6 {
				fmt.Println("Wolf object found in state that didn't have the correct number of variables.")
				continue
			}

			temp, err := strconv.Atoi(obj[0])
			if err != nil {
				fmt.Println(err)
				continue
			}
			state.allWolves[i].energy = temp

			temp, err = strconv.Atoi(obj[1])
			if err != nil {
				fmt.Println(err)
				continue
			}
			state.allWolves[i].x = temp

			temp, err = strconv.Atoi(obj[2])
			if err != nil {
				fmt.Println(err)
				continue
			}
			state.allWolves[i].y = temp

			temp, err = strconv.Atoi(obj[3])
			if err != nil {
				fmt.Println(err)
				continue
			}
			state.allWolves[i].mateCooldown = temp

			tempFloat, err := strconv.ParseFloat(obj[4], 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			state.allWolves[i].angle = tempFloat

			tempBrain, err := convFromStr(obj[5])
			if err != nil {
				fmt.Println(err)
				continue
			}
			state.allWolves[i].brain = tempBrain
		}

		// Allocate the state
		simulations[index] = state
	}

	fmt.Println("Loaded.")
}

func getKeys() {
	fmt.Println("There are currently " + strconv.Itoa(len(simulations)) + " valid keys.")

	for key := range simulations {
		fmt.Println(key)
	}
}

func help() {
	fmt.Println("help\t=>\tLists all commands.")
	fmt.Println("quit\t=>\tEnds the server (will freeze all client simulations).")
	fmt.Println("save\t=>\tSaves all simulations into a CSV (on this computer).")
	fmt.Println("load\t=>\tLoads save file and initialises all simulations.")
	fmt.Println("keys\t=>\tReturns all the keys for every active simulation.")
}
