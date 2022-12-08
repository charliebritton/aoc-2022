package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var rows [][]int

func main() {

	f, err := os.OpenFile("../data/day8a", os.O_RDONLY, 0755)
	if err != nil {
		panic(err)
	}

	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {

		line := fileScanner.Text()

		row := strings.Split(line, "")

		var rowInt []int

		for _, col := range row {

			colInt, _ := strconv.Atoi(col)
			rowInt = append(rowInt, colInt)

		}

		fmt.Printf("%v\n", row)

		rows = append(rows, rowInt)

	}

	vis, score := calcArrays(rows)

	printArr(vis)
	printScores(score)
	// fmt.Printf("%t", isVisible(3, 3))

	// fmt.Printf("%v\n", vis)

}

func printArr(rows [][]bool) {

	totalVisible := 0
	for _, row := range rows {

		for _, col := range row {

			if col {
				fmt.Printf("x ")
				totalVisible++
			} else {
				fmt.Printf("o ")
			}

		}
		fmt.Println()

	}

	fmt.Printf("TOTAL VISIBLE: %d\n", totalVisible)

}

func printScores(rows [][]int) {

	bestScore := 0
	for _, row := range rows {

		for _, col := range row {

			if col > bestScore {
				bestScore = col
			}

			// fmt.Printf("%03d ", col)

		}
		// fmt.Println()

	}

	fmt.Printf("BEST SCORE: %d\n", bestScore)

}

func calcArrays(rows [][]int) ([][]bool, [][]int) {

	vis := make([][]bool, len(rows))
	scores := make([][]int, len(rows))

	for i := range vis {
		vis[i] = make([]bool, len(rows[0]))
	}

	for i := range scores {
		scores[i] = make([]int, len(rows[0]))
	}

	for y := range rows {

		for x := 0; x < len(rows[1]); x++ {

			vis[y][x] = isVisible(x, y)
			scores[y][x] = getScore(x, y)

		}
	}

	return vis, scores

}

func getScore(x, y int) int {

	// totalRows := len(rows[0])
	// totalCols := len(rows)

	var topScore, leftScore, rightScore, bottomScore int

	leftScore = scoreLeft(x, y)
	rightScore = scoreRight(x, y)
	bottomScore = scoreBottom(x, y)
	topScore = scoreTop(x, y)
	// bottomScore, topScore = 1, 1

	return topScore * leftScore * rightScore * bottomScore
}

func isVisible(x int, y int) bool {

	totalRows := len(rows[0])
	totalCols := len(rows)

	if atEdge(x, y) {
		// fmt.Printf("testing (%d,%d): at edge\n", x, y)
		return true
	}

	// Visible from left, right
	if visRange(x, y, 0, y) || visRange(x, y, totalCols-1, y) || visRange(x, y, x, totalRows-1) || visRange(x, y, x, 0) {
		// fmt.Printf("  ==> Tested up down left and right and was visible from one of them\n\n")
		return true
	}

	// fmt.Printf("  ==> Not visible at all\n\n")
	return false

}

// Is the tallest in the whole line
func visRange(x0, y0, x1, y1 int) bool {

	// fmt.Printf("testing (%d,%d): ", x0, y0)
	visible := true

	// Working vertically, so change y
	if x0 == x1 {
		// fmt.Println("This is a vertical check")

		// Left to right
		if y0 < y1 {
			for i := y0 + 1; i <= y1; i++ {
				// fmt.Printf("testing %d < %d from top\n", rows[y0][x0], rows[i][x0])

				if rows[y0][x0] <= rows[i][x0] {
					visible = false
				}

			}
			// Right to left
		} else {
			for i := y1; i < y0; i++ {
				// fmt.Printf("testing %d < %d from bottom\n", rows[y0][x0], rows[i][x0])
				if rows[y0][x0] <= rows[i][x0] {
					visible = false
				}
			}
		}

	}

	// Working horizontally, change x
	if y0 == y1 {
		// fmt.Println("This is a horizontal check")

		// Top to bottom
		if x0 < x1 {
			for i := x0 + 1; i <= x1; i++ {
				// fmt.Printf("testing %d < %d left to right\n", rows[y0][x0], rows[y0][i])

				if rows[y0][x0] <= rows[y0][i] {
					visible = false
				}

			}

			// Right to left
		} else {
			for i := x1; i < x0; i++ {
				// fmt.Printf("testing %d < %d right to left\n", rows[y0][x0], rows[y0][i])
				if rows[y0][x0] <= rows[y0][i] {
					visible = false
				}
			}
		}
	}
	// fmt.Printf("==> Tested visibility of tree with height %d in range (%d,%d) to (%d,%d): %t\n", rows[y0][x0], x0, y0, x1, y1, visible)

	return visible

}

func atEdge(x int, y int) bool {

	if x == 0 || y == 0 || x == len(rows[0])-1 || y == len(rows)-1 {
		return true
	}

	return false

}

func scoreLeft(x, y int) int {

	score := 0

	// Add to the score while trees are smaller
	// From current to x = 0

	for c := x - 1; c >= 0; c-- {

		if rows[y][c] < rows[y][x] {
			score++
		} else if rows[y][c] >= rows[y][x] {
			score++
			break
		}

	}

	// fmt.Printf("(%d,%d) SCORE TO LEFT: %d\n", x, y, score)
	return score

}

func scoreRight(x, y int) int {

	score := 0

	// Add to the score while trees are smaller
	// From current to x = 0

	for c := x + 1; c < len(rows[0]); c++ {
		// fmt.Printf("[%d,%d] ", c, y)
		if rows[y][c] < rows[y][x] {
			// fmt.Println("less than")
			score++
		} else if rows[y][c] >= rows[y][x] {
			// fmt.Println("Greater or Equal")
			score++
			break
		}

	}

	// fmt.Printf("(%d,%d) SCORE TO RIGHT: %d\n", x, y, score)
	return score

}

func scoreTop(x, y int) int {

	score := 0

	// Add to the score while trees are smaller
	// From current to x = 0

	for c := y - 1; c >= 0; c-- {

		if rows[c][x] < rows[y][x] {
			score++
		} else if rows[c][x] >= rows[y][x] {
			score++
			break
		}

	}

	// fmt.Printf("(%d,%d) SCORE TO TOP: %d\n", x, y, score)
	return score

}

func scoreBottom(x, y int) int {

	score := 0

	// Add to the score while trees are smaller
	// From current to x = 0

	for c := y + 1; c < len(rows); c++ {

		if rows[c][x] < rows[y][x] {
			score++
		} else if rows[c][x] >= rows[y][x] {
			score++
			break
		}

	}

	// fmt.Printf("(%d,%d) SCORE TO BOTTOM: %d\n", x, y, score)
	return score

}
