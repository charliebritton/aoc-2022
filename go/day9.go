package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {

	f, err := os.OpenFile("../data/day9a", os.O_RDONLY, 0755)
	if err != nil {
		panic(err)
	}

	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)

	game := &Game{
		head:        Position{0, 0},
		tail:        Position{0, 0},
		tailVisited: make(map[Position]struct{}),
		boardSize:   Size{6, 5},
	}

	fmt.Printf("== Initial State ==\n\n\n")
	printBoard(game)

	for fileScanner.Scan() {

		line := fileScanner.Text()

		// fmt.Printf("\n== %s ==\n\n", strings.Split(line, "\n")[0])

		numMoves := strings.Split(line, " ")[1]
		numMovesInt, _ := strconv.Atoi(numMoves)

		dirString := strings.Split(line, " ")[0]

		for i := 0; i < numMovesInt; i++ {

			switch dirString {
			case "U":
				step(Up, game)
			case "D":
				step(Down, game)
			case "L":
				step(Left, game)
			case "R":
				step(Right, game)
			default:
				panic("Error parsing input")
			}

			// printBoard(game)
			fmt.Println()

		}

	}

	printTailVisited(game)

}

// Starts from bottom left = 0,0
type Position struct {
	x int
	y int
}

type Direction int

const (
	Up = iota + 1
	Down
	Left
	Right
)

func (e Direction) String() string {
	switch e {
	case Up:
		return "Up"
	case Down:
		return "Down"
	case Left:
		return "Left"
	case Right:
		return "Right"
	default:
		return fmt.Sprintf("%d", int(e))
	}
}

// Abs returns the absolute value of x.
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// tailNeedsToMove takes a Position of the head and the position of the tail
// TODO possible bug if coordinate space changes to include negative
func tailNeedsToMove(g *Game) bool {

	// Head x = 2 Tail x = 0 needs to move
	// Head x = 2 Tail x = 1 no need to move
	// if (g.head.x-g.tail.x <= 1 || g.tail.x-g.head.x <= 1) && (g.head.y-g.tail.y <= 1 || g.tail.y-g.head.y <= 1) {
	// 	return false
	// }

	diffX := Abs(g.head.x - g.tail.x)
	diffY := Abs(g.head.y - g.tail.y)

	if diffX > 1 || diffY > 1 {
		return true
	}

	return false

}

type Size Position

type Game struct {
	head, tail  Position
	tailVisited map[Position]struct{} // always inc 00
	boardSize   Size
	boardCorner Position // For where the origin is negative
}

// Takes a direction to move
func step(d Direction, g *Game) {

	switch d {
	case Up:
		g.head.y += 1
	case Down:
		g.head.y -= 1
	case Left:
		g.head.x -= 1
	case Right:
		g.head.x += 1
	default:
		panic("invalid direction")
	}

	if tailNeedsToMove(g) {

		// Case they are on the same row
		if g.head.y == g.tail.y {

			// Head is to right
			if g.head.x > g.tail.x {
				g.tail.x++
			} else {
				g.tail.x--
			}

			// Case they are on the same column
		} else if g.head.x == g.tail.x {

			// Head is above
			if g.head.y > g.tail.y {
				g.tail.y++
			} else {
				g.tail.y--
			}

			// Head x,y bigger
		} else if g.head.x > g.tail.x && g.head.y > g.tail.y {

			g.tail.x++
			g.tail.y++

			// Head x bigger but y smaller
		} else if g.head.x > g.tail.x && g.head.y < g.tail.y {

			g.tail.x++
			g.tail.y--

			// Head x,y smaller
		} else if g.head.x < g.tail.x && g.head.y < g.tail.y {

			g.tail.x--
			g.tail.y--

			// Head x smaller but y bigger
		} else if g.head.x < g.tail.x && g.head.y > g.tail.y {

			g.tail.x--
			g.tail.y++

		} else {

			panic("clearly missed a trick here chaz")

		}

	}

	g.tailVisited[g.tail] = struct{}{}

	// OOB top
	if g.head.y >= g.boardCorner.y+g.boardSize.y {
		g.boardSize.y += 1
	}

	// OOB right
	if g.head.x >= g.boardCorner.x+g.boardSize.x {
		g.boardSize.x += 1
	}

	// OOB bottom
	if g.head.y < g.boardCorner.y {
		g.boardSize.y++
		g.boardCorner.y--
	}

	// OOB left
	if g.head.x < g.boardCorner.x {
		g.boardSize.x++
		g.boardCorner.x--
	}

	// Update locations tail has visited and board size if necessary
	fmt.Printf("Moved head %v\nHead at: %v\nTail at: %v\nBoard Size: %v\nBoard Corner: %v\n\n", d, g.head, g.tail, g.boardSize, g.boardCorner)

}

func printBoard(g *Game) {

	lineBuf := []string{}

	// Line by line, buffer output and then flip
	for y := g.boardCorner.y; y < g.boardSize.y; y++ {

		line := ""
		for x := g.boardCorner.x; x < g.boardSize.x; x++ {

			// If head at position
			if g.head.x == x && g.head.y == y {
				line += "H"
			} else if g.tail.x == x && g.tail.y == y {
				line += "T"
			} else if x == 0 && y == 0 {
				line += "s"
			} else {
				line += "."
			}

		}

		lineBuf = append(lineBuf, line)

	}

	// Reverse
	for i, j := 0, len(lineBuf)-1; i < j; i, j = i+1, j-1 {
		lineBuf[i], lineBuf[j] = lineBuf[j], lineBuf[i]
	}

	for _, line := range lineBuf {

		fmt.Println(line)

	}

}

func printTailVisited(g *Game) {

	totalVisits := 0
	lineBuf := []string{}

	// Line by line, buffer output and then flip
	for y := g.boardCorner.y; y < g.boardSize.y; y++ {

		line := ""
		for x := g.boardCorner.x; x < g.boardSize.x; x++ {

			// If head at position
			if x == 0 && y == 0 {
				totalVisits++
				line += "s"

			} else if _, visited := g.tailVisited[Position{x, y}]; visited {

				totalVisits++
				line += "#"

			} else {
				line += "."
			}

		}

		lineBuf = append(lineBuf, line)

	}

	// Reverse
	for i, j := 0, len(lineBuf)-1; i < j; i, j = i+1, j-1 {
		lineBuf[i], lineBuf[j] = lineBuf[j], lineBuf[i]
	}

	for _, line := range lineBuf {

		fmt.Println(line)

	}

	fmt.Printf("TOTAL TAIL VISITS: %d\n", totalVisits)
}
