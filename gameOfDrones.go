//Participating Game of Drones by CodinGame
package main

import (
	"fmt"
	"io"
	"math"
	"os"
)

const (
	DEBUG          = false //True iff traces are activated
	ZONE_RADIUS    = 100.0 //Radius of the zones
	DRONE_MOVEMENT = 100.0 //Maximum movement of a drone in a turn
	MAX_DISTANCE   = 44    //Number of turns to cross the board
)

var (
	numPlayers         int       //Number of players in the game
	numZones           int       //Number of zones in the game
	players            []player  //all the player of drones. Array index = player's ID
	numDronesPerplayer int       //Number of drones each player has
	whoami             int       //index of my player in the array of players
	zones              []zone    //all game zones
	distances          [][][]int //Distances for each of the players, for each of the drones to each of the zones
)

//Prints the movements of own drones
func play() {
	// Do not move any drone
	for _, m := range decideMovement() {
		fmt.Println(m.x, m.y)
	}
}

//Returns the movement for each of the drones
func decideMovement() []point {
	return stratEachDroneToNearestZone()
}

//Returns the movements based on the following strategy:
//- Each drone moves to its nearest zone
//- [NO] Only a drone per zone
//- Move to the center of the zone
func stratEachDroneToNearestZone() []point {
	result := make([]point, numDronesPerplayer)
	calculateDistances()
	//	zonesWithDrone := make(map[int]bool, numZones)
	for dId := 0; dId < numDronesPerplayer; dId += 1 {
		minDist := MAX_DISTANCE
		bestZone := -1
		for zId := 0; zId < numZones; zId += 1 {
			//if _, occupied := zonesWithDrone[zId]; !occupied && distances[whoami][dId][zId] <= minDist {
			if distances[whoami][dId][zId] <= minDist {
				minDist = distances[whoami][dId][zId]
				bestZone = zId
			}
		}
		result[dId] = zones[bestZone].pos
		//zonesWithDrone[bestZone] = true
	}

	return result
}

//Calculates the distances from each of my drones to each of the zones' centres
func calculateDistances() {
	for pId, _ := range players {
		for dId, d := range players[pId].drones {
			for zId, z := range zones {
				distances[pId][dId][zId] = turnBasedDistance(d, z.pos)
			}
		}
	}
}

//Calculates the number of turns that it would take to each drone to reach each zone
func turnBasedDistance(pointA, pointB point) int {
	euc := euclideanDistance(pointA, pointB)
	if euc < ZONE_RADIUS {
		return 0
	}
	return int(math.Ceil((euc - (ZONE_RADIUS)) / DRONE_MOVEMENT))
}

//Returns the euclidean distance between two points
func euclideanDistance(pointA, pointB point) float64 {
	return math.Floor(math.Sqrt((float64(pointB.x-pointA.x) * float64(pointB.x-pointA.x)) +
		(float64(pointB.y-pointA.y) * float64(pointB.y-pointA.y))))
}

type point struct {
	x, y int
}

type player struct {
	drones []point //position of each drone
}

//Contains the center of the zone and current owner
type zone struct {
	pos   point
	owner int
}

func newZone() zone {
	return zone{point{-1, -1}, -1}
}

//Reads the game initialization information
func readBoard(in io.Reader) {
	fmt.Scanf("%d %d %d %d\n", &numPlayers, &whoami, &numDronesPerplayer, &numZones)
	players = make([]player, numPlayers)
	for i, _ := range players {
		players[i].drones = make([]point, numDronesPerplayer)
	}
	zones = make([]zone, numZones)
	for i, _ := range zones {
		zones[i] = newZone()
		fmt.Fscanf(in, "%d %d\n", &zones[i].pos.x, &zones[i].pos.y)
	}
	distances = make([][][]int, numPlayers)
	for pId := 0; pId < numPlayers; pId += 1 {
		distances[pId] = make([][]int, numDronesPerplayer)
		for dId := 0; dId < numDronesPerplayer; dId += 1 {
			distances[pId][dId] = make([]int, numZones)
		}
	}
}

//Reads the information of a turn
func parseTurn(in io.Reader) bool {
	for i, _ := range zones {
		_, err := fmt.Fscanf(in, "%d\n", &zones[i].owner)
		if err != nil {
			fmt.Println("Error reading turn zones owners:", err)
			return false
		}
	}

	for i, _ := range players {
		for j, _ := range players[i].drones {
			_, err := fmt.Fscanf(in, "%d %d\n", &players[i].drones[j].x, &players[i].drones[j].y)
			if err != nil {
				fmt.Println("Error reading turn drones:", err)
				return false
			}
		}
	}
	return true
}

//Unleashes the beast
func main() {
	readBoard(os.Stdin)
	debug("Initial status:", players, whoami, zones)
	for parseTurn(os.Stdin) {
		debug("Current status:", players, whoami, zones)
		play()
	}
}

func debug(x ...interface{}) {
	if DEBUG {
		fmt.Println(x)
	}
}
