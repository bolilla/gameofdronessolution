//Participating Game of Drones by CodinGame
package main

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"time"
)

const (
	DEBUG          = true   //True iff traces are activated
	ZONE_RADIUS    = 100.0  //Radius of the zones
	DRONE_MOVEMENT = 100.0  //Maximum movement of a drone in a turn
	MAX_DISTANCE   = 44     //Number of turns to cross the board
	BOARD_DIAGONAL = 4387.0 //Length of the diagonal of the board (maximum distance between two points)
	UNRECLAIMED    = -1     //Owner of unreclaimed zones
)

var ( //board-related variables
	inputReader        io.Reader //Where the information is read (os.Stdin for play, a file for testing)
	numPlayers         int       //Number of players in the game
	numZones           int       //Number of zones in the game
	numDronesPerplayer int       //Number of drones each player has
	whoami             int       //index of my player in the array of players
	players            []player  //all the player of drones. Array index = player's ID
	zones              []zone    //all game zones
	centroid           point     //Centroid of the zones
)

var ( //turn-related variables
	distances      [][][]int    //Distances for each of the players, for each of the drones to each of the zones
	nextMove       []point      //destination for each of my drones
	assignedDrones map[int]bool //Drones that have a destination asigned in this turn
)

type point struct {
	x, y int
}

type player struct {
	score  int     //Number of points of this player
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

//Prints the movements of own drones
func play() {
	initializeTurnComputation()
	maintainAirSuperiority()
	colonizeTheUnexplored()
	goForUnguardedZones()
	attackGuardedZones()
	defaultToCentroid()
	//defaultEachDroneToNearestZone()
	for _, m := range nextMove {
		fmt.Println(m.x, m.y)
	}
}

type attack struct {
	target   int          //ID of the zone to attack
	distance int          // Number of turns the farthest drone must go
	force    map[int]bool //Set of drones that will make the attack
}

//Implements sort.Interface
type attackSorter []attack

//Necessary to implement sort.Interface
func (as attackSorter) Less(i, j int) bool {
	if len(as[i].force) < len(as[j].force) {
		return true
	}
	return as[i].distance < as[j].distance
}

//Necessary to implement sort.Interface
func (as attackSorter) Swap(i, j int) {
	as[i], as[j] = as[j], as[i]
}

//Necessary to implement sort.Interface
func (as attackSorter) Len() int {
	return len(as)
}

//Calculates the movements for unasigned drones based on the following strategy:
//- while there is an attackable zone
//  + For all attackable zones
//    * Define attack
//  + Choose best attack
func attackGuardedZones() {
	for numAttackableZones := 2; numAttackableZones > 1; {
		numAttackableZones = 0
		attacks := make([]attack, 0, numZones)
		for zId, _ := range zones {
			if attackable(zId) {
				numAttackableZones += 1
				attacks = append(attacks, defineAttack(zId))
			}
		}
		if numAttackableZones > 0 {
			sort.Sort(attackSorter(attacks))
			a := attacks[0]
			for dId, _ := range a.force {
				assignDestination(dId, zones[a.target].pos)
			}
		}
	}
}

//Defines the attack over a zone
func defineAttack(zId int) (result attack) {
	result.target = zId
	numDronesOwner := len(playerDronesInZone(zones[zId].owner, zId))
	numDronesMe := len(playerDronesInZone(whoami, zId))
	numFreeDrones := numDronesPerplayer - len(assignedDrones)
	result.force = make(map[int]bool, numDronesOwner+1)
	for required := 0; required <= numDronesOwner-numDronesMe+numFreeDrones; required += 1 {
		nearest := nearestFreeOwnDrone(zones[zId].pos)
		result.force[nearest] = true
	}
	for dId, _ := range result.force {
		if thisDistance := turnBasedDistance(players[whoami].drones[dId], zones[zId].pos); thisDistance > result.distance {
			result.distance = thisDistance
		}
	}
	return result
}

//Calculates the movements for unasigned drones based on the following strategy:
//- If there is an unguarded zone, nearest drone goes to take it
func goForUnguardedZones() {
	for zId, z := range zones {
		if z.owner != UNRECLAIMED && z.owner != whoami && len(playerDronesInZone(z.owner, zId)) == 0 {
			dId := nearestFreeOwnDrone(z.pos)
			if dId >= 0 {
				assignDestination(dId, z.pos)
			}
		}
	}
}

//Clears old turn's data and calculates this turn key information
func initializeTurnComputation() {
	calculateDistances()
	nextMove = make([]point, numDronesPerplayer)
	assignedDrones = make(map[int]bool, numDronesPerplayer)
}

//Calculates the movements for unasigned drones based on the following strategy:
//- If (1+ drone is inside an owned zone AND there are enemies in the same zone)
//    air superiority cannot be lost (cannot abandon zone and leave air superiority to the oponent)
func maintainAirSuperiority() {
	for zId, z := range zones {
		if z.owner == whoami {
			myDrones := playerDronesInZone(whoami, zId)
			numHostiles := mostDronesBySingleOponentInZone(zId)
			i := 0
			for dId, _ := range myDrones {
				if i >= numHostiles {
					break
				}
				assignDestination(dId, players[whoami].drones[dId])
				i += 1
			}
		}
	}
}

//Returns the number of drones of the oponent who has most oponents in the given zone
func mostDronesBySingleOponentInZone(zId int) int {
	result := 0
	for pId, _ := range players {
		if pId == whoami {
			continue
		}
		if currentPlayerDronesInZone := len(playerDronesInZone(pId, zId)); currentPlayerDronesInZone > result {
			result = currentPlayerDronesInZone
		}
	}
	return result
}

//Returns a set of ids of the drones of given player that are inside given zone
func playerDronesInZone(pId, zId int) map[int]bool {
	result := make(map[int]bool)
	for dId, d := range players[pId].drones {
		if turnBasedDistance(zones[zId].pos, d) == 0 {
			result[dId] = true
		}
	}
	return result
}

//Calculates the movements for unasigned drones based on the following strategy:
//- If there is an unreclaimed zone:
//  + Drones inside that zone stay put
//  + Nearest drone (outside the zone) goes for it
func colonizeTheUnexplored() {
	for zId, _ := range unreclaimedZones() {
		for dId, d := range players[whoami].drones {
			if turnBasedDistance(d, zones[zId].pos) == 0 {
				assignDestination(dId, zones[zId].pos)
			}
		}
		if dId := nearestFreeOwnDrone(zones[zId].pos); dId != -1 {
			assignDestination(dId, zones[zId].pos)
		}
	}
}

//Returns the Id of the nearest drone I control (and has not been sent to other duties) to the given point
func nearestFreeOwnDrone(p point) int {
	minDist := BOARD_DIAGONAL
	bestDrone := -1
	for dId, d := range players[whoami].drones {
		if _, isAssigned := assignedDrones[dId]; !isAssigned {
			if currentDistance := euclideanDistance(d, p); currentDistance <= minDist {
				minDist = currentDistance
				bestDrone = dId
			}
		}
	}
	return bestDrone
}

//Returns the zones that remain unreclaimed
func unreclaimedZones() map[int]bool {
	result := make(map[int]bool, numZones)
	for i, z := range zones {
		if z.owner == UNRECLAIMED {
			result[i] = true
		}
	}
	return result
}

//Calculates the movements for the remaining drones based on the following strategy:
//- Each remaining drone moves to the centroid of the board
func defaultToCentroid() {
	for dId := 0; dId < numDronesPerplayer; dId += 1 {
		if _, isDroneAsigned := assignedDrones[dId]; isDroneAsigned {
			continue
		}
		assignDestination(dId, centroid)
	}
}

//Calculates the movements for the remaining drones based on the following strategy:
//- Each remaining drone moves to the centre of its nearest zone
func defaultToNearestZone() {
	for dId := 0; dId < numDronesPerplayer; dId += 1 {
		if _, isDroneAsigned := assignedDrones[dId]; isDroneAsigned {
			continue
		}
		minDist := MAX_DISTANCE
		bestZone := -1
		for zId := 0; zId < numZones; zId += 1 {
			if distances[whoami][dId][zId] <= minDist {
				minDist = distances[whoami][dId][zId]
				bestZone = zId
			}
		}
		assignDestination(dId, zones[bestZone].pos)
	}
}

//Asigns a drone to a destination
func assignDestination(dId int, p point) {
	nextMove[dId] = p
	assignedDrones[dId] = true
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

//Calculates the number of turns that it would take a drone to move from pointA to pointB
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

//Returns the centroid of the zones. i.e. The point that is at the same minimum distance from all zones
func getCentroid(zs []zone) (result point) {
	for _, z := range zs {
		result.x += z.pos.x
		result.y += z.pos.y
	}
	result.x /= len(zs)
	result.y /= len(zs)
	return result
}

//Returns whether given zone can be attacked with remaining forces
func attackable(zId int) bool {
	if zones[zId].owner == whoami || zones[zId].owner == UNRECLAIMED {
		return false
	}
	numDronesOwner := len(playerDronesInZone(zones[zId].owner, zId))
	numDronesMe := len(playerDronesInZone(whoami, zId))
	numFreeDrones := numDronesPerplayer - len(assignedDrones)
	return numDronesOwner < numDronesMe+numFreeDrones
}

//Reads the game initialization information
func readBoard() {
	fmt.Fscanf(inputReader, "%d %d %d %d\n", &numPlayers, &whoami, &numDronesPerplayer, &numZones)
	players = make([]player, numPlayers)
	for i, _ := range players {
		players[i].drones = make([]point, numDronesPerplayer)
	}
	zones = make([]zone, numZones)
	for i, _ := range zones {
		zones[i] = newZone()
		fmt.Fscanf(inputReader, "%d %d\n", &zones[i].pos.x, &zones[i].pos.y)
	}
	distances = make([][][]int, numPlayers)
	for pId := 0; pId < numPlayers; pId += 1 {
		distances[pId] = make([][]int, numDronesPerplayer)
		for dId := 0; dId < numDronesPerplayer; dId += 1 {
			distances[pId][dId] = make([]int, numZones)
		}
	}
	centroid = getCentroid(zones)
}

//Reads the information of a turn
func parseTurn() bool {
	for i, _ := range zones {
		_, err := fmt.Fscanf(inputReader, "%d\n", &zones[i].owner)
		if err != nil {
			fmt.Println("Error reading turn zones owners:", err)
			return false
		}
		if zones[i].owner >= 0 {
			players[zones[i].owner].score += 1
		}
	}

	for i, _ := range players {
		for j, _ := range players[i].drones {
			_, err := fmt.Fscanf(inputReader, "%d %d\n", &players[i].drones[j].x, &players[i].drones[j].y)
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
	inputReader = os.Stdin
	letTheGameBegin() //..hear the starting gun
}

//Plays the came of drones
func letTheGameBegin() {

	tFrom := time.Now()
	readBoard()
	debug("Initial status:", status())
	debug(fmt.Sprintf("Initialization computation time: %v microseconds", time.Now().Sub(tFrom).Nanoseconds()/1000))
	for tFrom = time.Now(); parseTurn(); tFrom = time.Now() {
		debug("Current status:", status())
		play()
		debug(fmt.Sprintf("Turn computation time: %v microseconds", time.Now().Sub(tFrom).Nanoseconds()/1000))
	}
	debug("End status:", status())
}

//Returns the status of the play if debug is enabled
func status() string {
	var result bytes.Buffer
	if DEBUG {
		result.Write([]byte("Players:\n"))
		for pId, p := range players {
			var numZonesPlayer int
			for _, z := range zones {
				if z.owner == pId {
					numZonesPlayer += 1
				}
			}
			var playerName string
			if pId == whoami {
				playerName = "(ME)"
			} else {
				playerName = "    "
			}
			result.Write([]byte(fmt.Sprintf("  %d%s- score: %d numZones: %d Drones: ",
				pId, playerName, p.score, numZonesPlayer)))
			result.Write([]byte("["))
			for dId, d := range p.drones {
				if _, isAssigned := assignedDrones[dId]; isAssigned && pId == whoami {
					result.Write([]byte(fmt.Sprintf("%v* ", d)))
				} else {
					result.Write([]byte(fmt.Sprintf("%v  ", d)))
				}
			}
			result.Write([]byte("]\n"))
		}
		result.Write([]byte("Zones:\n"))
		for zId, z := range zones {
			result.Write([]byte(fmt.Sprintf("  %d - owner: %d location: %v\n", zId, z.owner, z.pos)))
		}
	}
	return result.String()
}
func debug(x ...interface{}) {
	if DEBUG {
		fmt.Fprintln(os.Stderr, x)
	}
}

/*
IDEAS:
- Calculate "centroid" of the board based on the zones' locations. Move drones to the position in the zone neares to the center of the board
- Different strategies depending on whether I am winning (based on actual score of all players, zones under my control and remaining turns)
- Take into account oponents' possible movements
- Take into account oponents' drones' distances to owned zones
- should air superiority include drones at distance 1?
- Last strategy: go for the "centroid" (the center of all zones)
- Decide nearest drone per euclidean distance instead of per turns?
*/
