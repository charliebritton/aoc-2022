package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

/*
Want to find a path in as short a time as possible...
Do we use heuristics to see if we're getting closer and DFS

From each square, we have a list of 'allowed' moves
Can calculate the Euclidian distance from the final square

Will be useless to re-visit a square that we've visited before, so store a bitmap or something of all squared we've already been to.
Rule out any squares that we have already visited, then use BFS and work towards the goal using heuristics that you want to reduce Euclidian distance first
Second heuristic would be keeping on the same level or going up

Translate the input grid into numbers first

*/

// No tuple in Golang so have to use this
type Pos struct {
	x int
	y int
	p *Pos
}

// Stores overall state
type Puzzle struct {
	VisitMap  [][]bool
	HeightMap HeightMap // Mapping of letters to numbers
	Curr      Pos       // Current position in the search
	Start     Pos       // Start
	End       Pos       // End
	Trail     []Pos     // Stack of previous Currs
}

// Stores the map as a slice of ints
// Once a node has been visited the height is updated to 0.
type HeightMap = [][]int

// Stores the directions that can be taken
type Direction = int

const (
	Up = iota + 1
	Down
	Left
	Right
)

// Takes a puzzle and returns the Euclidian distance from the current position to the end
func (p Puzzle) Euclidian() float64 {

	return math.Sqrt(math.Pow((float64(p.End.x)-float64(p.Curr.x)), 2) + math.Pow((float64(p.End.y)-float64(p.Curr.y)), 2))

}

func (p Puzzle) EuclidianFromPos(pos Pos) float64 {

	return math.Sqrt(math.Pow((float64(p.End.x)-float64(pos.x)), 2) + math.Pow((float64(p.End.y)-float64(pos.y)), 2))

}

// For easy printing to the console
func (p Puzzle) Format(f fmt.State, c rune) {

	for y, row := range p.HeightMap {

		for x, col := range row {

			if x == p.Curr.x && y == p.Curr.y {
				f.Write([]byte(fmt.Sprintf("\033[1m%02d\033[0m ", col)))
			} else if p.VisitMap[y][x] == true {
				f.Write([]byte(fmt.Sprintf("\033[3m%02d\033[0m ", col)))
			} else {
				f.Write([]byte(fmt.Sprintf("%02d ", col)))

			}
		}

		f.Write([]byte(fmt.Sprintln()))

	}

	f.Write([]byte(fmt.Sprintf("\nCurr: %v, End: %v, Size: (%d,%d), Euclid: %.02f\n", p.Curr, p.End, len(p.HeightMap[0]), len(p.HeightMap), p.Euclidian())))

}

func (p Pos) Format(f fmt.State, c rune) {

	f.Write([]byte(fmt.Sprintf("(%d,%d)", p.x, p.y)))

}

func (puz Puzzle) onBoard(p Pos) bool {

	if p.x > -1 && p.y > -1 && p.x < len(puz.HeightMap[0]) && p.y < len(puz.HeightMap) {
		return true
	}

	return false

}

// Is legal provided that the square not already visited (0) and is at most 1 higher
func (puz Puzzle) isLegalMove(p Pos) bool {

	if !puz.onBoard(p) {
		fmt.Println("Not on board")
		return false
	}

	// Already visited
	// if puz.HeightMap[p.y][p.x] == 0 {
	// 	fmt.Println("Already visisted")
	// 	return false
	// }

	// At most 1 lower
	if puz.HeightMap[puz.Curr.y][puz.Curr.x]-puz.HeightMap[p.y][p.x] <= 1 {
		fmt.Printf("difference is %d\n", puz.HeightMap[puz.Curr.y][puz.Curr.x]-puz.HeightMap[p.y][p.x])
		return true
	}

	return false

}

// Gets the height difference (tiebreaker), is +ve if new square higher
func (puz Puzzle) heightDiff(p Pos) int {
	return puz.HeightMap[p.y][p.x] - puz.HeightMap[puz.Curr.y][puz.Curr.x]
}

// Takes the given current position and checks if the move is legal
func (puz Puzzle) getLegalMoves() []Pos {

	p := puz.Curr
	up := Pos{p.x, p.y - 1, nil}
	down := Pos{p.x, p.y + 1, nil}
	left := Pos{p.x - 1, p.y, nil}
	right := Pos{p.x + 1, p.y, nil}
	var positions = []Pos{up, down, left, right}

	var legalMoves []Pos

	for _, new := range positions {

		if puz.isLegalMove(new) {
			legalMoves = append(legalMoves, new)
		}

	}

	fmt.Printf("There were %d legal moves from %v: %v", len(legalMoves), puz.Curr, legalMoves)
	return legalMoves

}

func parseInput(fileScanner *bufio.Scanner) *Puzzle {

	// DEBUG to work out the offset
	// fmt.Printf("a to int: %d (%d), z to int: %d (%d), so conversion offset to get a 1 and z 26 should be: %d", int('a'), int('a')-96, int('z'), int('z')-96, int('a')-1)

	var p = &Puzzle{}
	currRow := 0

	// Each line is a row in the map
	for fileScanner.Scan() {

		var row []int
		var rowVisit []bool

		line := fileScanner.Text()

		for col, c := range line {

			if c == 'S' {
				p.Start = Pos{col, currRow, nil}
				row = append(row, 0) // Doesn't matter the height of an already visited node
				rowVisit = append(rowVisit, false)
				continue
			}

			if c == 'E' {
				p.End = Pos{col, currRow, nil}
				p.Curr = Pos{col, currRow, nil}
				row = append(row, 27) // Input shows that the height is greater than z which will be 26, 0 is reserved for visited already
				rowVisit = append(rowVisit, false)
				continue
			}

			row = append(row, int(c)-96)
			rowVisit = append(rowVisit, false)

		}

		currRow++
		p.HeightMap = append(p.HeightMap, row)
		p.VisitMap = append(p.VisitMap, rowVisit)

	}

	return p

}

type Heuristic struct {
	Pos        Pos
	Direction  Direction
	Legal      bool
	Euclid     float64 // Main heuristic
	ElevChange int     // Tie breaker
}

/*
	calcMoves takes the current state of a puzzle and returns a list of the moves to take, with the best move first

From each square, we have a list of 'allowed' moves
Can calculate the Euclidian distance from the final square

Will be useless to re-visit a square that we've visited before, so store a bitmap or something of all squared we've already been to.
Rule out any squares that we have already visited, then use BFS and work towards the goal using heuristics that you want to reduce Euclidian distance first
Second heuristic would be keeping on the same level or going up
*/
func calcMoves(p *Puzzle) []Direction {

	up := Pos{p.Curr.x, p.Curr.y - 1, nil}
	down := Pos{p.Curr.x, p.Curr.y + 1, nil}
	left := Pos{p.Curr.x - 1, p.Curr.y, nil}
	right := Pos{p.Curr.x + 1, p.Curr.y, nil}
	var positions = []Pos{up, down, left, right}

	var dirs = []Direction{Up, Down, Left, Right}
	var heuristics []Heuristic

	for i, pos := range positions {
		heuristics = append(heuristics, p.calcMove(pos, dirs[i]))
	}

	fmt.Println(heuristics)

	var filtered []Heuristic
	for _, h := range heuristics {
		if h.Legal {
			filtered = append(filtered, h)
		}
	}

	fmt.Println(filtered)

	// Sort on elevation change first, bigger is less in index
	sort.SliceStable(filtered, func(i, j int) bool {
		return float64(filtered[i].ElevChange) > float64(filtered[j].ElevChange)
	})

	// Then sort on Euclid which is the main
	sort.SliceStable(filtered, func(i, j int) bool {
		return filtered[i].Euclid < filtered[j].Euclid
	})

	fmt.Println(filtered)

	var dirsRtn []Direction
	for _, d := range filtered {
		dirsRtn = append(dirsRtn, d.Direction)
	}

	return dirsRtn

}

// Was going to be used for a DFS implemenetation from the start
func (puz Puzzle) calcMove(p Pos, d Direction) Heuristic {
	h := Heuristic{
		Pos:       p,
		Direction: d,
		Legal:     puz.isLegalMove(p),
	}

	if h.Legal {
		h.Euclid = puz.EuclidianFromPos(p)
		h.ElevChange = puz.heightDiff(p)
	}

	return h

}

// A queue stores a list of positions that we need to visit
type Queue []Pos

func (q *Queue) Enqueue(e Pos) {
	// fmt.Printf("Enqueuing %v\n", e)
	*q = append(*q, e)
	// fmt.Printf("Queue is %v\n", q)
}

func (q *Queue) Dequeue() Pos {
	rtn := (*q)[0]
	// fmt.Printf("Dequeuing %v\n", rtn)
	*q = (*q)[1:]
	return rtn
}

func (q Queue) IsEmpty() bool {
	return len(q) == 0
}

/*
Implementation parses the data into a mapping for us
Using DFS, we want to follow a path until we reach a node that is equal to the end
We then kill any searches that take more than the number of moves to find the node
Once all searches have either exhausted their total moves (by being equal to the
numMoves for the best so far) or they have no remaining directions to choose, we
have a winner
*/

func (p Puzzle) searchBFS() Pos {

	// Setup a Queue
	var q Queue

	// Mark as explored
	// p.HeightMap[p.End.y][p.End.x] = 0

	// Add the starting node to the queue
	q.Enqueue(p.End)

	// Whilst the queue is not empty, we look for legal moves that we can follow and do BFS on them there
	for !q.IsEmpty() {

		fmt.Println("Queue is not empty")
		fmt.Printf("Current Queue: %v\n", q)

		// Map at the start of this turn
		fmt.Printf("Queue at start of this round:\n\n%v", p)

		v := q.Dequeue()
		p.Curr = v

		fmt.Printf("Checking %v, queue is %v\n", v, q)

		if v.x == p.Start.x && v.y == p.Start.y {
			fmt.Printf("We have reached the target at %v", p.Start)
			return v
		}

		// Get elements that are legal from current position
		for _, e := range p.getLegalMoves() {

			// Not already visited
			if !p.VisitMap[e.y][e.x] {

				fmt.Printf("Adding %v to the queue\n", e)

				// Mark as visited
				p.VisitMap[e.y][e.x] = true

				// Store the parent
				e.p = &v

				// Add this node to the queue
				q.Enqueue(e)

			}

		}

	}
	fmt.Println("Queue is empty")

	return Pos{0, 0, nil}

}

func part1(fileName string) {

	f, err := os.OpenFile("../data/"+fileName, os.O_RDONLY, 0755)
	if err != nil {
		panic(err)
	}

	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)
	p := parseInput(fileScanner)
	fmt.Println(p)

	fmt.Printf("%v\n", calcMoves(p))

	bfsRes := p.searchBFS()

	fmt.Printf("Following trail through START -> ")
	l := trailLength(&bfsRes)
	fmt.Printf("Had to explore %d nodes\n", l)

}

func trailLength(p *Pos) int {

	if p.p == nil {
		fmt.Printf("%v -> END ", p)
		fmt.Println()
		return 0
	}

	fmt.Printf("%v -> ", p)
	return 1 + trailLength(p.p)

}

func part2(fileName string) {

	f, err := os.OpenFile("../data/"+fileName, os.O_RDONLY, 0755)
	if err != nil {
		panic(err)
	}

	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)
	p := parseInput(fileScanner)
	fmt.Println(p)

	// fmt.Printf("%v\n", calcMoves(p))

	// bfsRes := p.searchBFS()

	// fmt.Printf("Following trail through START -> ")
	// l := trailLength(&bfsRes)
	// fmt.Printf("Had to explore %d nodes\n", l)

}

func main() {

	part1("day12t")
	// part1("day12a")
	part2("day12t")
	// part2("day12a")
}
