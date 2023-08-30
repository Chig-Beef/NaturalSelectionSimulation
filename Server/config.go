package main

import (
	"errors"
	"strconv"
	"strings"
)

type Config struct {
	sheepSpeed               float64
	sheepSize                int
	sheepRandWeight          float32
	sheepMateBarrier         int
	sheepMateLoss            int
	sheepMatePartnerCooldown int
	sheepMateChildCooldown   int
	sheepEnergyLoss          int
	sheepEnergyGive          int
	sheepChildEnergy         int
	sheepMaxAmt              int
	sheepViewDis             float64
	wolfSpeed                float64 // How fast the wolf is
	wolfSize                 int     // How big a wolf is
	wolfRandWeight           float32 // The chance that a new weight is random rather than inherited
	wolfMateBarrier          int     // How much energy a wolf needs to breed
	wolfMateLoss             int     // How much energy a wolf loses from breeding
	wolfMatePartnerCooldown  int     // How long from giving birth a wolf can attempt mating again
	wolfMateChildCooldown    int     // How long from being born a wolf can attempt to mate
	wolfChildEnergy          int     // How much energy a wolf has when it is born
	wolfMaxAmt               int     // Maximum amount of wolves
	wolfViewDis              float64 // How far a wolf can see
	grassEnergy              int     // How much energy grass has when made
	grassSize                int     // How big grass is
	grassEnergyLoss          int     // How much enrgy grass loses when eaten
	grassEnergyGive          int     // How much energy grass gives to the sheep when eaten
	grassMaxAmt              int     // The maximum amount of grass allowed in the simulation
}

func makeDefaultConfig() Config {
	return Config{
		5,
		25,
		float32(1) / 100,
		1800,
		900,
		240,
		360,
		100,
		80,
		1800,
		300,
		100,
		5,
		25,
		float32(1) / 100,
		1800,
		900,
		240,
		360,
		1800,
		300,
		100,
		10_000,
		25,
		100,
		90,
		300,
	}
}

func convConfigFromStr(data string) (Config, error) {
	newConfig := Config{}
	fields := strings.Split(data[1:len(data)-1], ",")

	// Just to make sure
	if len(fields) != 27 {
		return newConfig, errors.New("expected 27 fields in config")
	}

	// Converting each field
	newFloat, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return newConfig, err
	}
	newConfig.sheepSpeed = newFloat

	newInt, err := strconv.Atoi(fields[1])
	if err != nil {
		return newConfig, err
	}
	newConfig.sheepSize = newInt

	newFloat, err = strconv.ParseFloat(fields[2], 32)
	if err != nil {
		return newConfig, err
	}
	newConfig.sheepRandWeight = float32(newFloat)

	newInt, err = strconv.Atoi(fields[3])
	if err != nil {
		return newConfig, err
	}
	newConfig.sheepMateBarrier = newInt

	newInt, err = strconv.Atoi(fields[4])
	if err != nil {
		return newConfig, err
	}
	newConfig.sheepMateLoss = newInt

	newInt, err = strconv.Atoi(fields[5])
	if err != nil {
		return newConfig, err
	}
	newConfig.sheepMatePartnerCooldown = newInt

	newInt, err = strconv.Atoi(fields[6])
	if err != nil {
		return newConfig, err
	}
	newConfig.sheepMateChildCooldown = newInt

	newInt, err = strconv.Atoi(fields[7])
	if err != nil {
		return newConfig, err
	}
	newConfig.sheepEnergyLoss = newInt

	newInt, err = strconv.Atoi(fields[8])
	if err != nil {
		return newConfig, err
	}
	newConfig.sheepEnergyGive = newInt

	newInt, err = strconv.Atoi(fields[9])
	if err != nil {
		return newConfig, err
	}
	newConfig.sheepChildEnergy = newInt

	newInt, err = strconv.Atoi(fields[10])
	if err != nil {
		return newConfig, err
	}
	newConfig.sheepMaxAmt = newInt

	newFloat, err = strconv.ParseFloat(fields[11], 64)
	if err != nil {
		return newConfig, err
	}
	newConfig.sheepViewDis = newFloat

	newFloat, err = strconv.ParseFloat(fields[12], 64)
	if err != nil {
		return newConfig, err
	}
	newConfig.wolfSpeed = newFloat

	newInt, err = strconv.Atoi(fields[13])
	if err != nil {
		return newConfig, err
	}
	newConfig.wolfSize = newInt

	newFloat, err = strconv.ParseFloat(fields[14], 32)
	if err != nil {
		return newConfig, err
	}
	newConfig.wolfRandWeight = float32(newFloat)

	newInt, err = strconv.Atoi(fields[15])
	if err != nil {
		return newConfig, err
	}
	newConfig.wolfMateBarrier = newInt

	newInt, err = strconv.Atoi(fields[16])
	if err != nil {
		return newConfig, err
	}
	newConfig.wolfMateLoss = newInt

	newInt, err = strconv.Atoi(fields[17])
	if err != nil {
		return newConfig, err
	}
	newConfig.wolfMatePartnerCooldown = newInt

	newInt, err = strconv.Atoi(fields[18])
	if err != nil {
		return newConfig, err
	}
	newConfig.wolfMateChildCooldown = newInt

	newInt, err = strconv.Atoi(fields[19])
	if err != nil {
		return newConfig, err
	}
	newConfig.wolfChildEnergy = newInt

	newInt, err = strconv.Atoi(fields[20])
	if err != nil {
		return newConfig, err
	}
	newConfig.wolfMaxAmt = newInt

	newFloat, err = strconv.ParseFloat(fields[21], 64)
	if err != nil {
		return newConfig, err
	}
	newConfig.wolfViewDis = newFloat

	newInt, err = strconv.Atoi(fields[22])
	if err != nil {
		return newConfig, err
	}
	newConfig.grassEnergy = newInt

	newInt, err = strconv.Atoi(fields[23])
	if err != nil {
		return newConfig, err
	}
	newConfig.grassSize = newInt

	newInt, err = strconv.Atoi(fields[24])
	if err != nil {
		return newConfig, err
	}
	newConfig.grassEnergyLoss = newInt

	newInt, err = strconv.Atoi(fields[25])
	if err != nil {
		return newConfig, err
	}
	newConfig.grassEnergyLoss = newInt

	newInt, err = strconv.Atoi(fields[26])
	if err != nil {
		return newConfig, err
	}
	newConfig.grassMaxAmt = newInt

	// Return
	return newConfig, nil
}
