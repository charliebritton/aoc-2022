package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const ROPE_LENGTH = 10

func main() {

	f, err := os.OpenFile("../data/day9a", os.O_RDONLY, 0755)
	if err != nil {
		panic(err)
	}

	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)

	game := &Game{
		knots:       make([]Position, ROPE_LENGTH),
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
			// fmt.Println()

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
func tailNeedsToMove(head Position, tail Position) bool {

	// Head x = 2 Tail x = 0 needs to move
	// Head x = 2 Tail x = 1 no need to move
	// if (g.head.x-g.tail.x <= 1 || g.tail.x-g.head.x <= 1) && (g.head.y-g.tail.y <= 1 || g.tail.y-g.head.y <= 1) {
	// 	return false
	// }

	diffX := Abs(head.x - tail.x)
	diffY := Abs(head.y - tail.y)

	if diffX > 1 || diffY > 1 {
		return true
	}

	return false

}

type Size Position

type Game struct {
	knots       []Position
	tailVisited map[Position]struct{}
	boardSize   Size
	boardCorner Position // For where the origin is negative
}

// func (g Game) head() Position {
// 	return g.knots[0]
// }

// func (g Game) tail() Position {
// 	return g.knots[len(g.knots)-1]
// }

// Takes a direction to move
func step(d Direction, g *Game) {

	switch d {
	case Up:
		g.knots[0].y += 1
	case Down:
		g.knots[0].y -= 1
	case Left:
		g.knots[0].x -= 1
	case Right:
		g.knots[0].x += 1
	default:
		panic("invalid direction")
	}

	// Changes need to propogate through
	// For each knot acting as its own head, check tail needs to move for that specific, then if at end just return
	// So do H1, 12, 23, 34, 45, 56, 67, 78, 89

	for h := range g.knots[:len(g.knots)-1] {

		t := h + 1

		// fmt.Printf("[%d,%d] ", h, t)

		if tailNeedsToMove(g.knots[h], g.knots[t]) {

			moveTail(&g.knots[h], &g.knots[t])

		}

	}

	g.tailVisited[g.knots[len(g.knots)-1]] = struct{}{}

	// OOB top
	if g.knots[0].y >= g.boardCorner.y+g.boardSize.y {
		g.boardSize.y += 1
	}

	// OOB right
	if g.knots[0].x >= g.boardCorner.x+g.boardSize.x {
		g.boardSize.x += 1
	}

	// OOB bottom
	if g.knots[0].y < g.boardCorner.y {
		g.boardSize.y++
		g.boardCorner.y--
	}

	// OOB left
	if g.knots[0].x < g.boardCorner.x {
		g.boardSize.x++
		g.boardCorner.x--
	}

	// Update locations tail has visited and board size if necessary
	// fmt.Printf("Moved head %v\nKnot positions: %v\nBoard Size: %v\nBoard Corner: %v\n\n", d, g.knots, g.boardSize, g.boardCorner)

}

func printBoard(g *Game) {

	lineBuf := []string{}

	// Line by line, buffer output and then flip
	for y := g.boardCorner.y; y < g.boardSize.y; y++ {

		line := ""
		bracketed := ""
		for x := g.boardCorner.x; x < g.boardSize.x; x++ {

			// This is set to false when encountered
			top := true
			singular := true
			bracketedThis := ""
			isKnot := false

			// For each in the knots try and write there
			// If writing as 0, then use H

			for i, knot := range g.knots {

				// In this position
				if knot.x == x && knot.y == y {

					if top && i == 0 {
						line += "H"
						bracketedThis += "H covers "
						top = false
					} else if top {
						line += strconv.Itoa(i)
						bracketedThis += strconv.Itoa(i) + " covers "
						top = false
					} else {
						singular = false
						bracketedThis += strconv.Itoa(i) + ", "
					}

					isKnot = true

				}

			}

			// For s specifically
			if top && x == 0 && y == 0 {
				line += "s"
				bracketedThis += "s covers "
				top = false
			} else if x == 0 && y == 0 {
				singular = false
				bracketedThis += "s, "
			} else if !isKnot {
				line += "."
			}

			if singular {
				bracketedThis = ""
			}

			bracketed += bracketedThis

		}

		if bracketed != "" {
			line += "  (" + bracketed[:len(bracketed)-2] + ")"
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

func moveTail(h *Position, t *Position) {
	if h.y == t.y {

		// Head is to right
		if h.x > t.x {
			t.x++
		} else {
			t.x--
		}

		// Case they are on the same column
	} else if h.x == t.x {

		// Head is above
		if h.y > t.y {
			t.y++
		} else {
			t.y--
		}

		// Head x,y bigger
	} else if h.x > t.x && h.y > t.y {

		t.x++
		t.y++

		// Head x bigger but y smaller
	} else if h.x > t.x && h.y < t.y {

		t.x++
		t.y--

		// Head x,y smaller
	} else if h.x < t.x && h.y < t.y {

		t.x--
		t.y--

		// Head x smaller but y bigger
	} else if h.x < t.x && h.y > t.y {

		t.x--
		t.y++

	} else {

		panic("clearly missed a trick here chaz")

	}
}
