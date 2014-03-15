// Codingame - Game of Drones
package main

import (
	"os"
	"testing"
)

const FILE_TESTS_BASE = `C:\Users\borja\programacion\codinggame\GameOfDronesSolution\testInputs\`

//Tests method turnBasedDistance
func TestTurnBasedDistance(t *testing.T) {
	var testCases = []struct {
		in1, in2 point
		out      int
	}{
		{point{500, 500}, point{450, 500}, 0}, //up - 0
		{point{500, 500}, point{550, 500}, 0}, //down - 0
		{point{500, 500}, point{500, 450}, 0}, //left - 0
		{point{500, 500}, point{500, 550}, 0}, //right - 0

		{point{500, 500}, point{400, 500}, 0}, //up - 0
		{point{500, 500}, point{600, 500}, 0}, //down - 0
		{point{500, 500}, point{500, 400}, 0}, //left - 0
		{point{500, 500}, point{500, 600}, 0}, //right - 0

		{point{500, 500}, point{399, 500}, 1}, //up - 1
		{point{500, 500}, point{601, 500}, 1}, //down - 1
		{point{500, 500}, point{500, 399}, 1}, //left - 1
		{point{500, 500}, point{500, 601}, 1}, //right - 1

		{point{500, 500}, point{299, 500}, 2}, //up - 2
		{point{500, 500}, point{701, 500}, 2}, //down - 2
		{point{500, 500}, point{500, 299}, 2}, //left - 2
		{point{500, 500}, point{500, 701}, 2}, //right - 2

		{point{500, 500}, point{430, 570}, 0}, //up-righ - 0
		{point{500, 500}, point{430, 430}, 0}, //up-left - 0
		{point{500, 500}, point{570, 570}, 0}, //down-right - 0
		{point{500, 500}, point{570, 430}, 0}, //down-left - 0

		{point{500, 500}, point{428, 572}, 1}, //up-righ - 1
		{point{500, 500}, point{428, 428}, 1}, //up-left - 1
		{point{500, 500}, point{572, 572}, 1}, //down-right - 1
		{point{500, 500}, point{572, 428}, 1}, //down-left - 1

		{point{0, 0}, point{0, 0}, 0}, //border

		{point{500, 500}, point{500, 500}, 0}, //same point
	}
	for i, testCase := range testCases {
		if result := turnBasedDistance(testCase.in1, testCase.in2); testCase.out != result {
			t.Error("Error in item", i, "Got", result, "Expected", testCase.out, "Case:", testCase)
		}
	}
}

//Tests method euclideanDistance
func TestEuclideanDistance(t *testing.T) {
	var testCases = []struct {
		in1, in2 point
		out      float64
	}{
		{point{500, 500}, point{450, 500}, 50.0}, //up - 0
		{point{500, 500}, point{550, 500}, 50.0}, //down - 0
		{point{500, 500}, point{500, 450}, 50.0}, //left - 0
		{point{500, 500}, point{500, 550}, 50.0}, //right - 0

		{point{500, 500}, point{430, 570}, 98.0}, //up-righ - 0
		{point{500, 500}, point{430, 430}, 98.0}, //up-left - 0
		{point{500, 500}, point{570, 570}, 98.0}, //down-right - 0
		{point{500, 500}, point{570, 430}, 98.0}, //down-left - 0

		{point{500, 500}, point{429, 571}, 100.0}, //up-righ - 1
		{point{500, 500}, point{429, 429}, 100.0}, //up-left - 1
		{point{500, 500}, point{571, 571}, 100.0}, //down-right - 1
		{point{500, 500}, point{571, 429}, 100.0}, //down-left - 1

		{point{0, 0}, point{0, 0}, 0}, //border

		{point{500, 500}, point{500, 500}, 0}, //same point
	}
	for i, testCase := range testCases {
		if result := euclideanDistance(testCase.in1, testCase.in2); testCase.out != result {
			t.Error("Error in item", i, "Got", result, "Expected", testCase.out, "Case:", testCase)
		}
	}
}

//Tests method colonizeTheUnexplored with Zero unexplored zones
func TestColonizeTheUnexploredZero(t *testing.T) {
	if f, err := os.Open(FILE_TESTS_BASE + "colonizeTheUnexplored\\inputZero.txt"); err != nil {
		t.Error("Error opening input file")
	} else {
		inputReader = f
		defer f.Close()
	}
	readBoard()
	parseTurn()
	initializeTurnComputation()
	colonizeTheUnexplored()
	if len(dronesAsigned) != 0 {
		t.Error("Too many drones asigned:", len(dronesAsigned))
	}
}

//Tests method colonizeTheUnexplored with One unexplored zone
func TestColonizeTheUnexploredOne(t *testing.T) {
	if f, err := os.Open(FILE_TESTS_BASE + "colonizeTheUnexplored\\inputOne.txt"); err != nil {
		t.Error("Error opening input file")
	} else {
		inputReader = f
		defer f.Close()
	}
	readBoard()
	parseTurn()
	initializeTurnComputation()
	colonizeTheUnexplored()
	if len(dronesAsigned) != 1 {
		t.Error("Wrong number of drones asigned:", len(dronesAsigned))
	}
	for d, _ := range dronesAsigned {
		if nextMove[d].x != 100 || nextMove[d].y != 100 {
			t.Error("Wrong movement. Going to", nextMove[d])
		}
	}
}

//Tests method colonizeTheUnexplored with Two unexplored zones
func TestColonizeTheUnexploredTwo(t *testing.T) {
	if f, err := os.Open(FILE_TESTS_BASE + "colonizeTheUnexplored\\inputTwo.txt"); err != nil {
		t.Error("Error opening input file")
	} else {
		inputReader = f
		defer f.Close()
	}
	readBoard()
	parseTurn()
	initializeTurnComputation()
	colonizeTheUnexplored()
	if len(dronesAsigned) != 2 {
		t.Error("Wrong number of drones asigned:", len(dronesAsigned))
	}
	destinations := make(map[point]bool, 2)
	for d, _ := range dronesAsigned {
		if (nextMove[d].x == 100 && nextMove[d].y == 100) ||
			(nextMove[d].x == 200 && nextMove[d].y == 200) {
			destinations[nextMove[d]] = true
		}
	}
	if len(destinations) != 2 {
		t.Error("Wrong number of different destinations:", len(destinations))
	}
}
