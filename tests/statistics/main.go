package main

import (
	tests "tests/game-data"
	teststats "tests/test-stats"
)

func main() {
	// Load JSON from file
	pd, _ := tests.MakeGameData("kst-data.json")
	teststats.RunReelSlotMachine(1000, pd)

}
