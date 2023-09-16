package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var max_objects int = 2000

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

		// Edge case that could possibly cause a crash
		if len(text) < 2 {
			fmt.Println("A command went really wrong to hit this error condition.")
			continue
		}

		// Commands are returned with /r/n appended, so this gets rid of those
		command := strings.Split(text[:len(text)-2], " ")

		switch command[0] {
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
		case "remove":
			removeKey(command[1:])
		case "limit":
			limit()
		case "setMax":
			setMax(command[1:])
		default:
			fmt.Println("Invalid Command.")
		}
	}

	fmt.Println("Exiting")
	os.Exit(0) // Ends the program (error code 0 means no error)
}

func save() {
	// No point saving 0 simulations
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
	data := strings.Split(string(f), ",")

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

			// Checkign if values are valid integers, and if so, allocate them

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

			tempBrain, err := convBrainFromStr(obj[5]) // Gets the brain into brain form
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

			tempBrain, err := convBrainFromStr(obj[5])
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
	fmt.Println("help\t\t=>\tLists all commands.")
	fmt.Println("quit\t\t=>\tEnds the server (will freeze all client simulations).")
	fmt.Println("save\t\t=>\tSaves all simulations into a CSV (on this computer).")
	fmt.Println("load\t\t=>\tLoads save file and initialises all simulations.")
	fmt.Println("keys\t\t=>\tReturns all the keys for every active simulation.")
	fmt.Println("remove {keys,}\t=>\tAttempts to delete all the specified keys if the simulations are not being used.")
	fmt.Println("remove all\t=>\tSame as above, except does this check on all simulations.")
	fmt.Println("setMax (number)\t=>\tSets the maximum amount of objects each simulation is allowed (minimum 50).")
}

func removeKey(keys []string) {
	// No point running when no keys are specified
	if len(keys) == 0 {
		return
	}

	// Keyword "all" can be passed to attempt deletion of all keys
	if keys[0] == "all" {
		keys = make([]string, len(simulations))
		i := 0
		for key := range simulations {
			keys[i] = strconv.Itoa(key)
			i++
		}
	}

	for i, value := range keys {
		// Get the key
		key, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println("Key given was not valid int.")
			keys = append(keys[:i], keys[i+1:]...)
			continue
		}

		// Get the sim
		sim, exists := simulations[key]
		if !exists {
			fmt.Println("Key given was not a valid simulation.")
			keys = append(keys[:i], keys[i+1:]...)
			continue
		}

		// Check whether it is active
		sim.active = false
	}

	// Waits for all simulations, instead of waiting per sim
	fmt.Println("Waiting for simulations to respond.")
	time.Sleep(time.Second * 5)

	// allows the user to delete multiple keys in one command
	for _, value := range keys {
		key, _ := strconv.Atoi(value)
		sim := simulations[key]
		if !sim.active {
			// Delete it
			delete(simulations, key)
			fmt.Println("Successfully deleted key " + value + ".")
		} else {
			fmt.Println("Simulation " + value + " was still active.")
		}
	}
	fmt.Println("Done deletion.")
}

func limit() {
	// Get rid of a few objects to help computers die less
	for _, sim := range simulations {
		if len(sim.allGrass) > 150 {
			sim.allGrass = sim.allGrass[:100]
		}
		if len(sim.allSheep) > 150 {
			sim.allSheep = sim.allSheep[:100]
		}
		if len(sim.allWolves) > 150 {
			sim.allWolves = sim.allWolves[:100]
		}
	}
	fmt.Println("Simulations successfully limited.")
}

func setMax(nums []string) {
	if len(nums) != 1 {
		fmt.Println("Expected only one number.")
		return
	}

	// Get the number
	num, err := strconv.Atoi(nums[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	// Simulations have at least some objects
	if num < 50 {
		fmt.Println("Number must be at least 50.")
		return
	}
	max_objects = num

	// Set the value
	for _, sim := range simulations {
		if max_objects < sim.config.sheepMaxAmt {
			sim.config.sheepMaxAmt = max_objects
		}
		if max_objects < sim.config.wolfMaxAmt {
			sim.config.wolfMaxAmt = max_objects
		}
		if max_objects < sim.config.grassMaxAmt {
			sim.config.grassMaxAmt = max_objects
		}
	}

	fmt.Println("Max has been set to", nums[0])
}
