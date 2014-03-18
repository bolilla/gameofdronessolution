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

//Tests method strategyColonizeTheUnexplored with Zero unexplored zones
func TestColonizeTheUnexploredZero(t *testing.T) {
	setUpTestFromFile(FILE_TESTS_BASE+"colonizeTheUnexplored\\inputZero.txt", t)
	strategyColonizeTheUnexplored()
	if len(assignedDrones) != 0 {
		t.Error("Too many drones asigned:", len(assignedDrones))
	}
}

//Tests method strategyColonizeTheUnexplored with One unexplored zone
func TestColonizeTheUnexploredOne(t *testing.T) {
	setUpTestFromFile(FILE_TESTS_BASE+"colonizeTheUnexplored\\inputOne.txt", t)
	strategyColonizeTheUnexplored()
	if len(assignedDrones) != 1 {
		t.Error("Wrong number of drones asigned:", len(assignedDrones))
	}
	for d, _ := range assignedDrones {
		if nextMove[d].x != 100 || nextMove[d].y != 100 {
			t.Error("Wrong movement. Going to", nextMove[d])
		}
	}
}

//Tests method strategyColonizeTheUnexplored with Two unexplored zones
func TestColonizeTheUnexploredTwo(t *testing.T) {
	setUpTestFromFile(FILE_TESTS_BASE+"colonizeTheUnexplored\\inputTwo.txt", t)
	strategyColonizeTheUnexplored()
	if len(assignedDrones) != 2 {
		t.Error("Wrong number of drones asigned:", len(assignedDrones))
	}
	destinations := make(map[point]bool, 2)
	for d, _ := range assignedDrones {
		if (nextMove[d].x == 100 && nextMove[d].y == 100) ||
			(nextMove[d].x == 200 && nextMove[d].y == 200) {
			destinations[nextMove[d]] = true
		}
	}
	if len(destinations) != 2 {
		t.Error("Wrong number of different destinations:", len(destinations))
	}
}

//Tests method strategyColonizeTheUnexplored with an unreclaimed zone that already contains drones from different players
func TestColonizeTheUnexploredUnreclaimedPopulated(t *testing.T) {
	setUpTestFromFile(FILE_TESTS_BASE+"colonizeTheUnexplored\\inputUnconqueredButPopulated.txt", t)
	strategyColonizeTheUnexplored()
	if len(assignedDrones) != 2 {
		t.Error("Wrong number of drones asigned:", len(assignedDrones))
	}
	destinations := make(map[point]bool, 2)
	for d, _ := range assignedDrones {
		if nextMove[d].x != 100 || nextMove[d].y != 100 {
			t.Error("Wrong movement. Going to", nextMove[d])
		} else {
			destinations[nextMove[d]] = true
		}
	}
	if len(destinations) != 1 {
		t.Error("Wrong number of different destinations:", len(destinations), assignedDrones)
	}
}

//Tests method strategyMaintainAirSuperiority with zero zones owned
func TestMaintainAirSuperiority0(t *testing.T) {
	setUpTestFromFile(FILE_TESTS_BASE+"maintainAirSuperiority\\inputOwned0.txt", t)

	strategyMaintainAirSuperiority()
	if len(assignedDrones) != 0 {
		t.Error("Wrong number of drones asigned:", len(assignedDrones))
	}
}

//Tests method strategyMaintainAirSuperiority with one zone owned 1 Vs 0
func TestMaintainAirSuperiority1Vs0(t *testing.T) {
	setUpTestFromFile(FILE_TESTS_BASE+"maintainAirSuperiority\\inputOwned1Vs0.txt", t)

	strategyMaintainAirSuperiority()
	if len(assignedDrones) != 0 {
		t.Error("Wrong number of drones asigned:", len(assignedDrones))
	}
}

//Tests method strategyMaintainAirSuperiority with one zone owned 1 Vs 1
func TestMaintainAirSuperiority1Vs1(t *testing.T) {
	setUpTestFromFile(FILE_TESTS_BASE+"maintainAirSuperiority\\inputOwned1Vs1.txt", t)

	strategyMaintainAirSuperiority()
	if len(assignedDrones) != 1 {
		t.Error("Wrong number of drones asigned:", len(assignedDrones))
	}
	destinations := make(map[point]bool, 2)
	for d, _ := range assignedDrones {
		if nextMove[d].x == 300 && nextMove[d].y == 300 {
			destinations[nextMove[d]] = true
		}
	}
	if len(destinations) != 1 {
		t.Error("Wrong number of different destinations:", len(destinations))
	}
}

//Tests method strategyMaintainAirSuperiority with one zone owned 2 Vs 1
func TestMaintainAirSuperiority2Vs1(t *testing.T) {
	setUpTestFromFile(FILE_TESTS_BASE+"maintainAirSuperiority\\inputOwned2Vs1.txt", t)
	strategyMaintainAirSuperiority()
	if len(assignedDrones) != 1 {
		t.Error("Wrong number of drones asigned:", len(assignedDrones))
	}
	destinations := make(map[point]bool, 1)
	for d, _ := range assignedDrones {
		if nextMove[d].x == 300 && nextMove[d].y == 300 {
			destinations[nextMove[d]] = true
		}
	}
	if len(destinations) != 1 {
		t.Error("Wrong number of different destinations:", len(destinations))
	}
}

//Tests method strategyMaintainAirSuperiority with one zone owned 2 Vs 2
func TestMaintainAirSuperiority2Vs2(t *testing.T) {
	setUpTestFromFile(FILE_TESTS_BASE+"maintainAirSuperiority\\inputOwned2Vs2.txt", t)

	strategyMaintainAirSuperiority()
	if len(assignedDrones) != 2 {
		t.Error("Wrong number of drones asigned:", len(assignedDrones))
	}
	destinations := make(map[point]bool, 2)
	for d, _ := range assignedDrones {
		if (nextMove[d].x == 300 && nextMove[d].y == 300) ||
			(nextMove[d].x == 300 && nextMove[d].y == 300) {
			destinations[nextMove[d]] = true
		}
	}
	if len(destinations) != 1 {
		t.Error("Wrong number of different destinations:", len(destinations))
	}
}

//Tests method strategyMaintainAirSuperiority with one zone owned 2 Vs 1 + 1
func TestMaintainAirSuperiority2Vs1Plus1(t *testing.T) {
	setUpTestFromFile(FILE_TESTS_BASE+"maintainAirSuperiority\\inputOwned2Vs1plus1.txt", t)

	strategyMaintainAirSuperiority()
	if len(assignedDrones) != 1 {
		t.Error("Wrong number of drones asigned:", len(assignedDrones))
	}
	destinations := make(map[point]bool, 1)
	for d, _ := range assignedDrones {
		if nextMove[d].x == 300 && nextMove[d].y == 300 {
			destinations[nextMove[d]] = true
		}
	}
	if len(destinations) != 1 {
		t.Error("Wrong number of different destinations:", len(destinations))
	}
}

//Tests method playerDronesInZone when there is no drone in the zone
func TestPlayerDronesInZoneZero(t *testing.T) {
	setUpTestFromFile(FILE_TESTS_BASE+"playerDronesInZone\\input0.txt", t)
	drones := playerDronesInZone(0, 2)

	if len(drones) != 0 {
		t.Error("Wrong number of drones in zone:", len(drones))
	}
}

//Tests method playerDronesInZone when there is one drone in the zone
func TestPlayerDronesInZoneOne(t *testing.T) {
	setUpTestFromFile(FILE_TESTS_BASE+"playerDronesInZone\\input1.txt", t)

	drones := playerDronesInZone(0, 2)

	if len(drones) != 1 {
		t.Error("Wrong number of drones in zone:", len(drones))
	}
	if _, isThere := drones[0]; !isThere {
		t.Error("Wrong drone in the zone", drones)
	}
}

//Tests method playerDronesInZone when there are two drones in the zone
func TestPlayerDronesInZoneTwo(t *testing.T) {
	setUpTestFromFile(FILE_TESTS_BASE+"playerDronesInZone\\input2.txt", t)
	drones := playerDronesInZone(0, 2)

	if len(drones) != 2 {
		t.Error("Wrong number of drones in zone:", len(drones))
	}
	if _, isThere := drones[0]; !isThere {
		t.Error("Wrong drone in the zone", drones)
	}
	if _, isThere := drones[1]; !isThere {
		t.Error("Wrong drone in the zone", drones)
	}
}

//Tests method strategyGoForUnguardedZones when there are no unguarded zones
func TestGoForUnguardedZonesZero(t *testing.T) {
	setUpTestFromFile(FILE_TESTS_BASE+"goForUnguarded\\input0.txt", t)
	strategyGoForUnguardedZones()

	if len(assignedDrones) != 0 {
		t.Error("Wrong number of drones asigned:", len(assignedDrones))
	}
}

//Tests method strategyGoForUnguardedZones when there is one unguarded zone and three available drones
func TestGoForUnguardedZonesOneForThree(t *testing.T) {
	setUpTestFromFile(FILE_TESTS_BASE+"goForUnguarded\\input1For3.txt", t)
	strategyGoForUnguardedZones()

	if len(assignedDrones) != 1 {
		t.Error("Wrong number of drones asigned:", len(assignedDrones))
	}
	if _, isThere := assignedDrones[1]; !isThere {
		t.Error("Wrong drone in the zone", assignedDrones)
	}
	if nextMove[1].x != 300 || nextMove[1].y != 300 {
		t.Error("Wrong move", nextMove[1])
	}
}

//Tests method strategyGoForUnguardedZones when there are two unguarded zones and one available drone
func TestGoForUnguardedZonesTwoForOne(t *testing.T) {
	setUpTestFromFile(FILE_TESTS_BASE+"goForUnguarded\\input2For1.txt", t)
	assignDestinationPoint(0, point{0, 0}, "I say so")
	assignDestinationPoint(1, point{0, 0}, "I say so")
	strategyGoForUnguardedZones()

	if len(assignedDrones) != 3 {
		t.Error("Wrong number of drones asigned:", len(assignedDrones))
	}
	if _, isThere := assignedDrones[2]; !isThere {
		t.Error("Wrong drone in the zone", assignedDrones)
	}
	if nextMove[2].x != 200 || nextMove[2].y != 200 {
		t.Error("Wrong move", nextMove[2])
	}
}

//Tests method strategyGoForUnguardedZones when there are two unguarded zones and two available drones
func TestGoForUnguardedZonesTwoForTwo(t *testing.T) {
	setUpTestFromFile(FILE_TESTS_BASE+"goForUnguarded\\input2For2.txt", t)
	assignDestinationPoint(0, point{0, 0}, "I say so")
	strategyGoForUnguardedZones()

	if len(assignedDrones) != 3 {
		t.Error("Wrong number of drones asigned:", len(assignedDrones))
	}
	if _, isThere := assignedDrones[0]; !isThere {
		t.Error("Wrong drone in the zone", assignedDrones)
	}
	if _, isThere := assignedDrones[1]; !isThere {
		t.Error("Wrong drone in the zone", assignedDrones)
	}
	if nextMove[1].x != 200 || nextMove[1].y != 200 {
		t.Error("Wrong move", nextMove[1])
	}
	if nextMove[2].x != 300 || nextMove[2].y != 300 {
		t.Error("Wrong move", nextMove[2])
	}
}

//Sets up the test reading current status from a certain file
func setUpTestFromFile(path string, t *testing.T) {
	if f, err := os.Open(path); err != nil {
		t.Error("Error opening input file", path)
	} else {
		inputReader = f
		defer f.Close()
	}
	readBoard()
	parseTurn()
	initializeTurnComputation()
}

//Tests method getCentroid
func TestGetCentroid(t *testing.T) {
	var testCases = []struct {
		in  []zone
		out point
	}{
		{[]zone{{point{500, 500}, -1}}, point{500, 500}},
		{[]zone{{point{400, 400}, -1}, {point{600, 600}, -1}}, point{500, 500}},
		{[]zone{{point{400, 400}, -1}, {point{600, 600}, -1}, {point{500, 200}, -1}}, point{500, 400}},
	}
	for i, testCase := range testCases {
		if result := getCentroid(testCase.in); testCase.out != result {
			t.Error("Error in item", i, "Got", result, "Expected", testCase.out, "Case:", testCase)
		}
	}
}

//Tests method defaultToCentroid
func TestDefaultToCentre(t *testing.T) {
	setUpTestFromFile(FILE_TESTS_BASE+"toCentre\\input.txt", t)
	assignDestinationPoint(0, point{0, 0}, "I say so")
	strategyDefaultToCentroid()

	if len(assignedDrones) != 3 {
		t.Error("Wrong number of drones asigned:", len(assignedDrones))
	}
	if _, isThere := assignedDrones[0]; !isThere {
		t.Error("Wrong drone in the zone", assignedDrones)
	}
	if _, isThere := assignedDrones[1]; !isThere {
		t.Error("Wrong drone in the zone", assignedDrones)
	}
	if _, isThere := assignedDrones[2]; !isThere {
		t.Error("Wrong drone in the zone", assignedDrones)
	}
	if nextMove[0].x != 0 || nextMove[0].y != 0 {
		t.Error("Wrong move", nextMove[0])
	}
	if nextMove[1].x != getCentroid(zones).x || nextMove[1].y != getCentroid(zones).x {
		t.Error("Wrong move", nextMove[1])
	}
	if nextMove[2].x != getCentroid(zones).x || nextMove[2].y != getCentroid(zones).x {
		t.Error("Wrong move", nextMove[2])
	}
}

//Tests method attackable:
//- Test 1
//  + Zone 0 is not attackable because it is unreclaimed
//  + If I set zone0 to the enemy, it is attackable
//  + If I set zone0 to myself, it is no longer attackable
// - Test2
//  + Zone 2 is not attackable because it contains three enemy drones and I only have three available free
//  + If I remove one of enemy's drones, it is attackable
//  + If I asign a drone to any destination, it is not attackable
//  + If I put all drones in the zone and unassigned, it is NOT attackable because I cannot gain air superiority
//  + If I remove one of enemy's drones, it is attackable
func TestAttackable(t *testing.T) {
	setUpTestFromFile(FILE_TESTS_BASE+"attackable\\input.txt", t)
	attackableTest1(t)
	attackableTest2(t)
}

//Described abobe
func attackableTest1(t *testing.T) {
	if !attackable(0) {
		t.Error("Zone 0 should be attackable because it should be unreclaimed.", zones[0])
	}
	zones[0].owner = whoami + 1
	if !attackable(0) {
		t.Error("Zone 0 should be attackable.", zones[0])
	}
	zones[0].owner = whoami
	if attackable(0) {
		t.Error("Zone 0 should NOT be attackable. It belongs to me", zones[0])
	}
}

//Described abobe
func attackableTest2(t *testing.T) {
	if attackable(2) {
		t.Error("Zone 2 should NOT be attackable. I cannot gain air superiority", status())
	}
	players[1].drones[0] = point{0, 0}
	if !attackable(2) {
		t.Error("Zone 2 should be attackable. I should be able to gain air superiority (even when all my drones are unassigned).", status())
	}
	assignDestinationZone(0, 2, "I say so") //should work even when destination is the zone I pretend to atack
	if attackable(2) {
		t.Error("Zone 2 should NOT be attackable. I have assigned one drone and only two should be available.", status())
	}
	players[whoami].drones[0] = zones[2].pos //If drone 0 is inside the zone, but will remain there, I can achieve air superiority
	if !attackable(2) {
		t.Error("Zone 2 should be attackable. I have assigned one drone and only two should be available, but assigned drone is to remain inside the zone.", status())
	}
	assignDestinationPoint(0, point{0, 0}, "I say so") //If drone 0 is inside the zone, but will go out (assigned to go out), I cannot gain air superiority
	if attackable(2) {
		t.Error("Zone 2 should NOT be attackable. I have assigned one drone to go out of the zone and only two should be available.", status())
	}
	initializeTurnComputation() //To free up all drones
	for pId, _ := range players {
		for dId, _ := range players[pId].drones {
			players[pId].drones[dId] = zones[2].pos
		}
	}
	for dId, _ := range players[whoami].drones {
		assignDestinationZone(dId, 2, "I say so")
	}
	if attackable(2) {
		t.Error("Zone 2 should NOT be attackable. I cannot gain air superiority (even when all my drones are inside the zone", status())
	}
	players[1].drones[0] = point{0, 0}
	if !attackable(2) {
		t.Error("Zone 2 should be attackable. I have enough drones inside the zone", status())
	}
}

//Tests method nearestFreeOwnDrone
func TestNearestFreeOwnDrone(t *testing.T) {
	setUpTestFromFile(FILE_TESTS_BASE+"NearestFreeOwnDrone\\input.txt", t)
	if result := nearestFreeOwnDrone(zones[0].pos); result != 0 {
		t.Error("Nearest unasigned drone:", result)
	}
	assignDestinationPoint(0, point{0, 0}, "I say so")
	if result := nearestFreeOwnDrone(zones[0].pos); result != 1 {
		t.Error("Nearest unasigned drone:", result)
	}
	assignDestinationPoint(1, point{0, 0}, "I say so")
	if result := nearestFreeOwnDrone(zones[0].pos); result != 2 {
		t.Error("Nearest unasigned drone:", result)
	}
	assignDestinationPoint(2, point{0, 0}, "I say so")
	if result := nearestFreeOwnDrone(zones[0].pos); result != -1 {
		t.Error("Nearest unasigned drone:", result)
	}
}

//Tests method nearestOwnDroneToGoFromSet
func TestNearestOwnDroneFromSet(t *testing.T) {
	setUpTestFromFile(FILE_TESTS_BASE+"NearestFreeOwnDrone\\input.txt", t)
	set := make(map[int]bool, 2)
	set[0] = true
	set[2] = true
	set[1] = true
	if result := nearestOwnDroneToGoFromSet(zones[0].pos, set); result != 0 {
		t.Error("Nearest unasigned drone:", result)
	}
	assignDestinationZone(0, 0, "I say so")
	if result := nearestOwnDroneToGoFromSet(zones[0].pos, set); result != 0 {
		t.Error("Nearest unasigned drone:", result)
	}
	assignDestinationPoint(0, point{0, 0}, "I say so")
	if result := nearestOwnDroneToGoFromSet(zones[0].pos, set); result != 1 {
		t.Error("Nearest unasigned drone:", result)
	}
	delete(set, 0)
	if result := nearestOwnDroneToGoFromSet(zones[0].pos, set); result != 1 {
		t.Error("Nearest unasigned drone:", result)
	}
	assignDestinationPoint(1, point{0, 0}, "I say so")
	if result := nearestOwnDroneToGoFromSet(zones[0].pos, set); result != 2 {
		t.Error("Nearest unasigned drone:", result)
	}
	assignDestinationPoint(2, point{0, 0}, "I say so")
	if result := nearestOwnDroneToGoFromSet(zones[0].pos, set); result != -1 {
		t.Error("Nearest unasigned drone:", result)
	}
}

//Tests method isAssigned
func TestIsAssigned(t *testing.T) {
	setUpTestFromFile(FILE_TESTS_BASE+"toCentre\\input.txt", t)
	if isAssigned(0) {
		t.Error("Drone zero should NOT be assigned")
	}
	assignDestinationPoint(0, point{0, 0}, "I say so")
	if !isAssigned(0) {
		t.Error("Drone zero should be assigned")
	}
}

//Tests method numAssignedDrones
func TestNumAssignedDrones(t *testing.T) {
	setUpTestFromFile(FILE_TESTS_BASE+"toCentre\\input.txt", t)
	if numAssignedDrones() != 0 {
		t.Error("All drones should be free")
	}
	assignDestinationPoint(0, point{0, 0}, "I say so")
	if numAssignedDrones() != 1 {
		t.Error("All drones but one should be free")
	}
}
