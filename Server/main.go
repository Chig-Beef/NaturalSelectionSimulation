package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

/*
	This file is the routing logic of the server, and connects everything together in a port for the client
*/

const framerate int = 24

func main() {
	// This is asynchronously using the "go" keyword.
	// This means that it starts executing, but when "runCommand"
	// is not busy, it continues executing the rest of this function.
	go runCommand()

	fmt.Println("Server Online (type \"help\" for commands).")

	// Server
	// All this basically says is that if they go to the domain, they are send the "Frontend" folder
	// which holds all the html, css, js, etc. etc.
	// And in a few special cases such as when starting a simulation a certain function is called.
	r := mux.NewRouter()
	r.HandleFunc("/sim/{id}", takeRequest).Methods("GET")       // When getting the next frame in simulation
	r.HandleFunc("/start/{id}", startSimulation).Methods("GET") // When starting a simulation
	r.HandleFunc("/remove/{id}", removeSimulation).Methods("GET")
	r.HandleFunc("/readout/{id}", createReadout).Methods("GET")
	r.HandleFunc("/config/{id}/{data}", changeConfig).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("../Frontend"))) // Getting the regular html
	http.Handle("/", r)

	http.ListenAndServe(":9090", nil)
}

func changeConfig(w http.ResponseWriter, r *http.Request) {
	// Get variables from request
	vars := mux.Vars(r)
	id := vars["id"]
	data := vars["data"]

	// Get the simulation
	simRequest, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("A simulation request wasn't an integer.")
		return
	}
	sim, result := simulations[simRequest]
	if !result {
		fmt.Println("No simulation found.")
		return
	}

	// Make changes in simulation config
	config, err := convConfigFromStr(data)
	if err != nil {
		fmt.Println("The config received from request was bad.")
		fmt.Println(err)
		return
	}
	sim.config = config

	// Respond with success
	_, err = fmt.Fprint(w, "\"Success\"")
	if err != nil {
		fmt.Println("Issue occured in response writing.")
	}
	fmt.Println("Config changed.")
}

func createReadout(w http.ResponseWriter, r *http.Request) {
	// What do we want in the readout?
	// How many of each object there are.
	// The total energy of the whole system

	vars := mux.Vars(r)
	id := vars["id"]

	// Check if it's an integer
	simRequest, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("A simulation request wasn't an integer.")
		return
	}

	sim, result := simulations[simRequest]
	if !result {
		return
	}

	// <br>s are used for newlines in the HTML
	outputData := "\""
	outputData += "Grass: " + strconv.Itoa(len(sim.allGrass)) +
		"<br>" +
		"Sheep: " + strconv.Itoa(len(sim.allSheep)) +
		"<br>" +
		"Wolf: " + strconv.Itoa(len(sim.allWolves)) +
		"<br>" +
		"Energy: " + addCommasToNumber(strconv.Itoa(sim.getEnergy())) +
		"\""
	_, err = fmt.Fprint(w, outputData)
	if err != nil {
		fmt.Println("Issue occured in response writing.")
	}
	fmt.Println("Readout given.")
}

func addCommasToNumber(num string) string {
	// Expects a positive integer as its input
	if len(num) >= 4 {
		num = num[0:len(num)-4] + "," + num[len(num)-4:len(num)-1]
	}
	return num
}

func removeSimulation(w http.ResponseWriter, r *http.Request) {
	// If you saw this line
	// 		r.HandleFunc("/remove/{id}", removeSimulation).Methods("GET")
	// You would notice the "{id}"
	// The next 2 lines just gets that variable.
	// Having {id} allows them to put in any number
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if it's an integer
	simRequest, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("A simulation request wasn't an integer.")
		return
	}

	// Delete the simulation, the delete function
	// takes a map and deletes a key value pair.
	delete(simulations, simRequest)

	_, err = fmt.Fprint(w, "\"Success\"")
	if err != nil {
		fmt.Println("Issue occured in response writing.")
	}

	fmt.Println("Connection successfully ended.")
}

func startSimulation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	simRequest, err := strconv.Atoi(id)
	if err != nil { // Either someone manually entered the wrong address or I wrote the wrong code.
		fmt.Println("A simulation request wasn't an integer.")
		return
	}

	// Create the simulation and send it back
	sim := createNewSimulation(simRequest)
	fmt.Println("Connection successfully made.")
	output := sim.toJson()

	_, err = fmt.Fprint(w, output)
	if err != nil {
		fmt.Println("Issue occured in response writing.")
	}
}

func takeRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	simRequest, err := strconv.Atoi(id)
	if err != nil { // Either someone manually entered the wrong address or I wrote the wrong code.
		fmt.Println("A simulation request wasn't an integer.")
		return
	}

	// result = whether the simulation the client is asking for exists
	sim, result := simulations[simRequest]
	if !result {
		return
	}

	// Calculate all the object's movements and send it back
	output := "["
	for i := 0; i < framerate; i++ {
		sim.step()
		output += sim.toJson() + ","
	}
	output = output[:len(output)-1]
	output += "]"

	sim.active = true

	_, err = fmt.Fprint(w, output)
	if err != nil {
		fmt.Println("Issue occured in response writing.")
	}
}
