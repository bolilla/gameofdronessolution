//Participating Game of Drones by CodinGame
package main

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"runtime/debug"
	"runtime/pprof"
	"sort"
	"time"
)

/************************************************************************************** CONSTANTS AND VARIABLES BEGIN */
const (
	DEBUG              = false  //True iff traces are activated
	PROFILING          = false  //True iff profiling must be activated
	ZONE_RADIUS        = 100.0  //Radius of the zones
	DRONE_MOVEMENT     = 100.0  //Maximum movement of a drone in a turn
	MAX_DISTANCE       = 44     //Number of turns to cross the board
	BOARD_DIAGONAL     = 4387.0 //Length of the diagonal of the board (maximum distance between two points)
	UNRECLAIMED        = -1     //Owner of unreclaimed zones
	NUM_TURNS_TO_CHECK = 5      //Number of turns to try to look into the future
)

const PROFILE_PATH = "C:\\Users\\borja\\programacion\\codinggame\\GameOfDronesSolution\\profile.pprof" //Path to the file that stores the profiling information

var ( //board-related variables
	realExecution      bool      // True iff the execution is made via "main" function
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
	distances    [][][]int //Distances for each of the players, for each of the drones to each of the zones
	nextMove     []point   //destination for each of my drones
	availability struct {
		numAvailables int   //number of drones with some degree of availability
		drones        []int //Number of turns the drone can be traveling
	}
	//assignedDrones map[int]bool //Drones that have a destination asigned in this turn
)

/* CONSTANTS AND VARIABLES END ********************************************************************* DATA TYPES BEGIN */

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

type attack struct {
	target   int          //ID of the zone to attack
	distance int          //Number of turns the farthest drone must travel
	force    map[int]bool //Set of drones that will make the attack
}

//updates the length of the attack
func (a *attack) calculateLength() {
	a.distance = 0
	for dId, _ := range a.force {
		if calculatedDistance := turnBasedDistance(players[whoami].drones[dId], zones[a.target].pos); calculatedDistance > a.distance {
			a.distance = calculatedDistance
		}
	}
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

/* DATA TYPES END ********************************************************************* STATEGIES BEGIN */

//Calculates the movements for unasigned drones based on the following strategy:
//- while there is an attackable zone
//  + For all attackable zones
//    * Define attack
//  + Choose best attack
func strategyAttack() {
	attackableZones := make(map[int]bool, numZones)
	for zId, _ := range zones {
		attackableZones[zId] = true
	}
	for len(attackableZones) > 0 {
		trace("availability", availability)
		attacks := make([]attack, 0, numZones)
		for zId, _ := range attackableZones {
			if a, attackable := bestAttackToZone(zId); attackable {
				trace("Adding attack", a)
				attacks = append(attacks, a)
			} else {
				delete(attackableZones, zId)
			}
		}
		if len(attacks) > 0 {
			sort.Sort(attackSorter(attacks))
			for dId, _ := range attacks[0].force {
				assignDestinationZone(dId, attacks[0].target, "Zone must be ours!!!")
			}
			delete(attackableZones, attacks[0].target)
		}
	}
	/*
		attackableZones := make(map[int]bool, numZones)
		for zId, _ := range zones {
			if attackable(zId) {
				attackableZones[zId] = true
			}
		}
		for len(attackableZones) > 1 {
			attacks := make([]attack, 0, numZones)
			for zId, _ := range attackableZones {
				if attackable(zId) {
					attacks = append(attacks, defineAttack(zId))
				} else {
					delete(attackableZones, zId)
				}
			}
			if len(attacks) > 0 {
				sort.Sort(attackSorter(attacks))
				a := attacks[0]
				for dId, _ := range a.force {
					assignDestinationZone(dId, a.target, "Zone must be ours!!!")
				}
				delete(attackableZones, a.target)
			}
		}
	*/
}

//Calculates the movements for unasigned drones based on the following strategy:
//- If (1+ drone is inside an owned zone AND there are enemies in the same zone)
//    air superiority cannot be lost (cannot abandon zone and leave air superiority to the oponent)
func strategyMaintainAirSuperiority() {
	for zId, z := range zones {
		if z.owner == whoami {
			myDrones := playerDronesNearZone(whoami, zId, 0)
			numHostiles := mostDronesBySingleOponentInZone(zId)
			i := 0
			for dId, _ := range myDrones {
				if !isAssigned(dId) {
					if i >= numHostiles {
						break
					}
					assignDestinationZone(dId, zId, "Zone air supperiority must be maintained")
					i += 1
				}
			}
		}
	}
}

//Calculates the movements for the remaining drones based on the following strategy:
//- Each remaining drone moves to the centre of its nearest zone
func strategyDefaultToNearestZone() {
	for dId := 0; dId < numDronesPerplayer; dId += 1 {
		if isAssigned(dId) {
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
		assignDestinationZone(dId, bestZone, "It is my nearest zone")
	}
}

//Calculates the movements for the remaining drones based on the following strategy:
//- Each remaining drone moves to the centroid of the board
func strategyDefaultToCentroid() {
	for dId := 0; dId < numDronesPerplayer; dId += 1 {
		if isAssigned(dId) {
			continue
		}
		assignDestinationPointNoisy(dId, centroid, "Going to the centroid to support my comrades")
	}
}

/* STRATEGIES END **************************************************************************** ATTACK UTILITIES BEGIN */

//Returns the best available strategy to attack base zId. If it cannot be attacked, isAttackable is false
func bestAttackToZone(zId int) (result attack, isAttackable bool) {
	trace("Analyzing attack to", zId)
	if zones[zId].owner != whoami {
		result.target = zId
		var ordersMustBeGiven bool
		trace("Zone", zId, "is not mine. Let's try to get it")
		result.force = make(map[int]bool, numDronesPerplayer)
		dist := 0
		enemies := maxEnemiesNearZone(zId, dist)
		for len(result.force) <= enemies && dist <= MAX_DISTANCE {
			trace("Iterating because we still do not have enoug forces at distance", dist, ":", enemies, "Vs", len(result.force))
			ownPossibilities := playerDronesNearZone(whoami, zId, dist)
			droneForTheAttack, mustMove := nearestOwnDroneToGoFromSet(zones[zId].pos, ownPossibilities)
			ordersMustBeGiven = ordersMustBeGiven || mustMove
			for droneForTheAttack != -1 && len(result.force) <= enemies {
				result.force[droneForTheAttack] = true
				delete(ownPossibilities, droneForTheAttack)
				droneForTheAttack, mustMove = nearestOwnDroneToGoFromSet(zones[zId].pos, ownPossibilities)
				ordersMustBeGiven = ordersMustBeGiven || mustMove
			}
			if enemies < len(result.force) {
				break
			}
			dist++
			enemies = maxEnemiesNearZone(zId, dist)
		}
		if dist > MAX_DISTANCE || !ordersMustBeGiven {
			result.force = make(map[int]bool)
		}
	}
	result.calculateLength()
	return result, len(result.force) > 0
}

/*
//Defines the attack over a zone
func defineAttack(zId int) (result attack) {
	result.target = zId
	var numDronesOwner int
	if zones[zId].owner != UNRECLAIMED {
		numDronesOwner = len(playerDronesNearZone(zones[zId].owner, zId, 0))
	} else {
		numDronesOwner = mostDronesBySingleOponentInZone(zId)
	}
	result.force = make(map[int]bool, numDronesOwner+1)
	availableDrones := make(map[int]bool, numDronesPerplayer-numAssignedDrones())
	for dId, _ := range players[whoami].drones {
		if isAssigned(dId) || turnBasedDistance(zones[zId].pos, nextMove[dId]) > 0 {
			availableDrones[dId] = true
		}
	}
	for required := 0; required < numDronesOwner+1; required += 1 {
		nearest := nearestOwnDroneToGoFromSet(zones[zId].pos, availableDrones)
		delete(availableDrones, nearest)
		result.force[nearest] = true
	}
	for dId, _ := range result.force {
		if thisDistance := turnBasedDistance(players[whoami].drones[dId], zones[zId].pos); thisDistance > result.distance {
			result.distance = thisDistance
		}
	}
	return result
}

//Returns whether given zone can be attacked with remaining forces
func attackable(zId int) bool {
	if zones[zId].owner == whoami {
		return false
	}
	var numDronesOwner int
	if zones[zId].owner != UNRECLAIMED {
		numDronesOwner = len(playerDronesNearZone(zones[zId].owner, zId, NUM_TURNS_TO_CHECK-1))
	} else {
		numDronesOwner = mostDronesBySingleOponentInZone(zId)
	}
	myDronesNearZone := playerDronesNearZone(whoami, zId, NUM_TURNS_TO_CHECK)
	for dId, _ := range myDronesNearZone {
		if turnBasedDistance(zones[zId].pos, nextMove[dId]) > availableDistance(dId) {
			delete(myDronesNearZone, dId)
		}
	}
	numDronesMe := len(myDronesNearZone)
	numFreeDrones := numDronesPerplayer - numAssignedDrones()
	return numDronesOwner < numDronesMe+numFreeDrones
}
*/
/* ATTACK UTILITIES END ********************************************************************* GENERAL UTILITIES BEGIN */
//Clears old turn's data and calculates this turn key information
func initializeTurnComputation() {
	calculateDistances()
	nextMove = make([]point, numDronesPerplayer)
	availability.numAvailables = numDronesPerplayer
	availability.drones = make([]int, numDronesPerplayer)
	for i, _ := range availability.drones {
		availability.drones[i] = MAX_DISTANCE
	}
	trace("Availability XXX", availability)
}

//Calculates the movements of the drones can make without geopardizing the zone they protect
func calculateDonesAvailableDistances() {
	for zId, _ := range zones {
		if zones[zId].owner == whoami {
			for i := 0; i < MAX_DISTANCE; i += 1 {
				myDronesSet := playerDronesNearZone(whoami, zId, i)
				/*
					myDrones := make([]bool, numDronesPerplayer)
					for dId, _ := range myDronesSet {
						myDrones[dId] = true
					}
				*/
				numEnemies := maxEnemiesNearZone(zId, i)
				if numEnemies > len(myDronesSet) {
					break //Nothing to do. Air superiority is lost at this distance.
				}
				numDronesLocked := 0
				for dId, _ := range myDronesSet {
					if i == 0 {
						assignDestinationZone(dId, zId, "Too risky to move")
					} else {
						setAvailableDistance(dId, i-1)
					}
					numDronesLocked++
					if numDronesLocked > numEnemies { //enoug drones locked for this menace
						break
					}
				}
			}
		}
	}
	trace("availability xxx", availability)
}

//Calculates the maximum number of foes from the same enemy at given distance of given zone
func maxEnemiesNearZone(zId, dist int) (result int) {
	for pId, _ := range players {
		if pId == whoami {
			continue
		}
		if num := len(playerDronesNearZone(pId, zId, dist)); num > result {
			result = num
		}
	}
	return
}

//Returns the number of drones of the oponent who has most oponents in the given zone
func mostDronesBySingleOponentInZone(zId int) int {
	result := 0
	for pId, _ := range players {
		if pId == whoami {
			continue
		}
		if currentPlayerDronesInZone := len(playerDronesNearZone(pId, zId, 0)); currentPlayerDronesInZone > result {
			result = currentPlayerDronesInZone
		}
	}
	return result
}

//Returns a set of ids of the drones of given player that are inside given distance of given zone
func playerDronesNearZone(pId, zId, dist int) map[int]bool {
	result := make(map[int]bool)
	for dId, d := range players[pId].drones {
		if turnBasedDistance(zones[zId].pos, d) <= dist {
			result[dId] = true
		}
	}
	return result
}

//Returns the Id of the nearest drone from the set of drones suplied.
//- The drone is free to do the movement: returns the drone id and true
//- The drone is inside the zone and assigned to remain still: returns the drone id and false
//- There is no suitable drone: returns -1 and false
func nearestOwnDroneToGoFromSet(p point, set map[int]bool) (int, bool) {
	minDist := BOARD_DIAGONAL
	bestDrone := -1
	for dId, _ := range set {
		if isAssigned(dId) && turnBasedDistance(p, players[whoami].drones[dId]) == 0 && turnBasedDistance(nextMove[dId], p) == 0 {
			return dId, false
		}
		if currentDistance := euclideanDistance(players[whoami].drones[dId], p); currentDistance <= minDist && availableDistance(dId) >= turnBasedDistance(players[whoami].drones[dId], p) && !isAssigned(dId) {
			minDist = currentDistance
			bestDrone = dId
		}
	}
	return bestDrone, bestDrone >= 0
}

//Returns the Id of the nearest drone I control (and has not been sent to other duties) to the given point
func nearestFreeOwnDrone(p point) int {
	minDist := BOARD_DIAGONAL
	bestDrone := -1
	for dId, d := range players[whoami].drones {
		if availableDistance(dId) >= turnBasedDistance(d, p) {
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

//Asigns a drone to a zone
func assignDestinationZone(dId, zId int, reason string) {
	turnInfo("Moving drone", dId, "to zone", zId, "because", reason)
	assignDestinationPoint(dId, zones[zId].pos)
}

//Assigns a drone to a point in the map and says so into the output
func assignDestinationPointNoisy(dId int, p point, reason string) {
	turnInfo("Moving drone", dId, "to point", p, "because", reason)
	assignDestinationPoint(dId, p)
}

//Assigns a drone to a point in the map
func assignDestinationPoint(dId int, p point) {
	trace("Assigning destination of drone", dId, ":", p)
	availability.drones[dId] = 0
	nextMove[dId] = p
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
	result := math.Floor(math.Sqrt((float64(pointB.x-pointA.x) * float64(pointB.x-pointA.x)) +
		(float64(pointB.y-pointA.y) * float64(pointB.y-pointA.y))))
	return result
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

//returns true iff given drone is already assigned to a destination
func isAssigned(dId int) bool {
	return availableDistance(dId) == 0
}

//Returns the number of assigned drones
func numAssignedDrones() (result int) {
	for _, avail := range availability.drones {
		if avail == 0 {
			result++
		}
	}
	return result
}

//Returns the available distance for the drone. i.e. the distance the drone can safely move
func availableDistance(dId int) int {
	return availability.drones[dId]
}

//Sets the distance the drone can safely fly
func setAvailableDistance(dId, dist int) {
	if availability.drones[dId] > dist {
		availability.drones[dId] = dist
	}
}

/* GENERAL UTILITIES END   *********************************************** INPUT PARSING - RELATED OPERATIONS BEGIN ***/
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

/* INPUT PARSING - RELATED OPERATIONS END ********************************TURN BEGIN/END - RELATED OPERATIONS BEGIN***/

//Unleashes the beast
func main() {
	realExecution = true
	if PROFILING {
		f, err := os.Create(PROFILE_PATH)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "\n\n\nRecovered while panicking. Status:\n", importableStatus(), "\n\n\n")
			fmt.Fprint(os.Stderr, string(debug.Stack()))
		}
	}()

	inputReader = os.Stdin
	letTheGameBegin() //..hear the starting gun
}

//Plays the came of drones
func letTheGameBegin() {

	tFrom := time.Now()
	readBoard()
	turnInfo("Initial status:", status())
	turnInfo(fmt.Sprintf("Initialization computation time: %v microseconds", time.Now().Sub(tFrom).Nanoseconds()/1000))
	for tFrom = time.Now(); parseTurn(); tFrom = time.Now() {
		turnInfo("XXX")
		turnInfo(importableStatus())
		turnInfo("XXX")
		turnInfo("Current status:", status())
		play()
		turnInfo(fmt.Sprintf("Turn computation time: %v microseconds", time.Now().Sub(tFrom).Nanoseconds()/1000))
	}
	turnInfo("End status:", status())
}

//Prints the movements of own drones
func play() {
	initializeTurnComputation()

	calculateDonesAvailableDistances()

	strategyMaintainAirSuperiority()
	//strategyColonizeTheUnexplored()
	//strategyGoForUnguardedZones()
	strategyAttack()
	strategyDefaultToCentroid()
	//strategyDefaultToNearestZone()
	for _, m := range nextMove {
		fmt.Println(m.x, m.y)
	}
}

/* TURN BEGIN/END - RELATED OPERATIONS END ****************************************DEBUG - RELATED OPERATIONS BEGIN***/

//Returns the status in a format that can be directly imported for testing
func importableStatus() string {
	var result bytes.Buffer
	result.Write([]byte(fmt.Sprintf("\n%d %d %d %d\n", numPlayers, whoami, numDronesPerplayer, numZones)))
	for _, z := range zones {
		result.Write([]byte(fmt.Sprintf("%d %d\n", z.pos.x, z.pos.y)))
	}
	for _, z := range zones {
		result.Write([]byte(fmt.Sprintf("%d\n", z.owner)))
	}
	for _, p := range players {
		for _, d := range p.drones {
			result.Write([]byte(fmt.Sprintf("%d %d\n", d.x, d.y)))
		}
	}
	return result.String()
}

//Returns the status of the play if debug is enabled
func status() string {
	var result bytes.Buffer
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
		for _, d := range p.drones {
			result.Write([]byte(fmt.Sprintf("%v ", d)))
		}
		result.Write([]byte("]\n"))
	}
	result.Write([]byte("Zones:\n"))
	for zId, z := range zones {
		result.Write([]byte(fmt.Sprintf("  %d - owner: %d location: %v\n", zId, z.owner, z.pos)))
	}
	return result.String()
}

//Traces information to Standard Error output
func trace(x ...interface{}) {
	if DEBUG {
		fmt.Fprintln(os.Stderr, x)
	}
}

//Writes the main information regarding the actions taken in the turn
func turnInfo(x ...interface{}) {
	if realExecution {
		fmt.Fprintln(os.Stderr, x)
	}
}

/* DEBUG - RELATED OPERATIONS END ******************************************************************* IDEAS BEGIN *****/

/*
IDEAS:
- Calculate "centroid" of the board based on the zones' locations. Move drones to the position in the zone neares to the center of the board
- Different strategies depending on whether I am winning (based on actual score of all players, zones under my control and remaining turns)
- Take into account oponents' possible movements
- Take into account oponents' drones' distances to owned zones
- should air superiority include drones at distance 1?
- Last strategy: go for the "centroid" (the center of all zones)
- Decide nearest drone per euclidean distance instead of per turns?
- ¿Would it be better to begin conquering non-obvious zones?
- IMPORTANT!!! Implement a strategy to protect zones unsufficiently guarded
*/
